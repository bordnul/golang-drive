package driveapi

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"google.golang.org/api/drive/v3"
)

var folderMap = make(map[string]string)
var pathFiller string

const mimeFolder string = "application/vnd.google-apps.folder"

func errFunc(err error) bool {
	if err != nil {
		fmt.Printf(err.Error())
		return true
	}
	return false
}

func makeFolder(driveS *drive.Service, parentId string, name string, mapName string) (string, error) {

	meta := drive.File{Name: name,
		MimeType: mimeFolder, Parents: []string{parentId}}
	temp, err := driveS.Files.Create(&meta).Do()
	if errFunc(err) {
		return "", err
	}
	folderMap[mapName] = temp.Id
	fmt.Println(mapName)

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

func deleteLocal(filePath string) error {

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("target to delete does not exist")
		}
		return err
	}

	fmt.Println("deleting", filePath)

	err = os.RemoveAll(filePath)
	if err != nil {
		return errors.New("error removing files/folders")
	}

	fmt.Println("success")
	return nil

}

func fullLocalDownload(driveS *drive.Service, localArg string, driveT string, pathFiller string) error {

	upRoot := ""

	gFiles, err := driveS.Files.List().PageSize(1000).Fields("files(name,id,trashed,parents,starred, mimeType,md5Checksum)").Do()
	if err != nil {
		return err
	}
	//////////////////////////////////////////////////////////
	for _, i := range gFiles.Files {

		if i.Name == driveT && !i.Trashed && i.Starred {
			upRoot = i.Id
			break

		}

	}
	if upRoot == "" {
		return errors.New("google root folder not found")
	}

	t := strings.Split(localArg, pathFiller)
	test := 0
	for _, i := range gFiles.Files {

		if i.Name == t[len(t)-1] && !i.Trashed {
			upRoot = i.Id
			test++
			break

		}

	}
	if test == 0 {
		return errors.New("failed to find folder to download")
	}
	//////////////////////////////////////////////////////////

	err = deleteLocal(localArg)
	if err != nil {
		fmt.Println(err)
	}

	getFolders(gFiles, upRoot, "")

	for i := range folderMap {
		fmt.Println(i, folderMap[i])
	}

	return nil

}

func deleteRemote(driveS *drive.Service, fList *drive.FileList, deleteT string, rootT string) {

	r := ""
	for _, i := range fList.Files {
		if i.Name == rootT && i.Starred {
			r = i.Id
			break
		}
	}
	if r == "" {
		fmt.Println("Google Drive root not found")
		return
	}
	for _, i := range fList.Files {
		if i.Name == deleteT && i.Parents[0] == r && !i.Trashed {
			driveS.Files.Delete(i.Id).Do()
			r = ""
			fmt.Println("Succcessfully deleted old folder")
			break
		}
	}
	if r != "" {
		fmt.Println("Folder to-be-deleted not found.")
		return
	}

	return

}

//Start of all program to be called with info!
//func StartData(tList []string,driveS *drive.Service, *callArgs[]) error {
func StartData(tList []string, driveT string, gMode string, pathSlash string, driveS *drive.Service, localArg string) error {
	//(fileList, driveService, &callArgs)
	pathFiller = pathSlash

	fList, err := driveS.Files.List().PageSize(1000).Fields("files(id,name,parents,starred,trashed)").Do()
	if errFunc(err) {
		return err
	}

	switch gMode {
	/////////////////////
	case "fullupload":
		//deletes old folder
		fmt.Println("Fullupload started")
		t := strings.Split(tList[0], pathFiller)
		deleteRemote(driveS, fList, t[len(t)-1], driveT)

		for i := range tList {
			tList[i] = strings.Replace(tList[i], localArg, "", 1)

		}
		//uploads the new contents of the new folder
		err = uploadList(tList, driveS, driveT, pathFiller, localArg)
		if errFunc(err) {
			return err
		}
	/////////////////////
	case "update":
		fmt.Println("Update started")
		for i := range tList {
			tList[i] = strings.Replace(tList[i], localArg, "", 1)

		}
		updateFiles(driveS, driveT, localArg, tList, pathSlash)
		//needs driveservice, root target, local target, file list, path slash,
		//TODO

	case "delete":
		fmt.Println("Delete started")
		t := strings.Split(localArg, pathSlash)
		deleteRemote(driveS, fList, t[len(t)-1], driveT)

	case "fulldownload":
		fmt.Println("Fulldownload started")
		err = fullLocalDownload(driveS, localArg, driveT, pathFiller)
		if errFunc(err) {
			return err
		}
		return nil

	default:
		fmt.Println("Unknown command - driveapi.go")
	}

	return nil

}

