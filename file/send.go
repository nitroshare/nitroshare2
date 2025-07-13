package file

import (
	"io"
	"os"
	"path"
)

// SendFile implements Item for transferring files.
type SendFile struct {
	basePath string
	name     string
	info     os.FileInfo
}

// NewSendFile creates a new file item for the provided file base path & name.
func NewSendFile(basePath, name string) (*SendFile, error) {
	i, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return &SendFile{
		basePath: basePath,
		name:     name,
		info:     i,
	}, nil
}

func (s *SendFile) Type() string {
	return Type
}

func (s *SendFile) Size() int64 {
	return s.info.Size()
}

// TODO: cross-platform method for obtaining file read / mod / permissions

func (s *SendFile) Meta() (any, error) {
	return map[string]any{
		"filename": path.Join(s.basePath, s.info.Name()),
	}, nil
}

func (s *SendFile) Open() (io.ReadCloser, error) {
	return os.Open(s.name)
}
