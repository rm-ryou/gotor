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

func (f *FileIO) Write(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (f *FileIO) Delete(path string) error {
	return os.Remove(path)
}
