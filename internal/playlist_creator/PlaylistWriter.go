package playlist_creator

import (
	"errors"
	"fmt"
	"iter"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"playlistCreator/internal/playlist_creator/config"
	"slices"
)

func WritePlaylist(config *config.Config, fileData *FileData) error {
	if config.ShufflePlaylist {
		rand.Shuffle(len(fileData.FilesList), func(counter1, counter2 int) {
			fileData.FilesList[counter1], fileData.FilesList[counter2] = fileData.FilesList[counter2], fileData.FilesList[counter1]
		})
	}

	println("Writing playlists")

	err := createFolderIfNotExists(config.OutputPath)
	if err != nil {
		return err
	}

	var chunks iter.Seq[[]FileEntry]
	if config.SplitPlaylist {
		chunks = slices.Chunk(fileData.FilesList, config.ChunkSize)
	} else {
		chunks = slices.Chunk(fileData.FilesList, len(fileData.FilesList))
	}

	var playlistNum = 0
	for entry := range chunks {
		playlistNum++
		log.Printf("Writing playlist %d\n", playlistNum)
		var err = writePlaylistFile(config, playlistNum, &entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func writePlaylistFile(config *config.Config, playlistNum int, fileEntries *[]FileEntry) error {
	file, err := os.Create(config.OutputPath + "/" + config.PlaylistName + "_" + fmt.Sprintf("%02d", playlistNum) + ".m3u8")
	if err != nil {
		err = fmt.Errorf("could not create playlist file %s_%s.m3u8: %s", config.PlaylistName, fmt.Sprintf("%02d", playlistNum), err.Error())
		log.Println(err.Error())
		return err
	}

	defer file.Close()

	_, err = file.WriteString("#EXTM3U" + "\n")
	if err != nil {
		log.Printf("Error writing header to playlist file %s_%s.m3u8: %s", config.PlaylistName, fmt.Sprintf("%02d", playlistNum), err.Error())
		return err
	}

	for _, entry := range *fileEntries {
		title := ""
		if config.ReadTags {
			title = createPlaylistEntryName(entry)
		}
		if title == "" {
			title = entry.FileName[:len(entry.FileName)-len(filepath.Ext(entry.FileName))]
		}

		_, err = file.WriteString("#EXTINF:" + title + "\n")
		if err != nil {
			log.Printf("Error writing EXTINF to playlist file %s_%s.m3u8. Pitle: %s. Error:: %s", config.PlaylistName, title, fmt.Sprintf("%02d", playlistNum), err.Error())
			continue
		}
		_, err = file.WriteString(entry.Path + "/" + entry.FileName + "\n")
		if err != nil {
			log.Printf("Error writing path to playlist file %s_%s.m3u8. Path: %s/%s. Error: %s", config.PlaylistName, fmt.Sprintf("%02d", playlistNum), entry.Path, entry.FileName, err.Error())
		}
	}
	return nil
}

func createFolderIfNotExists(pathAndFolder string) error {
	_, err := os.Stat(pathAndFolder)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.MkdirAll(pathAndFolder, os.FileMode(0777))
	if err != nil {
		log.Printf("Error: Cannot create folder: %s, %s", pathAndFolder, err.Error())
		return err
	}
	return nil
}
