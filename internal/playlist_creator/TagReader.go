package playlist_creator

import (
	"github.com/bogem/id3v2/v2"
	"log"
)

func createPlaylistEntryName(fileEntry FileEntry) string {

	tag, err := id3v2.Open(fileEntry.Path+"/"+fileEntry.FileName, id3v2.Options{Parse: true})
	if err != nil {
		log.Printf("Warning: unable to get taf from file: %s", err.Error())
		return ""
	}
	defer tag.Close()

	if tag.Artist() == "" && tag.Title() == "" {
		return ""
	}

	artist := "Unknown Artist"
	if tag.Artist() != "" {
		artist = tag.Artist()
	}

	title := "Unknown Title"
	if tag.Title() != "" {
		title = tag.Title()
	}

	return artist + " - " + title
}
