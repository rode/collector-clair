package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rode/collector-sonatype/sonatype"
	"github.com/rode/rode/protodeps/grafeas/proto/v1beta1/common_go_proto"
	"github.com/rode/rode/protodeps/grafeas/proto/v1beta1/package_go_proto"
	"github.com/rode/rode/protodeps/grafeas/proto/v1beta1/vulnerability_go_proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/rode/rode/proto/v1alpha1"
	"github.com/rode/rode/protodeps/grafeas/proto/v1beta1/grafeas_go_proto"
)

type listener struct {
	rodeClient pb.RodeClient
	logger     *zap.Logger
}

type Listener interface {
	ProcessEvent(http.ResponseWriter, *http.Request)
}

// NewListener instantiates a listener including a zap logger and the rodeclient connection
func NewListener(logger *zap.Logger, client pb.RodeClient) Listener {
	return &listener{
		rodeClient: client,
		logger:     logger,
	}
}

// ProcessEvent handles incoming webhook events
func (l *listener) ProcessEvent(w http.ResponseWriter, request *http.Request) {
	log := l.logger.Named("ProcessEvent")

	event := &sonatype.Event{}
	if err := json.NewDecoder(request.Body).Decode(event); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error reading webhook event")
		log.Error("error reading webhook event", zap.NamedError("error", err))
		return
	}

	log.Debug("received sonatype event", zap.Any("event", event), zap.Any("project", event.Image), zap.Any("image", event.Image))

	var occurrences []*grafeas_go_proto.Occurrence
	for _, vulnerability := range event.Vulnerabilities {
		log.Debug("sonatype vulnerability received", zap.Any("vulnerability", vulnerability.Vulnerability))
		// TODO determine method for getting the repo name if necessary
		occurrence := createQualityGateOccurrence(&vulnerability, "temp-repo")
		occurrences = append(occurrences, occurrence)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	response, err := l.rodeClient.BatchCreateOccurrences(ctx, &pb.BatchCreateOccurrencesRequest{
		Occurrences: occurrences,
	})
	if err != nil {
		log.Error("error sending occurrences to rode", zap.NamedError("error", err))
		w.WriteHeader(500)
		return
	}

	log.Debug("response payload", zap.Any("response", response.GetOccurrences()))
	w.WriteHeader(200)
}

func createQualityGateOccurrence(vulnerability *sonatype.Vulnerability, repo string) *grafeas_go_proto.Occurrence {
	occurrence := &grafeas_go_proto.Occurrence{
		Name: vulnerability.Vulnerability,
		Resource: &grafeas_go_proto.Resource{
			Name: repo,
			Uri:  repo,
		},
		// TODO change note name
		NoteName:    "projects/notes_project/notes/s",
		Kind:        common_go_proto.NoteKind_NOTE_KIND_UNSPECIFIED,
		Remediation: vulnerability.Fixedby,
		CreateTime:  timestamppb.Now(),

		Details: &grafeas_go_proto.Occurrence_Vulnerability{
			Vulnerability: &vulnerability_go_proto.Details{
				Type:             "container-vulnerability",
				Severity:         vulnerability_go_proto.Severity(vulnerability_go_proto.Severity_value[strings.ToUpper(vulnerability.Severity)]),
				ShortDescription: vulnerability.Description,
				LongDescription:  vulnerability.Description,
				RelatedUrls: []*common_go_proto.RelatedUrl{
					{
						Url:   vulnerability.Link,
						Label: vulnerability.Link,
					},
					{
						Url:   vulnerability.Link,
						Label: vulnerability.Link,
					},
				},
				EffectiveSeverity: vulnerability_go_proto.Severity(vulnerability_go_proto.Severity_value[strings.ToUpper(vulnerability.Severity)]),
				PackageIssue: []*vulnerability_go_proto.PackageIssue{
					{
						SeverityName: vulnerability.Vulnerability,
						AffectedLocation: &vulnerability_go_proto.VulnerabilityLocation{
							CpeUri:  vulnerability.Vulnerability,
							Package: vulnerability.Namespace,
							Version: &package_go_proto.Version{
								Name:     vulnerability.Featurename,
								Revision: vulnerability.Featureversion,
								Epoch:    35,
								Kind:     package_go_proto.Version_MINIMUM,
							},
						},
					},
				},
			},
		},
	}
	return occurrence
}
