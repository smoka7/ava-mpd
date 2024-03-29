package playlist

type Queue struct {
	Albums          Albums `json:",omitempty"`
	Length          uint
	CurrentSongPage uint
	CurrentPage     uint
	LastPage        uint
	Duration        uint
}

type Album struct {
	Album, Artist, Date string
	Songs               Songs
}

type Albums []Album

type Song struct {
	Album    string `json:",omitempty"`
	Artist   string `json:",omitempty"`
	Title    string `json:",omitempty"`
	Pos      string `json:",omitempty"`
	Id       string `json:",omitempty"`
	File     string `json:",omitempty"`
	Track    string `json:",omitempty"`
	Duration string `json:",omitempty"`
	Liked    bool
}

type Songs []Song

type Playlist struct {
	Name       string
	SongsCount uint
	Duration   uint
}

type Playlists []Playlist

type Directory struct {
	Directory string
}

type Directories []Directory

type File struct {
	File string
}

type Files []File

type ServerList struct {
	Directories Directories
	Files       Files
}

type SearchResult map[string]Files
