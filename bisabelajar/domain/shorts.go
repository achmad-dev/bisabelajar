package domain

import "time"

type Shorts struct {
	ID        int       `json:"id"`
	SeriesID  int       `json:"series_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
