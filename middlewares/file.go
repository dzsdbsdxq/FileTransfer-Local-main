package middlewares

type File struct {
	Id        string `json:"-"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Type      int    `json:"type"`
	Driver    string `json:"driver"`
	UpdatedAt string `json:"updated_at"`
	Thumbnail string `json:"thumbnail"`
	Url       string `json:"url"`
	SizeStr   string `json:"size_str"`
	TimeStr   string `json:"time_str"`
}

func (f File) GetSize() uint64 {
	return uint64(f.Size)
}

func (f File) GetName() string {
	return f.Name
}

func (f File) ModTime() string {
	return f.UpdatedAt
}

func (f File) IsDir() bool {
	return f.Type == 1
}

func (f File) GetType() int {
	return f.Type
}
