package playlist_creator

import (
	"encoding/csv"
	"os"
	"playlistCreator/internal/playlist_creator/config"
	"slices"
	"strings"
	"testing"
)

func TestPlaylistWriter_writePlaylist_ShouldWriteAPlaylist(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:      outputPath,
		PlaylistName:    "testPlaylist",
		ShufflePlaylist: false,
		ReadTags:        false,
	}

	var fileData = FileData{
		FilesList: []FileEntry{
			{"pathOne", "fileOne.mp3"},
			{"pathOne", "fileTwo.flac"},
			{"pathTwo", "fileOne.qt"},
			{"pathTwo", "fileTwo.a"},
		},
	}

	err = WritePlaylist(&configData, &fileData)
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

func TestPlaylistWriter_writePlaylist_ShouldWriteAShuffledPlaylist(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:      outputPath,
		PlaylistName:    "testPlaylist",
		ShufflePlaylist: true,
		ReadTags:        false,
	}

	var fileData = FileData{
		FilesList: []FileEntry{
			{"pathOne", "fileOne.mp3"},
			{"pathOne", "fileTwo.flac"},
			{"pathTwo", "fileOne.qt"},
			{"pathTwo", "fileTwo.a"},
		},
	}

	err = WritePlaylist(&configData, &fileData)
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

	var expectedLines = []string{
		"#EXTM3U",
		"#EXTINF:fileOne",
		"pathOne/fileOne.mp3",
		"#EXTINF:fileTwo",
		"pathOne/fileTwo.flac",
		"#EXTINF:fileOne",
		"pathTwo/fileOne.qt",
		"#EXTINF:fileTwo",
		"pathTwo/fileTwo.a",
	}

	for index := range lines {
		if !slices.Contains(expectedLines, strings.Join(lines[index], "")) {
			t.Fatal("Expected to find a match for '" + strings.Join(lines[index], "") + "'")
		}
	}
}

func TestPlaylistWriter_writePlaylist_ShouldSplitAPaylistWhenEntriesExceedSplitSize(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:      outputPath,
		PlaylistName:    "testPlaylist",
		ShufflePlaylist: false,
		ReadTags:        false,
		SplitPlaylist:   true,
		ChunkSize:       3,
	}

	var fileData = FileData{
		FilesList: []FileEntry{
			{"pathOne", "fileOne.mp3"},
			{"pathOne", "fileTwo.flac"},
			{"pathTwo", "fileOne.qt"},
			{"pathTwo", "fileTwo.a"},
		},
	}

	err = WritePlaylist(&configData, &fileData)
	if err != nil {
		t.Fatal("Expected no error when writing playlist", err)
	}

	file1, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	reader := csv.NewReader(file1)
	lines, err := reader.ReadAll()

	// 1 header, 3 entries *2 = 8 + 1 = 7
	if len(lines) != 7 {
		t.Fatal("Expected number of lines in file to be 7")
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

	file2, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_02.m3u8")
	reader = csv.NewReader(file2)
	lines, err = reader.ReadAll()

	// 1 header, 1 entries *2 = 2 + 1 = 3
	if len(lines) != 3 {
		t.Fatal("Expected number of lines in file to be 3")
	}

	if strings.Join(lines[0], "") != "#EXTM3U" {
		t.Fatal("Expected first line to be '#EXTM3U'")
	}

	if strings.Join(lines[1], "") != "#EXTINF:fileTwo" {
		t.Fatal("Expected eighth line to be '#EXTINF:fileTwo'")
	}
	if strings.Join(lines[2], "") != "pathTwo/fileTwo.a" {
		t.Fatal("Expected fourth line to be 'pathTwo/fileTwo.a'")
	}
}

func TestPlaylistWriter_writePlaylist_ShouldNotSplitAPaylistWhenEntriesEqualsSplitSize(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:      outputPath,
		PlaylistName:    "testPlaylist",
		ShufflePlaylist: false,
		ReadTags:        false,
		SplitPlaylist:   true,
		ChunkSize:       4,
	}

	var fileData = FileData{
		FilesList: []FileEntry{
			{"pathOne", "fileOne.mp3"},
			{"pathOne", "fileTwo.flac"},
			{"pathTwo", "fileOne.qt"},
			{"pathTwo", "fileTwo.a"},
		},
	}

	err = WritePlaylist(&configData, &fileData)
	if err != nil {
		t.Fatal("Expected no error when writing playlist", err)
	}

	file1, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	reader := csv.NewReader(file1)
	lines, err := reader.ReadAll()

	// 1 header, 4 entries *2 = 8 + 1 = 9
	if len(lines) != 9 {
		t.Fatal("Expected number of lines in file to be 9")
	}

	_, err = os.Stat(configData.OutputPath + "/" + configData.PlaylistName + "_02.m3u8")
	if !os.IsNotExist(err) {
		t.Fatal("second playlist should not exist")
	}
}

func TestPlaylistWriter_writePlaylist_ShouldNotSplitAPaylistWhenEntriesExceedSplitSizeAndSplitPlaylistIsFalse(t *testing.T) {
	var outputPath, err = os.MkdirTemp("", "PlaylistWriterTest")
	if err != nil {
		t.Fatal("Expected err to be nil", err)
	}
	if outputPath == "" {
		t.Fatal("Expected parent directory to not be empty")
	}

	var configData = config.Config{
		OutputPath:      outputPath,
		PlaylistName:    "testPlaylist",
		ShufflePlaylist: false,
		ReadTags:        false,
		SplitPlaylist:   false,
		ChunkSize:       2,
	}

	var fileData = FileData{
		FilesList: []FileEntry{
			{"pathOne", "fileOne.mp3"},
			{"pathOne", "fileTwo.flac"},
			{"pathTwo", "fileOne.qt"},
			{"pathTwo", "fileTwo.a"},
		},
	}

	err = WritePlaylist(&configData, &fileData)
	if err != nil {
		t.Fatal("Expected no error when writing playlist", err)
	}

	file1, err := os.Open(configData.OutputPath + "/" + configData.PlaylistName + "_01.m3u8")
	reader := csv.NewReader(file1)
	lines, err := reader.ReadAll()

	// 1 header, 4 entries *2 = 8 + 1 = 9
	if len(lines) != 9 {
		t.Fatal("Expected number of lines in file to be 9")
	}

	_, err = os.Stat(configData.OutputPath + "/" + configData.PlaylistName + "_02.m3u8")
	if !os.IsNotExist(err) {
		t.Fatal("second playlist should not exist")
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
