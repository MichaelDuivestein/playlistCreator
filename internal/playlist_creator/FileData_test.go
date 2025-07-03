package playlist_creator

import (
	"bytes"
	"os"
	"playlistCreator/internal/playlist_creator/config"
	"testing"
)

func TestFileData_addFileExtension_ShouldAddAnExtensionThatDoesNotExist(t *testing.T) {
	var fileData = FileData{}

	if fileData.UniqueExtensions != nil {
		t.Fatal("Expected UniqueExtensions to be nil")
	}

	fileData.addFileExtension("abc")
	if fileData.UniqueExtensions == nil {
		t.Fatal("Expected UniqueExtensions to not be nil")
	}

	if len(fileData.UniqueExtensions) != 1 {
		t.Fatal("Expected UniqueExtensions to contain 1 element")
	}

	if _, ok := fileData.UniqueExtensions["abc"]; !ok {
		t.Fatal("Expected UniqueExtensions key 'abc' to be present")
	}

	if count, _ := fileData.UniqueExtensions["abc"]; count != 1 {
		t.Fatal("Expected UniqueExtensions key 'abc' to contain a value of '1'")
	}
}

func TestFileData_addFileExtension_ShouldAddMultipleExtensions(t *testing.T) {
	var fileData = FileData{}

	fileData.addFileExtension("abc")
	fileData.addFileExtension("123")

	if len(fileData.UniqueExtensions) != 2 {
		t.Fatal("Expected UniqueExtensions to contain 2 elements")
	}

	if _, ok := fileData.UniqueExtensions["abc"]; !ok {
		t.Fatal("Expected UniqueExtensions key 'abc' to be present")
	}
	if count, _ := fileData.UniqueExtensions["abc"]; count != 1 {
		t.Fatal("Expected UniqueExtensions key 'abc' to contain a value of '1'")
	}

	if _, ok := fileData.UniqueExtensions["123"]; !ok {
		t.Fatal("Expected UniqueExtensions key '123' to be present")
	}
	if count, _ := fileData.UniqueExtensions["123"]; count != 1 {
		t.Fatal("Expected UniqueExtensions key '123' to contain a value of '1'")
	}
}

func TestFileData_addFileExtension_ShouldIncrementExtensionCount(t *testing.T) {
	var fileData = FileData{}

	fileData.addFileExtension("abc")
	fileData.addFileExtension("123")
	fileData.addFileExtension("abc")

	if _, ok := fileData.UniqueExtensions["abc"]; !ok {
		t.Fatal("Expected UniqueExtensions key 'abc' to be present")
	}
	if count, _ := fileData.UniqueExtensions["abc"]; count != 2 {
		t.Fatal("Expected UniqueExtensions key 'abc' to contain a value of '2'")
	}

	if _, ok := fileData.UniqueExtensions["123"]; !ok {
		t.Fatal("Expected UniqueExtensions key '123' to be present")
	}
	if count, _ := fileData.UniqueExtensions["123"]; count != 1 {
		t.Fatal("Expected UniqueExtensions key '123' to contain a value of '1'")
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

func TestFileData_listFiles_ShouldPrintFilesList(t *testing.T) {
	var configData = config.Config{
		ListLimit: -1,
	}
	var fileData = FileData{}
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "someFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"anotherPath", "anotherFileName"})

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFiles(&configData)

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("Files in list of length 2:\n")) {
		t.Fatal("Expected output to contain 'Files in list of length 2:\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("somePath - someFileName\n")) {
		t.Fatal("Expected output to contain 'somePath - someFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("anotherPath - anotherFileName\n")) {
		t.Fatal("Expected output to contain 'anotherPath - anotherFileName\\n'")
	}
}

func TestFileData_listFiles_ShouldHandleNilFilesList(t *testing.T) {
	var configData = config.Config{
		ListLimit: -1,
	}
	var fileData = FileData{}

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFiles(&configData)

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("Files in list of length 0:\n")) {
		t.Fatal("Expected output to contain 'Files in list of length 0:\\n'")
	}
}

func TestFileData_listFiles_ShouldObeyListLimit(t *testing.T) {
	var configData = config.Config{
		ListLimit: 3,
	}
	var fileData = FileData{}
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "someFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"anotherPath", "anotherFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "aDifferentFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"someOtherPath", "someFileName"})

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFiles(&configData)

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("First 3 files in list of length 4:\n")) {
		t.Fatal("Expected output to contain 'First 3 files in list of length 4:\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("somePath - someFileName\n")) {
		t.Fatal("Expected output to contain 'somePath - someFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("anotherPath - anotherFileName\n")) {
		t.Fatal("Expected output to contain 'anotherPath - anotherFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("somePath - aDifferentFileName\n")) {
		t.Fatal("Expected output to contain 'somePath - aDifferentFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("...")) {
		t.Fatal("Expected output to contain '...\\n'")
	}
}

func TestFileData_listFiles_ShouldHandleListLimitLargerThanActualFiles(t *testing.T) {
	var configData = config.Config{
		ListLimit: 10,
	}
	var fileData = FileData{}
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "someFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"anotherPath", "anotherFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "aDifferentFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"someOtherPath", "someFileName"})

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFiles(&configData)

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("First 4 files in list of length 4:\n")) {
		t.Fatal("Expected output to contain 'First 4 files in list of length 4:\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("somePath - someFileName\n")) {
		t.Fatal("Expected output to contain 'somePath - someFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("anotherPath - anotherFileName\n")) {
		t.Fatal("Expected output to contain 'anotherPath - anotherFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("somePath - aDifferentFileName\n")) {
		t.Fatal("Expected output to contain 'somePath - aDifferentFileName\\n'")
	}
	if bytes.Contains([]byte(actualOutput), []byte("...")) {
		t.Fatal("Expected output not to contain '...\\n'")
	}
}

func TestFileData_listFiles_ShouldHandleListLimitOfZero(t *testing.T) {
	var configData = config.Config{
		ListLimit: 0,
	}
	var fileData = FileData{}
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "someFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"anotherPath", "anotherFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"somePath", "aDifferentFileName"})
	fileData.FilesList = append(fileData.FilesList, FileEntry{"someOtherPath", "someFileName"})

	var buf bytes.Buffer
	old := os.Stdout
	input, output, _ := os.Pipe()
	os.Stdout = output

	fileData.ListFiles(&configData)

	output.Close()
	os.Stdout = old
	buf.ReadFrom(input)

	actualOutput := buf.String()

	if !bytes.Contains([]byte(actualOutput), []byte("First 0 files in list of length 4:\n")) {
		t.Fatal("Expected output to contain 'somePath - aDifferentFileName\\n'")
	}
	if bytes.Contains([]byte(actualOutput), []byte("somePath - someFileName\n")) {
		t.Fatal("Expected output to not contain 'somePath - aDifferentFileName\\n'")
	}
	if bytes.Contains([]byte(actualOutput), []byte("anotherPath - anotherFileName\n")) {
		t.Fatal("Expected output to not contain 'somePath - aDifferentFileName\\n'")
	}
	if bytes.Contains([]byte(actualOutput), []byte("somePath - aDifferentFileName\n")) {
		t.Fatal("Expected output to not contain 'somePath - aDifferentFileName\\n'")
	}
	if !bytes.Contains([]byte(actualOutput), []byte("...")) {
		t.Fatal("Expected output to contain 'somePath - aDifferentFileName\\n'")
	}
}
