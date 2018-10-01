package main

import (
	"flag"
	"fmt"
	"golang_drive/authorize"
	"golang_drive/controller"
	"io/ioutil"
	"os"

	"runtime"
	"strings"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

func manualInit(clientPath string, tokenPath string) (*drive.Service, error) {

	//
	// add len(os.Arg) != 1 check
	//
	//

	//sets trailing slash type depending on operating system

	//Initiates Google Drive service
	c, err := ioutil.ReadFile(clientPath)
	if err != nil {
		return nil, err
	}

	creds, err := google.ConfigFromJSON(c, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	driveClient, err := authorize.GetClient(creds, tokenPath, clientPath)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.New(driveClient)
	if err != nil {
		return nil, err
	}

	fmt.Println("service started successfully")
	return driveService, nil

}

func main() {
	slashType := ""
	if runtime.GOOS == "windows" {
		slashType = "\\"
	} else {
		slashType = "/"
	}
	////////////////////////////
	//call to get token and client info
	var tokenPath string = "token.json"
	var clientPath string = "client.json"
	////////////////////////////

	//mode
	modeF := flag.String("m", "", "use to define mode")
	//local target path\file
	localF := flag.String("l", "C:\\Users\\game_game\\go\\test\\", "local file target")
	//root google folder
	googleF := flag.String("g", "shared_golang", "Google root folder")
	//target in google folder
	googleD := flag.String("d", "", "download target in Google root folder")
	flag.Parse()

	callArgs := []string{slashType, *modeF, *localF, *googleF, *googleD}

	driveService, err := manualInit(clientPath, tokenPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//trims trailing slash from local target
	if strings.HasSuffix(callArgs[2], callArgs[0]) {
		callArgs[2] = callArgs[2][:len(callArgs[2])-1]
	}

	switch {
	case *modeF == "localupload" || *modeF == "drivedownload" && *googleD != "":

		err := controller.StartAPI(driveService, &callArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

}