func getFolders(gFiles *drive.FileList, startId string, parName string) {
	//fmt.Println("Building foldermap of Google Drive")

	c := 0
	for _, i := range gFiles.Files {
		//fmt.Println(i.Id)
		//fmt.Println(startId)

		if i.Parents[0] == startId && i.MimeType == mimeFolder && !i.Trashed {
			fmt.Println(parName + pathFiller + i.Name)
			folderMap[parName+pathFiller+i.Name] = i.Id
			c++
			//fmt.Println("grabbed:", i.Name)
			getFolders(gFiles, i.Id, parName+pathFiller+i.Name)
		}

	}

	if c == 0 {
		return
	}

}

func updateAFile(driveS *drive.Service, toUpId string, filePath string) error {

	fmt.Println("Updating file", filePath, "on Google Drive.")
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	meta := drive.File{}

	_, err = driveS.Files.Update(toUpId, &meta).Media(f).Do()
	if err != nil {
		return err
	}
	f.Close()

	return nil
}

func checkSum(filePath string, fileMd5 string) bool {
	fmt.Println("Checking hash for: ", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		errFunc(err)
	}
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		errFunc(err)
	}

	hashBytes := h.Sum(nil)[:16]
	sHash := hex.EncodeToString(hashBytes)
	err = f.Close()
	if err != nil {
		errFunc(err)
	}

	if fileMd5 == sHash {
		return true
	}
	fmt.Println("Hashes do not match.")

	return false
}

func updateFiles(driveS *drive.Service, rootT string, localArg string, tList []string, pathSlash string) error {

	fmt.Println("Update function started.")
	upRoot := ""

	gFiles, err := driveS.Files.List().PageSize(1000).Fields("files(name,id,parents,starred,mimeType,md5Checksum)").Do()
	if errFunc(err) {
		return err
	}

	//grabs google drive root folder and sets id
	for _, i := range gFiles.Files {
		if i.Name == rootT && i.Starred {
			upRoot = i.Id
			fmt.Println("root:", i.Id)
			break
		}
	}
	//returns if root in drive folder isn't found
	if upRoot == "" {
		fmt.Println("Root folder was not found!")
		return nil
	}

	t := strings.Split(localArg, pathFiller)
	//grabs goal folder in root folder to be updated and saves its id
	for _, i := range gFiles.Files {
		gName := t[len(t)-1]
		if i.Name == gName && i.Parents[0] == upRoot && !i.Trashed {
			fmt.Println("target:", i.Id, gName)
			folderMap[pathFiller] = i.Id
			upRoot = i.Id
			break
		}
	}

	getFolders(gFiles, upRoot, "")

	for i := range tList {
		if strings.HasPrefix(tList[i], pathFiller) {

			tV := strings.LastIndex(tList[i], pathFiller)
			tStr := tList[i]
			if tStr[:tV] == "" {
				//fmt.Println(tStr[tV:], "is in root")
				//fmt.Println(folderMap[pathFiller])
				tStr = pathFiller + tStr
				tV = strings.LastIndex(tStr, pathFiller)

			}

			if folderMap[tStr[:tV]] != "" {

				//fmt.Println(folderMap[tList[i][:tV]], "exists!")

				for _, g := range gFiles.Files {
					switch {

					//if google drive file is found via foldermap ids that has the same name and a parent that matches the foldermap id
					//check to see if the file needs to be updated, if so update it, if not don't
					case g.Name == tStr[tV+1:] && g.Parents[0] == folderMap[tStr[:tV]] && g.MimeType != mimeFolder:
						//fmt.Println(tStr[tV+1:], "belongs in", g.Parents[0], "which is", tStr[:tV])
						//fmt.Println("Checking sum")

						if checkSum(localArg+tList[i], g.Md5Checksum) {
							fmt.Println(tList[i], "does not need to be updated")
							tList[i] = "_"
							break

						} else {
							fmt.Println(tList[i], "needs to be updated")
							err = updateAFile(driveS, g.Id, localArg+tList[i])
							if err != nil {
								fmt.Println(err)
							}
							tList[i] = "_"
							break
						}

					case g.Name == tStr[tV+1:] && g.Parents[0] == folderMap[tStr[:tV]] && g.MimeType == mimeFolder:
						tList[i] = "_"
						break

					}

				}

			} else {
				fmt.Println(tList[i], "can't find home")
			}
		}
	}

	//Any files that get added to newFileList are new files that can't be updated because they are not on the Google Drive
	var newFileList []string
	for i := range tList {

		if !strings.HasPrefix(tList[i], "_") && tList[i] != "" {
			newFileList = append(newFileList, tList[i])
		}
	}

	//checks to see if there are any files that need updating, if length is one or more then there are files or folders
	if len(newFileList) != 0 {

		fmt.Println("Need to create", len(newFileList), "files.")

		sort.Strings(newFileList)

		for i := range newFileList {

			if strings.HasPrefix(newFileList[i], pathFiller) {

				tV := strings.LastIndex(newFileList[i], pathFiller)
				tStr := newFileList[i]

				////below is only if it is in root///
				if tStr[:tV] == "" {
					//fmt.Println(tStr[tV:], "is in root")
					//fmt.Println(folderMap[pathFiller])
					tStr = pathFiller + tStr
					tV = strings.LastIndex(tStr, pathFiller)

				}
				////////////////////////////////////

				if folderMap[tStr[:tV]] != "" {
					f, err := os.Open(localArg + newFileList[i])
					if err != nil {
						fmt.Println(err)
						continue
					}
					fStat, err := f.Stat()
					if err != nil {
						fmt.Println(err)
						continue
					}

					if fStat.IsDir() {

						_, err = makeFolder(driveS, folderMap[tStr[:tV]], fStat.Name(), tStr)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("Made folder:", fStat.Name(), "\n", tStr[:tV], "mapname:", tStr)
						fmt.Println(tStr + pathFiller)
						continue

					} else {
						err = makeFile(driveS, folderMap[tStr[:tV]], fStat.Name(), f)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("Made file:", fStat.Name())

						continue
					}

				}

			}

		}

	}

	return nil
}

