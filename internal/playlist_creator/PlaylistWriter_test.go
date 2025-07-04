package playlist_creator

import (
	"encoding/csv"
	"fmt"
	"os"
	"playlistCreator/internal/playlist_creator/config"
	"slices"
	"strings"
	"testing"
)

func TestPlaylistWriter_writePlaylist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		shufflePlaylist      bool
		readTags             bool
		filesList            []FileEntry
		splitPlaylist        bool
		chunkSize            int
		expectedPlaylistData [][]string
	}{
		{
			name:            "Should write a playlist",
			shufflePlaylist: false,
			readTags:        false,
			filesList: []FileEntry{
				{"folderOne", "fileOne.mp3"},
				{"folderOne", "fileTwo.flac"},
				{"folderTwo", "fileThree.qt"},
				{"folderTwo", "fileFour.a"},
			},
			expectedPlaylistData: [][]string{
				{"#EXTINF:fileOne", "folderOne/fileOne.mp3", "#EXTINF:fileTwo", "folderOne/fileTwo.flac", "#EXTINF:fileThree", "folderTwo/fileThree.qt", "#EXTINF:fileFour", "folderTwo/fileFour.a"},
			},
		},
		{
			name:            "Should write a shuffled playlist",
			shufflePlaylist: true,
			readTags:        false,
			filesList: []FileEntry{
				{"folderOne", "fileOne.mp3"},
				{"folderOne", "fileTwo.flac"},
				{"folderTwo", "fileThree.qt"},
				{"folderTwo", "fileFour.a"},
			},
			expectedPlaylistData: [][]string{
				{"#EXTINF:fileOne", "folderOne/fileOne.mp3", "#EXTINF:fileTwo", "folderOne/fileTwo.flac", "#EXTINF:fileThree", "folderTwo/fileThree.qt", "#EXTINF:fileFour", "folderTwo/fileFour.a"},
			},
		},
		{
			name:          "Should split a playlist when entries exceed split size",
			readTags:      false,
			splitPlaylist: true,
			chunkSize:     3,
			filesList: []FileEntry{
				{"folderOne", "fileOne.mp3"},
				{"folderOne", "fileTwo.flac"},
				{"folderTwo", "fileThree.qt"},
				{"folderTwo", "fileFour.a"},
			},
			expectedPlaylistData: [][]string{
				{"#EXTINF:fileOne", "folderOne/fileOne.mp3", "#EXTINF:fileTwo", "folderOne/fileTwo.flac", "#EXTINF:fileThree", "folderTwo/fileThree.qt"},
				{"#EXTINF:fileFour", "folderTwo/fileFour.a"},
			},
		},
		{
			name:          "Should not split a playlist when entries equal split size",
			readTags:      false,
			splitPlaylist: true,
			chunkSize:     4,
			filesList: []FileEntry{
				{"folderOne", "fileOne.mp3"},
				{"folderOne", "fileTwo.flac"},
				{"folderTwo", "fileThree.qt"},
				{"folderTwo", "fileFour.a"},
			},
			expectedPlaylistData: [][]string{
				{"#EXTINF:fileOne", "folderOne/fileOne.mp3", "#EXTINF:fileTwo", "folderOne/fileTwo.flac", "#EXTINF:fileThree", "folderTwo/fileThree.qt", "#EXTINF:fileFour", "folderTwo/fileFour.a"},
			},
		},
		{
			name:          "Should not split a playlist when entries exceed split size splitPlaylist is false",
			readTags:      false,
			splitPlaylist: false,
			chunkSize:     2,
			filesList: []FileEntry{
				{"folderOne", "fileOne.mp3"},
				{"folderOne", "fileTwo.flac"},
				{"folderTwo", "fileThree.qt"},
				{"folderTwo", "fileFour.a"},
			},
			expectedPlaylistData: [][]string{
				{"#EXTINF:fileOne", "folderOne/fileOne.mp3", "#EXTINF:fileTwo", "folderOne/fileTwo.flac", "#EXTINF:fileThree", "folderTwo/fileThree.qt", "#EXTINF:fileFour", "folderTwo/fileFour.a"},
			},
		},
	}

	for _, testData := range tests {
		testData := testData
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()
		})

		var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
		if err != nil {
			t.Error("Expected err to be nil", err)
		}
		if outputPath == "" {
			t.Error("Expected parent directory to not be empty")
		}

		var configData = config.Config{
			OutputPath:      outputPath,
			PlaylistName:    "testPlaylist",
			ShufflePlaylist: testData.shufflePlaylist,
			ReadTags:        testData.readTags,
			SplitPlaylist:   testData.splitPlaylist,
			ChunkSize:       testData.chunkSize,
		}

		var fileData = FileData{
			FilesList: testData.filesList,
		}

		err = WritePlaylist(&configData, &fileData)
		if err != nil {
			t.Error("Expected no error when writing playlist", err)
		}

		for playlistIndex := 0; playlistIndex < len(testData.expectedPlaylistData); playlistIndex++ {
			file, err := os.Open(fmt.Sprintf("%s/%s_0%d.m3u8", outputPath, configData.PlaylistName, playlistIndex+1))
			if err != nil {
				t.Errorf("Error opening file for writing: %s", err)
			}

			reader := csv.NewReader(file)

			actualLines, err := reader.ReadAll()
			if err != nil {
				t.Error("Expected no error when writing playlist", err)
			}

			var expectedPlaylistLines = testData.expectedPlaylistData[playlistIndex]

			// 1 header line; each entry takes 2 lines
			var expectedNumberOfLines = 1 + len(expectedPlaylistLines)
			if len(actualLines) != expectedNumberOfLines {
				t.Errorf("Expected number of lines in file to be %d", expectedNumberOfLines)
			}

			if strings.Join(actualLines[0], "") != "#EXTM3U" {
				t.Errorf("Expected first line to be '#EXTM3U'")
			}

			if !testData.shufflePlaylist {
				for lineIndex := 0; lineIndex < expectedNumberOfLines-1; lineIndex += 2 {
					if strings.Join(actualLines[lineIndex+1], "") != expectedPlaylistLines[lineIndex] {
						t.Errorf("Expected playlist to contain '%s'", expectedPlaylistLines[lineIndex])
					}
					if strings.Join(actualLines[lineIndex+2], "") != expectedPlaylistLines[lineIndex+1] {
						t.Errorf("Expected playlist to contain '%s'", expectedPlaylistLines[lineIndex+1])
					}
				}
			} else {
				for lineIndex := 0; lineIndex < expectedNumberOfLines-1; lineIndex += 2 {
					if !slices.Contains(expectedPlaylistLines, strings.Join(actualLines[lineIndex+1], "")) {
						t.Errorf("Expected to find a match for '%s'", expectedPlaylistLines[lineIndex])
					}
					if !slices.Contains(expectedPlaylistLines, strings.Join(actualLines[lineIndex+2], "")) {
						t.Errorf("Expected to find a match for '%s'", expectedPlaylistLines[lineIndex+1])
					}
				}
			}
		}
	}
}

