package playlist_creator

import (
	"awesomeProject/internal/playlist_creator/config"
	"fmt"
	"strings"
)

type FileEntry struct {
	Path     string
	FileName string
}

type FileData struct {
	FilesList        []FileEntry
	UniqueExtensions map[string]int
}

func (fileData *FileData) addFileExtension(extension string) {
	var extensionLowercase = strings.ToLower(extension)
	if fileData.UniqueExtensions == nil { //TODO: Can this initialisation be moved to a getter?
		fileData.UniqueExtensions = make(map[string]int)
	}

	if _, ok := fileData.UniqueExtensions[extensionLowercase]; ok {
		fileData.UniqueExtensions[extensionLowercase] += 1
	} else {
		fileData.UniqueExtensions[extensionLowercase] = 1
	}
}

func (fileData *FileData) ListFileExtensions() {
	fmt.Println("Extensions:")
	for extension, count := range fileData.UniqueExtensions {
		fmt.Printf("%s: %d\n", extension, count)
	}
}

func (fileData *FileData) ListFiles(config *config.Config) {
	var listLimit = len(fileData.FilesList)
	if config.ListLimit != -1 {
		listLimit = min(config.ListLimit, len(fileData.FilesList))
		fmt.Printf("First %d files in list of length %d:\n", listLimit, len(fileData.FilesList))
	} else {
		fmt.Printf("Files in list of length %d:\n", len(fileData.FilesList))
	}

	for counter := 0; counter < listLimit; counter++ {
		var fileEntry = fileData.FilesList[counter]
		fmt.Printf("%s - %s\n", fileEntry.Path, fileEntry.FileName)
	}
	if listLimit < len(fileData.FilesList) {
		fmt.Println("...")
	}
}
