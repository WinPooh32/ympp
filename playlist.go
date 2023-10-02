package ympp

type Tracklist struct {
	Playlist PlaylistWithTracks `json:"playlist"`
}

type Playlist struct {
	Revision             int    `json:"revision"`
	Kind                 int    `json:"kind"`
	Title                string `json:"title"`
	Description          string `json:"description"`
	DescriptionFormatted string `json:"descriptionFormatted"`
	TrackCount           int    `json:"trackCount"`
	Available            bool   `json:"available"`
}

type PlaylistWithTracks struct {
	Playlist
	Tracks []Track `json:"tracks"`
}

type Track struct {
	ID      string   `json:"id"`
	RealID  string   `json:"realID"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
