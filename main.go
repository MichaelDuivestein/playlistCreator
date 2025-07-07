package main

import (
	"log"
	"playlistCreator/internal/playlist_creator"
	"playlistCreator/internal/playlist_creator/config"
)

func main() {
	var configData, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	err, fileData := playlist_creator.ReadFiles(configData)
	if err == nil {
		log.Println("Error while reading files")
		return
	}
	if fileData.FilesList == nil || len(fileData.FilesList) == 0 {
		log.Println("No files found")
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
