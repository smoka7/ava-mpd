package playlist

type Queue struct {
	Length          int
	CurrentSongPage uint
	CurrentPage     uint
	LastPage        uint
	Duration        float64
	Albums          []Album
}

type Album struct {
	Songs               []Song
	Album, Artist, Date string
}

type Song struct {
	Title, Pos, Id, Track, Duration string
}

type Folder struct {
	Directory string
}

type File struct {
	File string
}

type ServerList struct {
	Folders []Folder
	Files   []File
}

type SearchResult map[string][]File

var err error
