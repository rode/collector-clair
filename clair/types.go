// Copyright 2021 The Rode Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
