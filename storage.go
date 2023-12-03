package sh

import (
	"errors"
	"os"
	"path/filepath"
)

type Storage struct {
	Dir          string
	NotAutoClean bool
}

func (s *Storage) RemoveOrTruncate(file *os.File, size int64) error {
	if file == nil {
		return nil
	}
	if s.NotAutoClean {
		return s.Truncate(file, size)
	}
	return os.Remove(file.Name())
}

func (s *Storage) Truncate(file *os.File, size int64) error {
	if file == nil || size == 0 {
		return nil
	}
	if stat, err := file.Stat(); err != nil {
		return err
	} else {
		if size >= stat.Size() {
			return file.Truncate(stat.Size() - 1)
		}
		return file.Truncate(stat.Size() - size)
	}
}

func (s *Storage) CreateFile(filename string) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("filename is empty")
	}
	return os.Create(filepath.Join(s.Dir, filename))
}

func CheckStorage(s *Storage) (bool, error) {
	if s == nil || s.Dir == "" {
		return false, nil
	}
	f, err := os.Stat(s.Dir)
	if err != nil {
		return false, err
	}
	if !f.IsDir() {
		return false, &os.PathError{Op: "IsDir", Path: s.Dir, Err: errors.New("no such directory")}
	}

	return true, nil
}
