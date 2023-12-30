package model

import "github.com/lib/pq"

type Track struct {
	ID           uint          `gorm:"primary_key"`
	ISRC         string        `gorm:"unique_index"`
	SpotifyImage string
	Title        string
	ArtistNames  pq.StringArray `gorm:"type:text[]"`
	Popularity   int
}

type RequestBody struct {
    ISRC string `json:"isrc" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TrackDetails struct {
	ID           uint   `json:"primary_key"`
	ISRC         string `json:"isrc"`
	SpotifyImage string `json:"image"`
	Title        string `json:"title"`
	ArtistNames  string		`json:"artist"`
	Popularity   int    `json:"popularity"`
}
