package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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

	bytes, err := json.Marshal(temp)
	assert.Nil(t, err, "Could not marshal json")

	file, err := os.CreateTemp("", "config.json")
	assert.Nil(t, err, "Could not create temp file")

	os.WriteFile(file.Name(), bytes, 0644)

	config, err := loadConfigFromFile(file.Name())
	assert.Nil(t, err, "Expected err to be nil")
	assert.NotNil(t, config, "Expected config to not be nil")
	assert.Equal(t, "someInputPath", config.InputPath, "Expected inputPath to be 'someInputPath'")
	assert.Equal(t, "someOutputPath", config.OutputPath, "Expected config.OutputPath to be 'someOutputPath'")
	assert.Equal(t, "somePlaylistName", config.PlaylistName, "Expected config.PlaylistName to be 'somePlaylistName'")

	assert.NotNil(t, config.ExtensionWhitelist, "Expected config.ExtensionWhitelist to be non-nil")
	assert.Equal(t, 3, len(config.ExtensionWhitelist), "Expected config.ExtensionWhitelist to contain 3 elements")
	assert.Equal(t, "ext1", config.ExtensionWhitelist[0], "Expected config.ExtensionWhitelist[0] to be 'ext1'")
	assert.Equal(t, "ext2", config.ExtensionWhitelist[1], "Expected config.ExtensionWhitelist[1] to be 'ext2'")
	assert.Equal(t, "ext3", config.ExtensionWhitelist[2], "Expected config.ExtensionWhitelist[2] to be 'ext3'")

	assert.True(t, config.SplitPlaylist, "Expected config.SplitPlaylist to be true")
	assert.Equal(t, 30, config.ChunkSize, "Expected config.ChunkSize to be 30")
	assert.True(t, config.ListExtensions, "Expected config.ListExtensions to be true")
	assert.True(t, config.ListFiles, "Expected config.ListFiles to be true")
	assert.Equal(t, 34, config.ListLimit, "Expected config.ListLimit to be 34")
	assert.True(t, config.ReadTags, "Expected config.ReadTags to be true")
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

	bytes, err := json.Marshal(temp)
	assert.Nil(t, err, "Expected err to be nil")

	file, err := os.CreateTemp("", "config.json")
	assert.Nil(t, err, "Expected err to be nil")
	os.WriteFile(file.Name(), bytes, 0644)

	_, err = loadConfigFromFile(file.Name())
	assert.NotNil(t, err, "Expected err to not be nil")

	assert.Equal(t, "config.InputPath is empty", err.Error(), "Expected error to be 'config.InputPath is empty'")
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

	bytes, err := json.Marshal(temp)
	assert.Nil(t, err, "Could not marshal json")

	file, err := os.CreateTemp("", "config.json")
	assert.Nil(t, err)
	os.WriteFile(file.Name(), bytes, 0644)

	_, err = loadConfigFromFile(file.Name())
	assert.NotNil(t, err, "Expected err to not be nil")

	assert.Equal(t, "config.OutputPath is empty", err.Error(), "Expected error to be 'config.OutputPath is empty'")
}
