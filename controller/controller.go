package controller

import (
	"fmt"
	"golang_drive/crawl"
	"golang_drive/create"
	"golang_drive/delete"
	"os"
	"time"

	"errors"
	"strings"

	drive "google.golang.org/api/drive/v3"
)

//because my computer sucks and it takes awhile to register folders
const timeToWait time.Duration = (time.Millisecond * 1000)

func grabFiles(driveS *drive.Service) (*drive.FileList, error) {
	gFiles, err := driveS.Files.List().PageSize(1000).Fields("files(name,id,parents,starred,mimeType,md5Checksum)").Do()
	if err != nil {
		return nil, err
	}
	return gFiles, nil
}

//gets the id of the folder in google drive where things will be uploaded
//bool is if the target is starred or not - the default is true
func getRoot(gFiles *drive.FileList, name string, starred bool) (gRoot string, err error) {

	for _, i := range gFiles.Files {
		if i.Name == name && i.Starred == starred && !i.Trashed {

			return i.Id, nil
		}
	}
	return "", errors.New("google folder: " + string(name) + " not found")

}

//removes trailing slash from target folder path
//this way the ends of files and folders are the same
func trimPath(callArgs *[]string) (trimmed string) {

	if strings.HasSuffix((*callArgs)[2], (*callArgs)[0]) {
		temp := strings.Split((*callArgs)[2], (*callArgs)[0])
		return temp[len(temp)-2]
	}

	temp := strings.Split((*callArgs)[2], (*callArgs)[0])

	return temp[len(temp)-1]
}

func localUpload(driveS *drive.Service, callArgs *[]string) error {

	localList, err := crawl.Local((*callArgs)[2])
	if err != nil {
		return err
	}

	getFiles, err := grabFiles(driveS)
	if err != nil {
		return err
	}

	driveList := make(map[string]string)
	driveRoot, err := getRoot(getFiles, (*callArgs)[3], true)
	if err != nil {
		return err
	}

	//removes trailing slash from target folder path
	//this way the ends of files and folders are the same
	trimTarget := trimPath(callArgs)

	err = delete.Drive(driveS, getFiles, trimTarget, driveRoot)
	if err != nil {
		return err
	}

	err = create.Drive(driveS, callArgs, &driveRoot, &localList, &driveList)
	if err != nil {
		return err
	}

	return nil

}

func driveDownload(driveS *drive.Service, callArgs *[]string) error {
	getFiles, err := grabFiles(driveS)
	if err != nil {
		return err
	}

	driveRoot, err := getRoot(getFiles, (*callArgs)[3], true)
	if err != nil {
		return err
	}

	driveList := make(map[string]string)
	for _, i := range getFiles.Files {

		//looks redundant, but index 0 of nil is an out of bounds slice
		if i.Parents != nil {

			if i.Parents[0] == driveRoot && i.Name == (*callArgs)[4] {

				//if call is to a folder, download the folder
				if i.MimeType == "application/vnd.google-apps.folder" {
					fmt.Println("download contents of folder:", i.Name)

					crawl.Drive(&driveList, getFiles, callArgs, i.Id, "")

					//make local target if it does not exist
					err = os.MkdirAll((*callArgs)[2], 'd')
					if err != nil {
						if os.IsExist(err) {

						} else {
							return err
						}
					}
					//because the folder can take a moment to register with the OS - at least for me
					time.Sleep(timeToWait)

					//sorts through driveList
					for k := range driveList {
						err := create.Local(driveS, callArgs, driveList[k], k)

						if err != nil {
							return err
						}

						//so folders have time to register being made
						time.Sleep(timeToWait)

					}

				} else {
					fmt.Println("downloading file:", i.Name)
					//if call is to a single file, just download that file
					err := create.Local(driveS, callArgs, i.Id, (*callArgs)[0]+i.Name)

					if err != nil {
						return err
					}
				}

				break

			}

		}
	}

	return nil
}

//StartAPI and this package control the other packages
func StartAPI(driveS *drive.Service, callArgs *[]string) error {

	switch {
	case (*callArgs)[1] == "localupload":

		err := localUpload(driveS, callArgs)
		if err != nil {
			return err
		}

	case (*callArgs)[1] == "drivedownload":
		err := driveDownload(driveS, callArgs)
		if err != nil {
			return err
		}

	}

	return nil
}
