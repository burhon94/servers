package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func ListFiles(dir string) (filesList string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("can't read dir: %v", err)
		return
	}
	for _, file := range files {
		isDir := file.IsDir()
		if isDir == false {
			filesList = filesList + "\n" + file.Name()
		}
	}
	filesList = strings.TrimPrefix(filesList, "\n")
	fmt.Println(filesList)
	return filesList, err
}