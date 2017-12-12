package api

import (
	"time"
)

// MilestoneData represents a milestone of a repo/PR
type MilestoneData struct {
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	LabelsURL    string    `json:"labels_url"`
	ID           int       `json:"id"`
	Number       int       `json:"number"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      UserData  `json:"creator"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}
