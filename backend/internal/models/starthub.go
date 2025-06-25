package models

type StartHub struct {
	ID                     int      `json:"id"`
	Name                   string   `json:"name"`
	Category               []string `json:"category"`
	Description            string   `json:"description"`
	Location               string   `json:"location"`
	TeamSize               int      `json:"team_size"`
	URL                    string   `json:"url"`
	Email                  string   `json:"email"`
	CollaboratingStarthubs []int    `json:"collaborating_starthubs"`
	CollaboratingCompanies []string `json:"collaborating_companies"`
	Investors              []int    `json:"investors"`
	Donators               []int    `json:"donators"`
	JoinDate               string   `json:"join_date"`
}
