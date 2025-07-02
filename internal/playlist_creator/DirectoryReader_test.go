package playlist_creator

import (
	config2 "awesomeProject/internal/playlist_creator/config"
	"testing"
)

func TestDirectoryReader_IsExtensionAllowed_ShouldAllowExtension(t *testing.T) {
	var config = config2.Config{
		ExtensionWhitelist: []string{"ext1", "ext2", "ext3"},
	}

	if !isExtensionAllowed(&config, "ext2") {
		t.Fatal("Expected extension 'ext2' to be allowed")
	}
}

func TestDirectoryReader_IsExtensionAllowed_ShouldDenyExtension(t *testing.T) {
	var config = config2.Config{
		ExtensionWhitelist: []string{"ext1", "ext2", "ext3"},
	}

	if isExtensionAllowed(&config, "ext4") {
		t.Fatal("Expected extension 'ext4' to not be allowed")
	}
}

func TestDirectoryReader_IsExtensionAllowed_ShouldAllowAnyExtensionIfTheListIsNil(t *testing.T) {
	var config = config2.Config{}

	if !isExtensionAllowed(&config, "ext4") {
		t.Fatal("Expected extension 'ext4' to be allowed")
	}
}

func TestDirectoryReader_IsExtensionAllowed_ShouldAllowAnyExtensionIfTheListIsEmpty(t *testing.T) {
	var config = config2.Config{
		ExtensionWhitelist: []string{},
	}

	if !isExtensionAllowed(&config, "ext4") {
		t.Fatal("Expected extension 'ext4' to be allowed")
	}
}
