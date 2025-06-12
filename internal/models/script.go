package models

import "time"

type TransformationScript struct {
	ID          string    `json:"id"`
	Identifier  string    `json:"identifier"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Script      string    `json:"script"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"createdBy"`
}
