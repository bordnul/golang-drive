package delete

import (
	"errors"
	"fmt"
	"os"

	drive "google.golang.org/api/drive/v3"
)

//Drive deletes target folder or file in Google Drive
func Drive(driveS *drive.Service, getFiles *drive.FileList, goal string, parent string) error {

	for _, i := range getFiles.Files {

		if i.Parents == nil {
			continue
		}

		if i.Parents[0] == parent && !i.Trashed && i.Name == goal {
			err := driveS.Files.Delete(i.Id).Do()

			if err != nil {
				return err
			}

			return nil
		}

	}
	fmt.Println("drive target " + goal + " with parent " + parent + " not found")

	return nil

}

//Local deletes target folder or file on local disk
func Local(filePath string) error {

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("local target " + filePath + " could not be found")
		}
		return err
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		return err
	}

	return nil

}
