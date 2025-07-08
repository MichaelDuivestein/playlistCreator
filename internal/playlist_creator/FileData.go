package playlist_creator

import (
	"log"
	"playlistCreator/internal/playlist_creator/config"
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

func NewFileData() *FileData { return &FileData{UniqueExtensions: make(map[string]int)} }

func (fileData *FileData) addFileExtension(extension string) {
	extensionLowercase := strings.ToLower(extension)
	if fileData.UniqueExtensions == nil {
		fileData.UniqueExtensions = make(map[string]int)
	}

	if _, ok := fileData.UniqueExtensions[extensionLowercase]; ok {
		fileData.UniqueExtensions[extensionLowercase] += 1
	} else {
		fileData.UniqueExtensions[extensionLowercase] = 1
	}
}

func (fileData *FileData) ListFileExtensions() {
	log.Println("Extensions:")
	for extension, count := range fileData.UniqueExtensions {
		log.Printf("%s: %d\n", extension, count)
	}
}

func (fileData *FileData) ListFiles(config *config.Config) {
	listLimit := len(fileData.FilesList)
	if config.ListLimit != -1 {
		listLimit = min(config.ListLimit, len(fileData.FilesList))
		log.Printf("First %d files in list of length %d:\n", listLimit, len(fileData.FilesList))
	} else {
		log.Printf("Files in list of length %d:\n", len(fileData.FilesList))
	}

	for counter := 0; counter < listLimit; counter++ {
		fileEntry := fileData.FilesList[counter]
		log.Printf("%s - %s\n", fileEntry.Path, fileEntry.FileName)
	}
	if listLimit < len(fileData.FilesList) {
		log.Println("...")
	}
}
