package models

import "time"

type StartHub struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Location               string    `json:"location"`
	TeamSize               int       `json:"team_size"`
	URL                    string    `json:"url"`
	Email                  string    `json:"email"`
	JoinDate               time.Time `json:"join_date"`
	Categories             []string  `json:"categories,omitempty"`
	CollaboratingStarthubs []string  `json:"collaborating_starthubs,omitempty"`
	ExternalCollaborators  []string  `json:"external_collaborators,omitempty"`
}

// CreateStartHubRequest represents the request body for creating a starthub
type CreateStartHubRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	TeamSize    int      `json:"team_size"`
	URL         string   `json:"url"`
	Email       string   `json:"email" validate:"required,email"`
	Categories  []string `json:"categories"`
}
