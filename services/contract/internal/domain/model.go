package domain

import "time"

type Contract struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Desc      string    `json:"description,omitempty"`
	Version   int64     `json:"version"`
}
