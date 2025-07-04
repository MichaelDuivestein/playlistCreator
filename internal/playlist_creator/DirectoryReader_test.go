package playlist_creator

import (
	"playlistCreator/internal/playlist_creator/config"
	"testing"
)

func TestIsExtensionAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		whitelist       []string
		extension       string
		expectedAllowed bool
	}{
		{
			name:            "Extension should be allowed",
			whitelist:       []string{"ext1", "ext2", "ext3"},
			extension:       "ext2",
			expectedAllowed: true,
		},
		{
			name:            "Extension should be denied",
			whitelist:       []string{"ext1", "ext2", "ext3"},
			extension:       "ext4",
			expectedAllowed: false,
		},
		{
			name:            "Should allow any extension if whitelist is nil",
			whitelist:       nil,
			extension:       "ext4",
			expectedAllowed: true,
		},
		{
			name:            "Should allow any extension if whitelist is empty",
			whitelist:       []string{},
			extension:       "ext4",
			expectedAllowed: true,
		},
		{
			name:            "Should match case-insensitive extensions",
			whitelist:       []string{"EXT1"},
			extension:       "ext1",
			expectedAllowed: true,
		},
	}

	for _, testData := range tests {
		testData := testData
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			configData := config.Config{
				ExtensionWhitelist: testData.whitelist,
			}

			result := isExtensionAllowed(&configData, testData.extension)
			if result != testData.expectedAllowed {
				t.Errorf("isExtensionAllowed(%v, %q) = %v; want %v",
					testData.whitelist, testData.extension, result, testData.expectedAllowed)
			}
		})
	}
}
