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
	ImageURL               string    `json:"image_url,omitempty"`
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

// PexelsResponse represents the response from Pexels API
type PexelsResponse struct {
	Photos []PexelsPhoto `json:"photos"`
}

type PexelsPhoto struct {
	Src PexelsSrc `json:"src"`
}

type PexelsSrc struct {
	Medium string `json:"medium"`
}
