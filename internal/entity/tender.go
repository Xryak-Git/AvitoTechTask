package entity

import "time"

type Tender struct {
	Id             int       `db:"id"`
	Name           string    `db:"balance"`
	Description    string    `db:"description"`
	ServiceType    string    `db:"service_type"`
	Status         string    `db:"status"`
	OrganizationId *int      `db:"organization_id"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
}
