package localfiles

import (
	"fmt"
	"os"
	"path/filepath"
)

func ListTarget(localArg string) (fileList []string, listErr error) {

	//fileList = make([]string, 1)

	//Filepath walk reads all of the files in the target folder
	err := filepath.Walk(localArg, func(path string, info os.FileInfo, err error) error {

		fileList = append(fileList, path)
		//fmt.Println(fileList)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return fileList, err
	}

	fmt.Println("Local file list built successfully!")
	return fileList, err
}
