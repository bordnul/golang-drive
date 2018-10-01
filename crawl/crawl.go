package crawl

import (
	"fmt"
	"os"
	"path/filepath"

	drive "google.golang.org/api/drive/v3"
)

const mimeFolder string = "application/vnd.google-apps.folder"

//Local crawls local file and builds a list
func Local(localArg string) (fileList []string, listErr error) {
	fmt.Println(localArg)

	//Filepath walk reads all of the files in the target folder
	err := filepath.Walk(localArg, func(path string, info os.FileInfo, err error) error {

		fileList = append(fileList, path)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	fmt.Println("local file list gathered successfully")
	return fileList, err
}

//Drive crawls google drive and builds a list
func Drive(folderMap *map[string]string, getFiles *drive.FileList, callArgs *[]string, startId string, parName string) {

	c := 0
	for _, i := range getFiles.Files {

		if i.Parents == nil {
			continue
		}

		if i.Parents[0] == startId && !i.Trashed {

			(*folderMap)[parName+(*callArgs)[0]+i.Name] = i.Id
			c++
			Drive(folderMap, getFiles, callArgs, i.Id, parName+(*callArgs)[0]+i.Name)

		}

	}
	if c == 0 {
		return
	}

}
