package file_storage

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateFileInLocalStorage(buf []byte, fileExt string) (url string, err error) {
	fileId, _ := uuid.NewRandom()
	filePath := fmt.Sprintf("static/%s/%s%s", string(fileId.String()[0]), fileId.String(), fileExt)
	path := fmt.Sprintf("static/%s", string(fileId.String()[0]))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	err = ioutil.WriteFile(filePath, buf, os.ModePerm)
	if err != nil {
		return
	}

	return "/" + filePath, nil
}

func LoadFile(name string, folder string) (data []byte, ext string) {
	data, err := ioutil.ReadFile("static/" + folder + "/" + name)
	if err != nil {
		return []byte{}, ""
	}
	return data, filepath.Ext(name)
}