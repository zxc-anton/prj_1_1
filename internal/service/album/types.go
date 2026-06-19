package album

import (
	"time"
)

type Author struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Cover *string `json:"cover"`
}

type Album struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Cover     *string   `json:"cover"`
	CreatedAt time.Time `json:"created_At"`
	Authors   []Author  `json:"author"`
}
