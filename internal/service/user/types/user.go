package types

type UserResponse struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Avatar_URL *string `json:"avatar_url"`
}

type Playlist struct {
	Id        int     `json:"id" db:"ID"`
	Name      string  `json:"name" db:"name"`
	Cover_URL *string `json:"cover_url" db:"cover_url"`
	Is_public bool    `json:"is_public" db:"is_public"`
}

type Track struct {
	Id       int     `json:"id"`
	File_url string  `json:"file_url"`
	Album_Id *int    `json:"album_id"`
	Album    *string `json:"album"`
	Title    string  `json:"title"`
}

type CreatePlaylist struct {
	IsPublic bool   `json:"is_public" validate:"required"`
	Name     string `json:"name" validate:"min=1"`
}

//easyjson:json
type ListenHistory []Track

//easyjson:json
type MyPLaylists []Playlist
