package playlist_creator

import (
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"playlistCreator/internal/playlist_creator/config"
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

		outputPath, err := os.MkdirTemp("", "PlaylistWriterTest")
		assert.Nil(t, err, "Expected err to be nil")
		assert.NotEmpty(t, outputPath, "Expected parent directory to not be empty")

		configData := config.Config{
			OutputPath:      outputPath,
			PlaylistName:    "testPlaylist",
			ShufflePlaylist: testData.shufflePlaylist,
			ReadTags:        testData.readTags,
			SplitPlaylist:   testData.splitPlaylist,
			ChunkSize:       testData.chunkSize,
		}

		fileData := FileData{
			FilesList: testData.filesList,
		}

		err = WritePlaylist(&configData, &fileData)
		assert.Nil(t, err, "Expected no error when writing playlist")

		for playlistIndex := 0; playlistIndex < len(testData.expectedPlaylistData); playlistIndex++ {
			file, err := os.Open(fmt.Sprintf("%s/%s_0%d.m3u8", outputPath, configData.PlaylistName, playlistIndex+1))
			assert.Nil(t, err, "Error opening file for writing: %s", err)

			reader := csv.NewReader(file)
			actualLines, err := reader.ReadAll()
			assert.Nil(t, err, "Expected no error when writing playlist")

			expectedPlaylistLines := testData.expectedPlaylistData[playlistIndex]

			// 1 header line; each entry takes 2 lines
			expectedNumberOfLines := 1 + len(expectedPlaylistLines)
			assert.Equal(t, expectedNumberOfLines, len(actualLines), "Expected number of lines in file to be %d", expectedNumberOfLines)
			assert.Equal(t, "#EXTM3U", strings.Join(actualLines[0], ""))

			if !testData.shufflePlaylist {
				for lineIndex := 0; lineIndex < expectedNumberOfLines-1; lineIndex += 2 {
					assert.Equal(t, expectedPlaylistLines[lineIndex], strings.Join(actualLines[lineIndex+1], ""), "Expected playlist to contain '%s'", expectedPlaylistLines[lineIndex])
					assert.Equal(t, expectedPlaylistLines[lineIndex+1], strings.Join(actualLines[lineIndex+2], ""), "Expected playlist to contain '%s'", expectedPlaylistLines[lineIndex+1])
				}
			} else {
				for lineIndex := 0; lineIndex < expectedNumberOfLines-1; lineIndex += 2 {
					assert.Contains(t, expectedPlaylistLines, strings.Join(actualLines[lineIndex+1], ""), "Expected to find a match for '%s'", expectedPlaylistLines[lineIndex])
					assert.Contains(t, expectedPlaylistLines, strings.Join(actualLines[lineIndex+2], ""), "Expected to find a match for '%s'", expectedPlaylistLines[lineIndex+1])
				}
			}
		}
	}
}

func TestPlaylistWriter_writePlaylistFile_ShouldWriteAPlaylistFile(t *testing.T) {
	t.Parallel()

	outputPath, err := os.MkdirTemp("", "PlaylistWriterTest")
	assert.Nil(t, err, "Expected err to be nil")
	assert.NotEmpty(t, outputPath, "Expected parent directory to not be empty")

	configData := config.Config{
		OutputPath:    outputPath,
		PlaylistName:  "testPlaylist",
		ReadTags:      false,
		SplitPlaylist: true,
		ChunkSize:     2,
	}

	fileEntries := []FileEntry{
		{"pathOne", "fileOne.mp3"},
		{"pathOne", "fileTwo.flac"},
		{"pathTwo", "fileOne.qt"},
		{"pathTwo", "fileTwo.a"},
	}

	err = writePlaylistFile(&configData, 1, &fileEntries)
	assert.Nil(t, err, "Expected no error when writing playlist")

	file, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	assert.Nil(t, err, "Expected playlist file to be created")

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	assert.Nil(t, err, "Expected error to be nil")

	// 1 header, 4 entries *2 = 8 + 1 = 9
	assert.Equal(t, 9, len(lines), "Expected number of lines in file to be 9")
	assert.Equal(t, "#EXTM3U", strings.Join(lines[0], ""), "Expected first line to be '#EXTM3U'")
	assert.Equal(t, "#EXTINF:fileOne", strings.Join(lines[1], ""), "Expected second line to be '#EXTINF:fileOne'")
	assert.Equal(t, "pathOne/fileOne.mp3", strings.Join(lines[2], ""), "Expected third line to be 'pathOne/fileOne.mp3'")
	assert.Equal(t, "#EXTINF:fileTwo", strings.Join(lines[3], ""), "Expected second line to be '#EXTINF:fileTwo'")
	assert.Equal(t, "pathOne/fileTwo.flac", strings.Join(lines[4], ""), "Expected third line to be 'pathOne/fileTwo.flac'")
	assert.Equal(t, "#EXTINF:fileOne", strings.Join(lines[5], ""), "Expected second line to be '#EXTINF:fileOne'")
	assert.Equal(t, "pathTwo/fileOne.qt", strings.Join(lines[6], ""), "Expected third line to be 'pathTwo/fileOne.qt'")
	assert.Equal(t, "#EXTINF:fileTwo", strings.Join(lines[7], ""), "Expected second line to be '#EXTINF:fileTwo'")
	assert.Equal(t, "pathTwo/fileTwo.a", strings.Join(lines[8], ""), "Expected third line to be 'pathTwo/fileTwo.a'")
}

func TestPlaylistWriter_writePlaylistFile_ShouldWriteAnEmptyPlaylistFile(t *testing.T) {
	t.Parallel()

	outputPath, err := os.MkdirTemp("", "PlaylistWriterTest")
	assert.NoError(t, err, "Expected err to be nil")
	assert.NotEmpty(t, outputPath, "Expected parent directory to not be empty")

	configData := config.Config{
		OutputPath:   outputPath,
		PlaylistName: "testPlaylist",
		ReadTags:     false,
	}

	var fileEntries []FileEntry

	err = writePlaylistFile(&configData, 1, &fileEntries)
	assert.Nil(t, err, "Expected no error when writing playlist")

	file, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	assert.NotNil(t, file, "Expected file to not be nil", err)

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	assert.Nil(t, err, "Expected error to be nil")

	assert.Equal(t, 1, len(lines), "Expected number of lines in file to be 1")
	assert.Equal(t, "#EXTM3U", strings.Join(lines[0], ""), "Expected first line to be '#EXTM3U'")
}

func TestPlaylistWriter_createFolderIfNotExists_ShouldCreateAFolderIfItDoesntExist(t *testing.T) {
	t.Parallel()

	parentDirectory, err := os.MkdirTemp("", "PlaylistWriterTest")
	assert.Nil(t, err, "Expected err to be nil")
	assert.NotNil(t, parentDirectory, "Expected parent directory to not be nil")

	_, err = os.Stat(parentDirectory)
	assert.False(t, os.IsNotExist(err), "Expected parent directory to be created")

	testPath := parentDirectory + "/testPlaylistFile"
	err = os.MkdirAll(testPath, os.FileMode(0777))
	assert.Nil(t, err, "Error: Cannot create folder: %s, %s", testPath, err)

	_, err = os.Stat(testPath)
	assert.False(t, os.IsNotExist(err), "Expected directory to be created")
}
