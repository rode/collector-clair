package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/rode/rode/proto/v1alpha1"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/rode/collector-sonatype/listener"
)

var (
	debug    bool
	port     int
	rodeHost string
)

func main() {
	flag.IntVar(&port, "port", 8081, "the port that the sonatype collector should listen on")
	flag.BoolVar(&debug, "debug", false, "when set, debug mode will be enabled")
	flag.StringVar(&rodeHost, "rode-host", "localhost:50051", "the host to use to connect to rode")

	flag.Parse()

	logger, err := createLogger(debug)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	conn, err := grpc.Dial(rodeHost, grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	if err != nil {
		logger.Fatal("failed to establish grpc connection to Rode API", zap.NamedError("error", err))
	}

	rodeClient := pb.NewRodeClient(conn)

	l := listener.NewListener(logger.Named("listener"), rodeClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/event", l.ProcessEvent)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "I'm healthy") })
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			logger.Fatal("could not start http server...", zap.NamedError("error", err))
		}
	}()

	logger.Info("listening for Sonatype events", zap.String("host", server.Addr))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	terminationSignal := <-sig
	logger.Info("shutting down...", zap.String("termination signal", terminationSignal.String()))

	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Fatal("could not shutdown http server...", zap.NamedError("error", err))
	}
}

func createLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}
