package playlist_creator

import (
	"log"
	"os"
	"path/filepath"
	"playlistCreator/internal/playlist_creator/config"
	"sort"
	"strings"
)

func ReadFiles(config *config.Config) (error, *FileData) {
	var fileData = NewFileData()

	err := readFilesRecursively(config, config.InputPath, fileData)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	if fileData.FilesList != nil && len(fileData.FilesList) > 1 {
		sort.Slice(fileData.FilesList, func(i, j int) bool {
			return fileData.FilesList[i].FileName < fileData.FilesList[j].FileName
		})
	}

	return nil, fileData
}

func readFilesRecursively(config *config.Config, directoryNameAndPath string, fileData *FileData) error {
	files, err := os.ReadDir(directoryNameAndPath)
	if err != nil {
		log.Printf("Warning: Could not read directory %s: %s", directoryNameAndPath, err.Error())
		return err
	}

	// https://golang.cafe/blog/how-to-list-files-in-a-directory-in-go.html
	// filepath.Walk
	for _, file := range files {
		if !file.IsDir() {
			extension := filepath.Ext(file.Name())
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
	return nil
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
