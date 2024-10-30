package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ShipStoreDirName = ".ship"
)

func GetPackageStorePathByNames(packageNames ...string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	packageNames = append([]string{home, ShipStoreDirName}, packageNames...)
	return filepath.Join(packageNames...), nil
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetFilePathAndName(path string) (string, string) {
	fileName := filepath.Base(path)
	filePath := path[:len(path)-len(fileName)]
	return filePath, fileName
}

type File struct {
	path       string
	name       string
	outputFile *os.File
	resume     bool
}

func NewFile(path, name string, resume bool) (*File, error) {
	if !resume {
		name = getRandomTimeString() + filepath.Ext(name)
	}
	f := &File{path: path, name: name, resume: resume}
	err := f.handlerPath()
	if err != nil {
		return nil, err
	}
	err = f.handlerFile()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *File) Write(chunk []byte) error {
	if f.outputFile == nil {
		return fmt.Errorf("file is not open")
	}
	_, err := f.outputFile.Write(chunk)
	return err
}

func (f *File) Read() ([]byte, error) {
	if f.outputFile == nil {
		return nil, fmt.Errorf("file is not open")
	}
	return io.ReadAll(f.outputFile)
}

func (f *File) Close() error {
	if f.outputFile == nil {
		return fmt.Errorf("file is not open")
	}
	err := f.outputFile.Close()
	if err != nil {
		return err
	}
	f.outputFile = nil
	return nil
}

func (f *File) GetFileName() string {
	return f.name
}

func (f *File) GetFilePath() string {
	return f.path
}

func (f *File) GetFileFullPath() string {
	return f.path + f.name
}

func (f *File) ClearFileContent() error {
	return os.Truncate(f.path+f.name, 0)
}

func (f *File) handlerPath() error {
	if f.path == "" {
		return fmt.Errorf("path is empty")
	}
	if f.path[len(f.path)-1:] != "/" {
		f.path += "/"
	}
	if f.checkIsObjExist(f.path) {
		return nil
	}
	return f.createDir()
}

func (f *File) handlerFile() (err error) {
	if f.name == "" {
		return fmt.Errorf("name is empty")
	}
	if f.checkIsObjExist(f.path + f.name) {
		if f.resume {
			f.outputFile, err = os.OpenFile(f.path+f.name, os.O_APPEND|os.O_WRONLY, 0644)
			return err
		}
		err = f.deleteFile()
		if err != nil {
			return err
		}
	}
	f.outputFile, err = f.createFile()
	return err
}

func (f *File) checkIsObjExist(obj string) bool {
	if _, err := os.Stat(obj); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *File) createDir() error {
	return os.MkdirAll(f.path, os.ModePerm)
}

func (f *File) createFile() (*os.File, error) {
	return os.Create(f.path + f.name)
}

func (f *File) deleteFile() error {
	return os.Remove(f.path + f.name)
}

func getRandomTimeString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randPart := r.Intn(1000)                        // Generate a random integer
	timePart := time.Now().Format("20060102150405") // Get the current time in the format YYYYMMDDHHMMSS
	return fmt.Sprintf("%s%d", timePart, randPart)
}

// check array contains
func ArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// ReadLastNLines reads the last n lines from a file.
func ReadLastNLines(file *os.File, n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("invalid number of lines: %d", n)
	}

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return "", nil
	}

	bufferSize := 1024
	buf := make([]byte, bufferSize)
	lines := make([]string, 0, n)
	offset := int64(0)
	lineCount := 0

	for offset < fileSize && lineCount < n {
		readSize := min(bufferSize, int(fileSize-offset))
		offset += int64(readSize)

		_, err := file.Seek(-offset, io.SeekEnd)
		if err != nil {
			return "", err
		}

		_, err = file.Read(buf[:readSize])
		if err != nil {
			return "", err
		}

		// Reverse the buffer to process lines from end to start
		for i := readSize - 1; i >= 0 && lineCount < n; i-- {
			if buf[i] == '\n' || i == 0 {
				start := i
				if buf[i] == '\n' {
					start++
				}
				line := string(buf[start:readSize])
				if line != "" || i == 0 {
					lines = append([]string{line}, lines...)
					lineCount++
					readSize = i
				}
			}
		}
	}
	return strings.Join(lines, "\n"), nil
}

func MergePath(paths ...string) string {
	pathArr := make([]string, 0)
	for _, path := range paths {
		pathArr = append(pathArr, strings.Split(path, "/")...)
	}
	return strings.Join(pathArr, "/")
}
