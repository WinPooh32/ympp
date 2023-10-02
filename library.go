package ympp

type OwnerInfo struct {
	Visibility string `json:"visibility"`
	HasTracks  bool   `json:"hasTracks"`
	Owner      Owner  `json:"owner"`
}

type Owner struct {
	UID   string `json:"uid"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Library struct {
	Playlists []LibraryPlaylist `json:"playlists"`
}

type LibraryPlaylist struct {
	Playlist
}
