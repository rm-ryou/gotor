package fs

import "os"

type FileIO struct{}

func NewFileIO() *FileIO {
	return &FileIO{}
}

func (f *FileIO) Read(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
