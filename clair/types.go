package clair

// Event received
type Event struct {
	Image           string          `json:"image"`
	Unapproved      []string        `json:"unapproved"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// Vulnerability indicating the severity and vulnerability name
type Vulnerability struct {
	Featurename    string `json:"featurename"`
	Featureversion string `json:"featureversion"`
	Vulnerability  string `json:"vulnerability"`
	Namespace      string `json:"namespace"`
	Description    string `json:"description"`
	Link           string `json:"link"`
	Severity       string `json:"severity"`
	Fixedby        string `json:"fixedby"`
}
