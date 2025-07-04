package playlist_creator

import (
	"bytes"
	"fmt"
	"os"
	"playlistCreator/internal/playlist_creator/config"
	"testing"
)

func TestFileData_addFileExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		extensions             []string
		expectedExtensionsSize int
		expectedExtensions     map[string]int
	}{
		{
			name:                   "Should add an extension that does not exist",
			extensions:             []string{"abc"},
			expectedExtensionsSize: 1,
			expectedExtensions:     map[string]int{"abc": 1},
		},
		{
			name:                   "Should add multiple extensions",
			extensions:             []string{"abc", "123"},
			expectedExtensionsSize: 2,
			expectedExtensions:     map[string]int{"abc": 1, "123": 1},
		},
		{
			name:                   "Should increment extensions",
			extensions:             []string{"abc", "123", "abc"},
			expectedExtensionsSize: 2,
			expectedExtensions:     map[string]int{"abc": 2, "123": 1},
		},
	}

	for _, testData := range tests {
		testData := testData
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			var fileData = FileData{}
			if fileData.UniqueExtensions != nil {
				t.Error("Expected UniqueExtensions to be nil")
			}

			for index := range testData.extensions {
				fileData.addFileExtension(testData.extensions[index])
			}

			if len(fileData.UniqueExtensions) != testData.expectedExtensionsSize {
				t.Errorf("Expected UniqueExtensions to contain '%d' elements. Actual: '%d'.", testData.expectedExtensionsSize, len(fileData.UniqueExtensions))
			}

			for extension, extensionCount := range testData.expectedExtensions {
				if _, ok := fileData.UniqueExtensions[extension]; !ok {
					t.Errorf("Expected UniqueExtensions key %s to be present", extension)
				}

				if actualCount, _ := fileData.UniqueExtensions[extension]; actualCount != extensionCount {
					t.Errorf("Expected UniqueExtensions key '%s' to contain a value of '%d'. Actual: '%d'", extension, extensionCount, actualCount)
				}
			}
		})
	}
}

func TestFileData_incrementUniqueExtensions_ShouldIgnoreCase(t *testing.T) {
	var fileData = FileData{}

	fileData.addFileExtension("abc")
	fileData.addFileExtension("ABC")
	fileData.addFileExtension("aBc")
	fileData.addFileExtension("Abc")

	if count, _ := fileData.UniqueExtensions["abc"]; count != 4 {
		t.Fatal("Expected UniqueExtensions key 'abc' to contain a value of '4'")
	}
}

func TestFileData_listFileExtensions_ShouldPrintFileExtensions(t *testing.T) {
	var fileData = FileData{}

	fileData.addFileExtension("abc")
	fileData.addFileExtension("123")
	fileData.addFileExtension("abc")

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFileExtensions()

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("Extensions:\n")) {
		t.Fatal("Expected output to contain 'Extensions:\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("abc: 2\n")) {
		t.Fatal("Expected output to contain 'abc: 2\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("123: 1\n")) {
		t.Fatal("Expected output to contain '123: 1\\n'")
	}
}

func TestFileData_listFileExtensions_ShouldHandleNilUniqueExtensions(t *testing.T) {
	var fileData = FileData{}

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFileExtensions()

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("Extensions:\n")) {
		t.Fatal("Expected output to contain 'Extensions:\\n'")
	}
}

func TestFileData_ListFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		listLimit                 int
		files                     []FileEntry
		expectedNumPrintedFiles   int
		expectedNumFiles          int
		expectTruncatedList       bool
		expectContinuationEllipse bool
	}{
		{
			name:      "Should print files list",
			listLimit: -1,
			files: []FileEntry{
				{"somePath", "someFileName"},
				{"anotherPath", "anotherFileName"},
			},
			expectedNumFiles: 2,
		},
		{
			name:             "Should handle nil files list",
			listLimit:        -1,
			expectedNumFiles: 0,
		},
		{
			name:      "Should obey list limit",
			listLimit: 3,
			files: []FileEntry{
				{"somePath", "someFileName"},
				{"anotherPath", "anotherFileName"},
				{"somePath", "aDifferentFileName"},
				{"someOtherPath", "someFileName"},
			},
			expectedNumPrintedFiles:   3,
			expectedNumFiles:          4,
			expectTruncatedList:       true,
			expectContinuationEllipse: true,
		},
		{
			name:      "Should handle list limit larger than actual file count",
			listLimit: 10,
			files: []FileEntry{
				{"somePath", "someFileName"},
				{"anotherPath", "anotherFileName"},
				{"somePath", "aDifferentFileName"},
				{"someOtherPath", "someFileName"},
			},
			expectedNumPrintedFiles: 4,
			expectedNumFiles:        4,
			expectTruncatedList:     true,
		},
		{
			name:      "Should handle list limit of zero",
			listLimit: 0,
			files: []FileEntry{
				{"somePath", "someFileName"},
				{"anotherPath", "anotherFileName"},
				{"somePath", "aDifferentFileName"},
				{"someOtherPath", "someFileName"},
			},
			expectedNumPrintedFiles: 0,
			expectedNumFiles:        4,
			expectTruncatedList:     true,
		},
	}

	for _, testData := range tests {
		testData := testData
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()
		})

		var configData = config.Config{
			ListLimit: testData.listLimit,
		}

		var fileData = FileData{
			FilesList: testData.files,
		}

		var buf bytes.Buffer
		old := os.Stdout
		input, output, _ := os.Pipe()
		os.Stdout = output

		fileData.ListFiles(&configData)

		output.Close()
		os.Stdout = old
		buf.ReadFrom(input)

		actualOutput := buf.String()

		var expectedOutput string
		if testData.expectTruncatedList {
			expectedOutput = fmt.Sprintf("First %d files in list of length %d:", testData.expectedNumPrintedFiles, testData.expectedNumFiles)
		} else {
			expectedOutput = fmt.Sprintf("Files in list of length %d:", testData.expectedNumFiles)
		}
		if !bytes.Contains([]byte(actualOutput), []byte(expectedOutput+"\n")) {
			t.Errorf("Actual output doesn't match expected output. Expecting: %s, \\n'", expectedOutput)
		}

		for index := range testData.files[:testData.expectedNumPrintedFiles] {
			var fileName = testData.files[index].Path + " - " + testData.files[index].FileName

			if !bytes.Contains([]byte(actualOutput), []byte(fileName+"\n")) {
				t.Errorf("Expected output to contain '%s\\n'", fileName)
			}
		}

		if testData.expectContinuationEllipse {
			if !bytes.Contains([]byte(actualOutput), []byte("...")) {
				t.Error("Expected output to contain '...\\n'")
			}
		}
	}
}
