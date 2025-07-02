package playlist_creator

import (
	"awesomeProject/internal/playlist_creator/config"
	"fmt"
	"iter"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
)

func WritePlaylist(config *config.Config, fileData *FileData) {
	if config.ShufflePlaylist {
		rand.Shuffle(len(fileData.FilesList), func(counter1, counter2 int) {
			fileData.FilesList[counter1], fileData.FilesList[counter2] = fileData.FilesList[counter2], fileData.FilesList[counter1]
		})
	}

	println("Writing playlists")

	createFolderIfNotExists(config.OutputPath)

	var chunks iter.Seq[[]FileEntry]

	if config.SplitPlaylist {
		chunks = slices.Chunk(fileData.FilesList, config.ChunkSize)
	} else {
		chunks = slices.Chunk(fileData.FilesList, len(fileData.FilesList))
	}

	var playlistNum = 0
	for entry := range chunks {
		playlistNum++
		fmt.Printf("Writing playlist %d\n", playlistNum)
		writePlaylistFile(config, playlistNum, &entry)
	}
}

func writePlaylistFile(config *config.Config, playlistNum int, fileEntries *[]FileEntry) {
	file, err := os.Create(config.OutputPath + "/" + config.PlaylistName + "_" + fmt.Sprintf("%02d", playlistNum) + ".m3u8")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.WriteString("#EXTM3U" + "\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range *fileEntries {
		var title = ""
		if config.ReadTags {
			title = createPlaylistEntryName(entry)
		}
		if title == "" {
			title = entry.FileName[:len(entry.FileName)-len(filepath.Ext(entry.FileName))]
		}

		_, err = file.WriteString("#EXTINF:" + title + "\n")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(entry.Path + "/" + entry.FileName + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createFolderIfNotExists(pathAndFolder string) {
	var _, err = os.Stat(pathAndFolder)

	if os.IsNotExist(err) {
		err := os.Mkdir(pathAndFolder, os.FileMode(0777))
		if err != nil {
			log.Fatal("Cannot create folder: "+pathAndFolder, err)
		}
	}
}
