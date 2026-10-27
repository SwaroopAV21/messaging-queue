package models

import (
	"time"

)

type Status int

const (
	NOTGOOD    Status = -1
	PROCESSING Status = 1
	COMPLETE   Status = 0
)

type Response struct {
	RequestID     int64           `json:"requestID"`
	Api           string          `json:"requestAPI"`
	ReturnCode    int             `json:"rc"`
	RequestStatus string          `json:"status"`
}

type ProfileResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ShortDesc   string     `json:"shortDesc"`
	TechStack   []string   `json:"techStack"`
	GithubURL   string     `json:"githubUrl"`
	LiveURL     string     `json:"liveUrl"`
	ImageURL    string     `json:"imageUrl"`
	Featured    bool       `json:"featured"`
	SortOrder   int        `json:"sortOrder"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

