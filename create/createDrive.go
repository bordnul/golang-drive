package create

import (
	"fmt"
	"os"
	"strings"

	drive "google.golang.org/api/drive/v3"
)

//builds files
func driveFile(driveS *drive.Service, callArgs *[]string, parentID string, name string, fileLink *os.File) error {

	meta := drive.File{Name: name, Parents: []string{parentID}}
	_, err := driveS.Files.Create(&meta).Media(fileLink).Do()
	if err != nil {
		return err
	}
	fileLink.Close()

	return nil
}

//builds folders
func driveFolder(driveS *drive.Service, callArgs *[]string, parentID string, name string, mapName string, driveList *map[string]string) error {

	meta := drive.File{Name: name,
		MimeType: "application/vnd.google-apps.folder", Parents: []string{parentID}}
	temp, err := driveS.Files.Create(&meta).Do()
	if err != nil {
		return err
	}

	(*driveList)[mapName] = temp.Id

	return nil
}

//Drive sets up to make folders and files in the google drive
func Drive(driveS *drive.Service, callArgs *[]string, driveRoot *string, localList *[]string, driveList *map[string]string) error {

	for i := range *localList {

		//open file
		fLink, err := os.Open((*localList)[i])
		if err != nil {
			return err
		}

		//trims file location to include target folder
		tempIn := strings.LastIndex((*callArgs)[2], (*callArgs)[0])
		tempName := (*callArgs)[2][:tempIn]

		//C:\Users\game_game\go gets cut out leaving \test\level two\level two - two
		//cuts off entry text to only get part we care about
		//mapName name of the to-be map
		tempPath := strings.Replace((*localList)[i], tempName, "", 1)
		fmt.Println(tempPath)

		//finds last slash to get map name
		//map name is called from driveList to get parent id
		//parent id ex driveList[mapName]
		tempIn = strings.LastIndex(tempPath, (*callArgs)[0])
		mapName := tempPath[:tempIn]

		//gets relative file path - split counts the other side of a slash ex \foo is length two
		//folder & file length
		tempLen := len(strings.Split(tempPath, (*callArgs)[0]))

		//gets stats
		fStat, err := fLink.Stat()
		if err != nil {
			return err
		}

		switch {

		//Seems like a sitched together thing, but the slash causes issues with split
		//resulting in \foo always being two and \foo\whatever being three
		//however, there needs to be slash as this is necessary for building the internal folder tree
		case fStat.IsDir() && tempLen == 2:
			err := driveFolder(driveS, callArgs, *driveRoot, fStat.Name(), tempPath, driveList)
			if err != nil {
				return err
			}
			fLink.Close()
			continue

			//this is here in case a single file is uploaded
		case !fStat.IsDir() && tempLen == 2:
			err := driveFile(driveS, callArgs, *driveRoot, fStat.Name(), fLink)
			if err != nil {
				return err
			}
			continue

		case fStat.IsDir() && tempLen > 2:
			err := driveFolder(driveS, callArgs, (*driveList)[mapName], fStat.Name(), tempPath, driveList)
			if err != nil {
				return err
			}
			fLink.Close()
			continue

		case !fStat.IsDir() && tempLen > 2:
			err := driveFile(driveS, callArgs, (*driveList)[mapName], fStat.Name(), fLink)
			if err != nil {
				return err
			}
			continue
		}

	}

	return nil

}

/*


	switch {
	case folder:
		meta := drive.File{Name: name,
			MimeType: "application/vnd.google-apps.folder", Parents: []string{parentID}}
		temp, err := driveS.Files.Create(&meta).Do()
		if err != nil {
			return "", err
		}

		filePoint.Close()

		return temp.Id, nil

	default:

		meta := drive.File{Name: name, Parents: []string{parentID}}
		_, err := driveS.Files.Create(&meta).Media(filePoint).Do()
		if err != nil {
			return "", err
		}

		filePoint.Close()

		return "", nil
	}


*/

/*
func makeFolder(driveS *drive.Service, parentId string, name string, mapName string) (string, error) {

	meta := drive.File{Name: name,
		MimeType: mimeFolder, Parents: []string{parentId}}
	temp, err := driveS.Files.Create(&meta).Do()
	if err != nil {
		return "failed to make folder" + name + "mapName" + mapName + "\n", err
	}
	folderMap[mapName] = temp.Id
	fmt.Println("attempted map name", mapName)

	return temp.Id, nil
}
func makeFile(driveS *drive.Service, parentId string, name string, filePoint *os.File) error {

	meta := drive.File{Name: name, Parents: []string{parentId}}

	_, err := driveS.Files.Create(&meta).Media(filePoint).Do()
	if errFunc(err) {
		filePoint.Close()
		return err
	}
	filePoint.Close()
	return nil

}

*/
