package main

import (
	"awesomeProject/internal/playlist_creator"
	"awesomeProject/internal/playlist_creator/config"
	"fmt"
)

func main() {
	var config, err = config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	var fileData = playlist_creator.ReadFiles(config)
	if fileData.FilesList == nil || len(fileData.FilesList) == 0 {
		fmt.Println("No files found")
		return
	}

	if config.ListExtensions {
		fileData.ListFileExtensions()
	}

	if config.ListFiles {
		fileData.ListFiles(config)
	}

	playlist_creator.WritePlaylist(config, fileData)
}
