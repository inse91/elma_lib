package e365_gateway

import "time"

type Common struct {
	ID                  string    `json:"__id,omitempty"`
	Name                string    `json:"__name,omitempty"`
	CreatedAt           time.Time `json:"__createdAt,omitempty"`
	DeletedAt           time.Time `json:"__deletedAt,omitempty"`
	CreatedBy           string    `json:"__createdBy,omitempty"`
	Index               int       `json:"__index,omitempty"`
	Subscribers         []string  `json:"__subscribers,omitempty"`
	UpdatedAt           time.Time `json:"__updatedAt,omitempty"`
	UpdatedBy           string    `json:"__updatedBy,omitempty"`
	Version             int       `json:"__version,omitempty"`
	Debug               bool      `json:"__debug,omitempty"`
	ExternalID          string    `json:"__externalId,omitempty"`
	ExternalProcessMeta any       `json:"__externalProcessMeta,omitempty"`
	Status              Status    `json:"__status,omitempty"`
}

type Status struct {
	Order  int `json:"order,omitempty"`
	Status int `json:"status,omitempty"`
}
