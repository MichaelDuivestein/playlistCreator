package playlist_creator

import (
	"awesomeProject/internal/playlist_creator/config"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ReadFiles(config *config.Config) *FileData {
	var fileData = FileData{}

	readFilesRecursively(config, config.InputPath, &fileData)

	sort.Slice(fileData.FilesList, func(i, j int) bool {
		return fileData.FilesList[i].FileName < fileData.FilesList[j].FileName
	})

	return &fileData
}

func readFilesRecursively(config *config.Config, directoryNameAndPath string, fileData *FileData) {
	files, err := os.ReadDir(directoryNameAndPath)
	if err != nil {
		log.Fatal(err)
	}

	// https://golang.cafe/blog/how-to-list-files-in-a-directory-in-go.html
	// filepath.Walk
	for _, file := range files {
		if !file.IsDir() {
			var extension = filepath.Ext(file.Name())
			if extension != "" {
				extension = extension[1:]

				if isExtensionAllowed(config, extension) {
					fileData.addFileExtension(extension)
				}

				var fileEntry = FileEntry{directoryNameAndPath, file.Name()}
				fileData.FilesList = append(fileData.FilesList, fileEntry)
			}
		} else {
			readFilesRecursively(config, directoryNameAndPath+"/"+file.Name(), fileData)
		}
	}
}

func isExtensionAllowed(config *config.Config, extension string) bool {
	// if there's no whitelist, enable all extensions
	if len(config.ExtensionWhitelist) == 0 {
		return true
	}

	for _, ext := range config.ExtensionWhitelist {
		if strings.EqualFold(ext, extension) {
			return true
		}
	}

	return false
}
