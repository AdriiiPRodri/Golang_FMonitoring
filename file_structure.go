package main

type File struct {
	Filename  string `json:"filename,omitempty"`
	Size      string `json:"size"`
	Mode      string `json:"mode"`
	ModTime   string `json:"modtime"`
	IsDir     string `json:"isdir"`
	Inode     string `json:"inode"`
	HLinks    string `json:"n_hardlinks"`
	IsSymLink string `json:"issym"`
	RSymLink  string `json:"resolev_sym"`
	Checksum  string `json:"checksum"`
}
