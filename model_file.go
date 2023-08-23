package e365_gateway

import "time"

type File struct {
	ID           string    `json:"__id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"originalName"`
	Directory    string    `json:"directory"`
	Size         int       `json:"size"`
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"__createdAt"`
	UpdatedAt    time.Time `json:"__updatedAt"`
	DeletedAt    time.Time `json:"__deletedAt"`
	CreatedBy    string    `json:"__createdBy"`
	UpdatedBy    string    `json:"__updatedBy"`
}

type DirectoryInfo struct {
	ID          string    `json:"__id"`
	Name        string    `json:"__name"`
	System      bool      `json:"system"`
	Directory   string    `json:"directory"`
	CreatedAt   time.Time `json:"__createdAt"`
	CreatedBy   string    `json:"__createdBy"`
	DeletedAt   time.Time `json:"__deletedAt"`
	UpdatedAt   time.Time `json:"__updatedAt"`
	UpdatedBy   string    `json:"__updatedBy"`
	ParentsList []string  `json:"parentsList"`
	UniqueNames bool      `json:"uniqueNames"`
}