func uploadList(fList []string, driveS *drive.Service, driveT string, pathSlash string, localArg string) error {

	fmt.Println("Upload function started.")
	upRoot := ""

	temp, _ := driveS.Files.List().PageSize(1000).Fields("files(id,name,starred)").Do()
	for _, i := range temp.Files {
		if i.Name == driveT && i.Starred {
			upRoot = i.Id
		}
	}
	for i := range fList {
		//fmt.Println(fList[i])
		tFile, err := os.Open(localArg + fList[i])
		if errFunc(err) {
			return err
		}
		//fmt.Println("problem opening starter file")
		//Gets information we need by comparingtarget path ex /blah/blah/yay/ to sub-directory paths like /blah/blah/yay/ohno/
		//and removing the extra to make /ohno/ Only things within the target directory matter
		//returns path string - short for True Path
		//tP := strings.Replace(fList[i], localArg, "", 1)

		//Gets last index of a slash - important for selecting path before file
		//returns integer - short for True Path Index
		tPI := strings.LastIndex(fList[i], pathFiller) + 1

		//splits tP by slashes get how many folders are involved, if tPS == 1 then it is in root, 2 it is in a folder within root, etc.
		//returns string array - short for True Path Split
		tPS := strings.Split(fList[i], pathFiller)
		tPS = tPS[1:] //Gets rid of the white space character from strings.Split - very annoying
		tStat, _ := tFile.Stat()

		switch {

		//Checks and makes folders not in the root target folder
		case len(tPS) != 1 && tStat.IsDir():

			//dictGoal, fullPath := pathFind(localArg, path)
			//This is the root folder, whose name we still want
			if fList[i][:tPI] == "" {
				//work on error handling
				//make this skip if the source target file is unwanted

				upRoot, err = makeFolder(driveS, upRoot, tStat.Name(), pathFiller+tStat.Name()+pathFiller)
				if errFunc(err) {
					fmt.Println("root error")
				}
				break
			}
			_, err = makeFolder(driveS, folderMap[fList[i][:tPI]], tStat.Name(), fList[i]+pathFiller)
			if errFunc(err) {
				tFile.Close()
				fmt.Println("make folder error")
			}
			tFile.Close()

		//Checks and makes files not in the root target folder
		case len(tPS) != 1 && !tStat.IsDir():
			err = makeFile(driveS, folderMap[fList[i][:tPI]], tStat.Name(), tFile)
			if errFunc(err) {
				return err
			}

		//Checks and makes folders in root target folder
		case len(tPS) == 1 && tStat.IsDir():
			_, err = makeFolder(driveS, upRoot, tStat.Name(), fList[i]+pathFiller)
			if errFunc(err) {
				tFile.Close()
				return err
			}
			tFile.Close()

		//Checks and makes files in root target folder
		case len(tPS) == 1 && !tStat.IsDir():
			err = makeFile(driveS, upRoot, tStat.Name(), tFile)
			if errFunc(err) {
				return err
			}
		}

	}
	return nil
}