func TestPlaylistWriter_writePlaylistFile_ShouldWriteAPlaylistFile(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:    outputPath,
		PlaylistName:  "testPlaylist",
		ReadTags:      false,
		SplitPlaylist: true,
		ChunkSize:     2,
	}

	var fileEntries = []FileEntry{
		{"pathOne", "fileOne.mp3"},
		{"pathOne", "fileTwo.flac"},
		{"pathTwo", "fileOne.qt"},
		{"pathTwo", "fileTwo.a"},
	}

	err = writePlaylistFile(&configData, 1, &fileEntries)
	if err != nil {
		t.Fatal("Expected no error when writing playlist", err)
	}

	file, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	if err != nil {
		t.Fatal("Expected playlist file to be created", err)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		t.Fatal("Expected error to be nil", err)
	}

	// 1 header, 4 entries *2 = 8 + 1 = 9
	if len(lines) != 9 {
		t.Fatal("Expected number of lines in file to be 9")
	}

	if strings.Join(lines[0], "") != "#EXTM3U" {
		t.Fatal("Expected first line to be '#EXTM3U'")
	}

	if strings.Join(lines[1], "") != "#EXTINF:fileOne" {
		t.Fatal("Expected second line to be '#EXTINF:fileOne'")
	}
	if strings.Join(lines[2], "") != "pathOne/fileOne.mp3" {
		t.Fatal("Expected third line to be 'pathOne/fileOne.mp3'")
	}

	if strings.Join(lines[3], "") != "#EXTINF:fileTwo" {
		t.Fatal("Expected fourth line to be '#EXTINF:fileTwo'")
	}
	if strings.Join(lines[4], "") != "pathOne/fileTwo.flac" {
		t.Fatal("Expected fifth line to be 'pathOne/fileTwo.flac'")
	}

	if strings.Join(lines[5], "") != "#EXTINF:fileOne" {
		t.Fatal("Expected sixth line to be '#EXTINF:fileOne'")
	}
	if strings.Join(lines[6], "") != "pathTwo/fileOne.qt" {
		t.Fatal("Expected seventh line to be 'pathTwo/fileOne.qt'")
	}

	if strings.Join(lines[7], "") != "#EXTINF:fileTwo" {
		t.Fatal("Expected eighth line to be '#EXTINF:fileTwo'")
	}
	if strings.Join(lines[8], "") != "pathTwo/fileTwo.a" {
		t.Fatal("Expected fourth line to be 'pathTwo/fileTwo.a'")
	}
}

func TestPlaylistWriter_writePlaylistFile_ShouldWriteAnEmptyPlaylistFile(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:   outputPath,
		PlaylistName: "testPlaylist",
		ReadTags:     false,
	}

	var fileEntries []FileEntry

	err = writePlaylistFile(&configData, 1, &fileEntries)
	if err != nil {
		t.Fatal("Expected no error when writing playlist", err)
	}

	file, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	if err != nil {
		t.Fatal("Expected playlist file to be created", err)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		t.Fatal("Expected error to be nil", err)
	}

	if len(lines) != 1 {
		t.Fatal("Expected number of lines in file to be 1")
	}

	if strings.Join(lines[0], "") != "#EXTM3U" {
		t.Fatal("Expected first line to be '#EXTM3U'")
	}
}

func TestPlaylistWriter_createFolderIfNotExists_ShouldCreateAFolderIfItDoesntExist(t *testing.T) {
	var parentDirectory, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if parentDirectory == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	_, err = os.Stat(parentDirectory)
	if os.IsNotExist(err) {
		t.Fatal("Expected parent directory to be created")
	}

	var testPath = parentDirectory + "/testPlaylistFile"
	err = createFolderIfNotExists(testPath)
	if err != nil {
		t.Fatal("Expected no error when creating folder", err)
	}

	_, err = os.Stat(testPath)
	if os.IsNotExist(err) {
		t.Fatal("Expected directory to be created")
	}
}
