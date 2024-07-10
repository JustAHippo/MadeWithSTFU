package main

import (
	"context"
	"encoding/hex"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Folder picker
func (a *App) FolderSelector(name string) string {
	result, err := cfdutil.ShowPickFolderDialog(cfd.DialogConfig{
		Title:  "Pick Folder",
		Role:   "FolderSelector",
		Folder: "C:\\",
	})
	if err == cfd.ErrorCancelled {
		println("Dialog was cancelled by the user.")
	} else if err != nil {
		log.Fatal(err)
	}

	return result
}

// Runs the actual process for patching
func (a *App) RunDeletion(path string, version string) string {
	warningMessages := ""
	dataDir := ""
	//Read the project directory for all folders
	entries, err := os.ReadDir(path)
	if err != nil {
		return "Couldn't read directory"
	}
	//Find the data folder within the others
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "_Data") {
			dataDir = path + "\\" + entry.Name()
		}
	}
	//Couldn't find any folder that has _Data
	if dataDir == "" {
		return "Couldn't find the _Data directory"
	}
	//globalgamemanagers contains information about pro version and splash screen
	gameManagerPath := dataDir + "/globalgamemanagers"
	managerContent, err := ioutil.ReadFile(gameManagerPath)
	// if no file could be found!
	if err != nil {
		return "Couldn't find the globalgamemanagers file"
	}
	//Read bytes and make hex for game manager
	hexManagerContent := hex.EncodeToString(managerContent)
	hexVersion := hex.EncodeToString([]byte(version))
	versionIndex := strings.LastIndex(hexManagerContent, hexVersion)
	if versionIndex == -1 {
		return "The version entered is likely incorrect. It could not be found in the file."
	}
	//20th byte before version has a boolean hasProVersion. We set it to 01.
	hexManagerContent = replaceAtIndex(hexManagerContent, "0", versionIndex-40)
	hexManagerContent = replaceAtIndex(hexManagerContent, "1", versionIndex-39)

	//This hex is for showSplashScreen. It is semi inconsistent in placement, but always has the same bytes leading to it
	//The byte following 3F at the end, 01, is the boolean we change.
	searchForHex := "8D8C0C3EF9F8F83D8180003E0000803F0101"
	hexReplacement := "8D8C0C3EF9F8F83D8180003E0000803F0001"
	//Golang parses in lowercase lol
	searchForHex = strings.ToLower(searchForHex)
	hexReplacement = strings.ToLower(hexReplacement)
	//Make sure one of them is present.
	if strings.Index(hexManagerContent, searchForHex) == -1 {
		if strings.Index(hexManagerContent, hexReplacement) != -1 {
			warningMessages += "WARNING: showSplashScreen was already set to false.\n"
		} else {
			return "Failed to find showSplashScreen pattern"
		}
	}
	//Replace the string with the patched one
	hexManagerContent = strings.ReplaceAll(hexManagerContent, searchForHex, hexReplacement)
	if strings.Index(hexManagerContent, searchForHex) != -1 {
		return "Failed to find showSplashScreen pattern... Didn't replace correctly."
	}

	err = os.Rename(gameManagerPath, gameManagerPath+".bak")
	if err != nil {
		return "Couldn't back up old game manager"
	}
	fileContent, _ := hex.DecodeString(hexManagerContent)
	file, _ := os.Create(gameManagerPath)
	_, err = file.Write(fileContent)
	if err != nil {
		return "Didn't write new gamemanager successfully?"
	}
	return warningMessages + "Succeeded in " + dataDir + " on version " + version
}

func replaceAtIndex(input string, replacement string, index int) string {
	out := []rune(input)

	out[index] = []rune(replacement)[0]
	return string(out)
}
