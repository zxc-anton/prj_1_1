package track

import "time"

type Author struct {
	Id       int     `json:"id"`
	PhotoUrl *string `json:"photo_url"`
	Name     string  `json:"name"`
}

type Track struct {
	Id        int       `json:"id"`
	FileURL   *string   `json:"file_url"`
	AlbumId   *int      `json:"album_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`

	Authors []Author `json:"authors"`
}
