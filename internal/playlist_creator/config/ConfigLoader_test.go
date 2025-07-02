package config

import (
	"encoding/json"
	"os"
	"testing"
)

func Test_loadConfigFromFile_ShouldLoadConfig(t *testing.T) {
	var temp = Config{
		InputPath:          "someInputPath",
		OutputPath:         "someOutputPath",
		PlaylistName:       "somePlaylistName",
		ExtensionWhitelist: []string{"ext1", "ext2", "ext3"},
		ShufflePlaylist:    true,
		SplitPlaylist:      true,
		ChunkSize:          30,
		ListExtensions:     true,
		ListFiles:          true,
		ListLimit:          34,
		ReadTags:           true,
	}

	var bytes, err = json.Marshal(temp)
	if err != nil {
		t.Fatal("Could not marshal json")
	}

	var file, error3 = os.CreateTemp("", "config.json")
	if error3 != nil {
		t.Fatal(error3)
	}
	os.WriteFile(file.Name(), bytes, 0644)

	var config, error2 = loadConfigFromFile(file.Name())
	if error2 != nil {
		t.Fatal("Expected err to be nil", err)
	}

	if config == nil {
		t.Fatal("Expected config to be non-nil")
	}

	if config.InputPath != "someInputPath" {
		t.Fatal("Expected config.InputPath to be 'someInputPath'")
	}

	if config.OutputPath != "someOutputPath" {
		t.Fatal("Expected config.OutputPath to be 'someOutputPath'")
	}

	if config.PlaylistName != "somePlaylistName" {
		t.Fatal("Expected config.PlaylistName to be 'somePlaylistName'")
	}

	if config.ExtensionWhitelist == nil {
		t.Fatal("Expected config.ExtensionWhitelist to be non-nil")
	}
	if len(config.ExtensionWhitelist) != 3 {
		t.Fatal("Expected config.ExtensionWhitelist to contain 3 elements")
	}
	if config.ExtensionWhitelist[0] != "ext1" {
		t.Fatal("Expected config.ExtensionWhitelist[0] to be 'ext1'")
	}
	if config.ExtensionWhitelist[1] != "ext2" {
		t.Fatal("Expected config.ExtensionWhitelist[0] to be 'ext2'")
	}
	if config.ExtensionWhitelist[2] != "ext3" {
		t.Fatal("Expected config.ExtensionWhitelist[0] to be 'ext3'")
	}
	if config.SplitPlaylist != true {
		t.Fatal("Expected config.SplitPlaylist to be true")
	}
	if config.ChunkSize != 30 {
		t.Fatal("Expected config.ChunkSize to be 30")
	}
	if config.ListExtensions != true {
		t.Fatal("Expected config.ListExtensions to be true")
	}
	if config.ListFiles != true {
		t.Fatal("Expected config.ListFiles to be true")
	}
	if config.ListLimit != 34 {
		t.Fatal("Expected config.ListLimit to be 34")
	}
	if config.ReadTags != true {
		t.Fatal("Expected config.ReadTags to be true")
	}
}

func Test_loadConfigFromFile_ShouldFailIfInputPathIsEmpty(t *testing.T) {
	var temp = Config{
		OutputPath:         "someOutputPath",
		PlaylistName:       "somePlaylistName",
		ExtensionWhitelist: []string{"ext1", "ext2", "ext3"},
		ShufflePlaylist:    true,
		SplitPlaylist:      true,
		ChunkSize:          30,
		ListExtensions:     true,
		ListFiles:          true,
		ListLimit:          34,
		ReadTags:           true,
	}

	var bytes, err = json.Marshal(temp)
	if err != nil {
		t.Fatal("Could not marshal json")
	}

	var file, error3 = os.CreateTemp("", "config.json")
	if error3 != nil {
		t.Fatal(error3)
	}
	os.WriteFile(file.Name(), bytes, 0644)

	var _, error2 = loadConfigFromFile(file.Name())
	if error2 == nil {
		t.Fatal("Expected err to not be nil", err)
	}

	if error2.Error() != "config.InputPath is empty" {
		t.Fatal("Expected error to be 'config.InputPath is empty'")
	}
}

func Test_loadConfigFromFile_ShouldFailIfOutputPathIsEmpty(t *testing.T) {
	var temp = Config{
		InputPath:          "someInputPath",
		PlaylistName:       "somePlaylistName",
		ExtensionWhitelist: []string{"ext1", "ext2", "ext3"},
		ShufflePlaylist:    true,
		SplitPlaylist:      true,
		ChunkSize:          30,
		ListExtensions:     true,
		ListFiles:          true,
		ListLimit:          34,
		ReadTags:           true,
	}

	var bytes, error1 = json.Marshal(temp)
	if error1 != nil {
		t.Fatal("Could not marshal json")
	}

	var file, error2 = os.CreateTemp("", "config.json")
	if error2 != nil {
		t.Fatal(error2)
	}
	os.WriteFile(file.Name(), bytes, 0644)

	var _, error3 = loadConfigFromFile(file.Name())
	if error3 == nil {
		t.Fatal("Expected err to not be nil", error3)
	}

	if error3.Error() != "config.OutputPath is empty" {
		t.Fatal("Expected error to be 'config.OutputPath is empty'")
	}
}
