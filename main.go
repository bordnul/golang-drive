package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"Mine/GoDrive/authorize"
	"Mine/GoDrive/driveapi"
	"Mine/GoDrive/localfiles"

	"runtime"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

//var fileList []string

func errFunc(err error) {
	if err != nil {
		fmt.Printf("Error!!: %v\nExiting program.", err.Error())
		os.Exit(404)
	}
}

func errFile(err error, name string) {

	fmt.Println(err, "encountered while uploading: ", name, "\nFile", name, " skipped.")

}

func manualInit(clientPath string, tokenPath string) (*drive.Service, error) {

	//
	// add len(os.Arg) != 1 check
	//
	//

	//sets trailing slash type depending on operating system

	//Initiates Google Drive service
	c, err := ioutil.ReadFile(clientPath)
	errFunc(err)

	creds, err := google.ConfigFromJSON(c, drive.DriveScope)
	errFunc(err)

	driveClient, err := authorize.GetClient(creds, tokenPath, clientPath)
	errFunc(err)

	driveService, err := drive.New(driveClient)
	errFunc(err)

	fmt.Printf("Service started successfully.\n")
	////////////////////////////////
	return driveService, nil

}

func commandFind() {

	testCommand := make(map[string]string)
	testCommand["upload"] = "/tmp/fake/upload"
	testCommand["download"] = "/tmp/fake/download"

	return

}

func main() {
	var folderArg string = "shared_golang"
	var modeArg string = "fulldownload"
	var localArg string = "C:\\Users\\game_game\\go\\test\\level two\\level three"
	//var folderMap = make(map[string]string)
	var pathFiller string
	////////////////////////////
	if runtime.GOOS == "windows" {
		pathFiller = "\\"
	} else {
		pathFiller = "/"
	}
	//callArgs 0 = pathFiller, 1 = folderArg, 2 = localArg, 3 = modeArg
	callArgs := []string{pathFiller, folderArg, localArg, modeArg}

	////////////////////////////
	//call to get token and client info
	var tokenPath string = "token.json"
	var clientPath string = "client.json"
	////////////////////////////

	driveService, err := manualInit(clientPath, tokenPath)
	if err != nil {
		errFunc(err)
	}

	switch modeArg {
	case "fullupload":

		fileList, err := localfiles.ListTarget(callArgs[2])
		if err != nil {
			errFunc(err)
		}
		fmt.Println("file list gathered successfully")

		//err = driveapi.StartData(fileList, driveService, &callArgs)
		err = driveapi.StartData(fileList, folderArg, modeArg, pathFiller, driveService, localArg)
		if err != nil {
			errFunc(err)
		}

	case "update":

		fileList, err := localfiles.ListTarget(localArg)
		if err != nil {
			errFunc(err)
		}

		err = driveapi.StartData(fileList, folderArg, modeArg, pathFiller, driveService, localArg)
		if err != nil {
			errFunc(err)
		}

	case "delete":

	case "fulldownload":
		err := driveapi.StartData([]string{}, folderArg, modeArg, pathFiller, driveService, localArg)
		if err != nil {
			errFunc(err)
		}

	case "downloadUpdate":

	default:
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		bleh, err := driveService.Files.List().PageSize(1000).Fields("files(name,id)").Do()

		k := ""

		if err != nil {
			errFunc(err)
		}
		for _, i := range bleh.Files {
			if i.Name == "Print this" {
				k = i.Id
				break
			}
		}
		fmt.Println(k)

		//fmt.Println(k)
		/*
				ey := driveService.Files.Get(k)
				file, err := ey.Download()
				if err != nil {
					errFunc(err)
				}
				fmt.Println(file)

				fmt.Println("Attempting download now")

			err = os.MkdirAll("C:\\Users\\game_game\\go\\src\\Mine\\Testing\\coocoo\\", 'd')
			if err != nil {
				errFunc(err)
			}
		*/

		/*
			_, err = io.Copy(make, file.Body)
			if err != nil {
				errFunc(err)
			}
			make.Close()
			//file.Close()
		*/
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////

	commandFind()

}
