package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func CheckPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

func RelPath(root string, path string) (string, error) {
	s, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	return s, nil
}

func RootPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("error: could not get root path: %s", err)
		return "", err
	}
	return wd, nil
}

func AbsolutePath(relPath string) (string, error) {
	root, err := RootPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, relPath), nil
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func CreateDirIfNotExist(path string) error {
	if err := CheckPath(path); err != nil {
		if err = CreateDir(path); err != nil {
			return err
		}
	}
	return nil
}

func RemoveFileIfExist(path string) error {
	if err := CheckPath(path); err != nil {
		if err = RemoveFile(path); err != nil {
			return err
		}
	}
	return nil
}

func CreateFile(path string, content []byte) (*os.File, error) {
	var (
		f   *os.File
		err error
	)

	if err = CheckPath(path); err != nil {
		if f, err = OpenFile(path); err != nil {
			return nil, err
		}
		if err = WriteFile(f, content); err != nil {
			return nil, err
		}
		defer CloseFile(f)
	}
	return f, nil
}

func WriteFile(file *os.File, content []byte) error {
	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}

func CloseFile(file *os.File) error {
	return file.Close()
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

func OpenFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, err
}

func RemoveRootDir(relPath string) error {
	relParent := strings.Split(relPath, "/")[0]
	parent, err := AbsolutePath(relParent)
	if err != nil {
		return err
	}
	return RemoveDir(parent)
}

func IsWrite(ev fsnotify.Event) bool {
	return ev.Op&fsnotify.Write == fsnotify.Write
}

func IsRemove(ev fsnotify.Event) bool {
	return ev.Op&fsnotify.Remove == fsnotify.Remove
}

func IsCreate(ev fsnotify.Event) bool {
	return ev.Op&fsnotify.Create == fsnotify.Create
}
