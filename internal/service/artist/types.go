package artist

type Artist struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Cover *string `json:"cover"`
}

type ArtistResponse struct {
	Artist             Artist `json:"artist"`
	LastMonthListeners int    `json:"lastMonthListeners"`
}

type Album struct {
	Id    int     `json:"id"`
	Title string  `json:"title"`
	Cover *string `json:"cover"`

	Artists []Artist `json:"artists"`
}
