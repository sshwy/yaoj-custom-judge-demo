package service

import (
	"io"
	"log"
	"os"
	"path"
)

// 核心文件存储服务
type FileService struct {
	tempDir string
}

var TempFile = FileService{
	tempDir: "/tmp/yaoj-go",
}

// store content into f.tmpDir/name
func (f *FileService) Store(name string, reader io.Reader) error {
	fullpath := path.Join(f.tempDir, name)

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}
	return nil
}

func (f *FileService) Touch(name string) error {
	fullpath := path.Join(f.tempDir, name)

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (f *FileService) Remove(name string) error {
	fullpath := path.Join(f.tempDir, name)

	if err := os.Remove(fullpath); err != nil {
		return err
	}
	return nil
}

func (f *FileService) Pathof(name string) string {
	return path.Join(f.tempDir, name)
}

func (f *FileService) Create(name string) (*os.File, error) {
	fullpath := path.Join(f.tempDir, name)
	return os.Create(fullpath)
}

func (f *FileService) ReadFile(name string) ([]byte, error) {
	fullpath := path.Join(f.tempDir, name)
	return os.ReadFile(fullpath)
}

func init() {
	log.SetPrefix("[yaoj-go/service] ")
	os.MkdirAll("/tmp/yaoj-go", os.ModePerm)
	log.Println("fs init")
}
