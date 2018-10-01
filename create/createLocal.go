package create

import (
	"fmt"
	"io"
	"os"

	drive "google.golang.org/api/drive/v3"
)

//Local creates a file locally, should be more advanced but redoing this entire thing isn't worth it to me
func Local(driveS *drive.Service, callArgs *[]string, downloadID string, downloadPath string) error {

	temp, err := driveS.Files.Get((downloadID)).Do()
	if err != nil {
		return err
	}

	if temp.MimeType == "application/vnd.google-apps.folder" {

		fmt.Println("making folder:", (*callArgs)[2]+downloadPath)

		err = os.Mkdir((*callArgs)[2]+downloadPath, 'd')

		if err != nil {
			if os.IsExist(err) {
				return nil
			} else {
				return err
			}
		}
		return nil

	} else {

		filePoint, err := driveS.Files.Get(downloadID).Download()
		if err != nil {
			return err
		}

		//checks to see if file already exists; makes it if it does not exist
		fmt.Println("downloading file:", (*callArgs)[2]+downloadPath)
		fileMake, err := os.Create((*callArgs)[2] + downloadPath)
		if err != nil {
			if os.IsExist(err) {
			} else {
				return err
			}
		}

		_, err = io.Copy(fileMake, filePoint.Body)
		if err != nil {
			return err
		}

		fileMake.Close()

		return nil
	}

}
