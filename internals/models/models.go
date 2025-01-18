package models

import (
	"time"
)

type Entry struct {
	ID       int       `json:"id"`
	Url      string    `json:"url"`
	ShortUrl string    `json:"short_url"`
	Date     time.Time `json:"date"`
}
