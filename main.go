package main

import (
	"fmt"
	"log"
	"playlistCreator/internal/playlist_creator"
	"playlistCreator/internal/playlist_creator/config"
)

func main() {
	var configData, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	var fileData = playlist_creator.ReadFiles(configData)
	if fileData.FilesList == nil || len(fileData.FilesList) == 0 {
		fmt.Println("No files found")
		return
	}

	if configData.ListExtensions {
		fileData.ListFileExtensions()
	}

	if configData.ListFiles {
		fileData.ListFiles(configData)
	}

	err = playlist_creator.WritePlaylist(configData, fileData)
	if err != nil {
		log.Fatalf("Could not write playlists: %s", err)
	}
}
