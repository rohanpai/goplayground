// Main package to running BagIns from the commandline.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
)

type TagFile struct {
	Filepath string            // Filepath for tag file.
	Data     map[string]string // key value pairs of data for the tagfile.
}

// Writes key value pairs to a tag file.
func (tf *TagFile) Create() {
	// Create directory if needed.
	basepath := path.Dir(tf.Filepath)
	filename := path.Base(tf.Filepath)
	if os.MkdirAll(basepath, 0777) != nil {
		panic("Unable to create directory for tagfile!")
	}

	// Create the tagfile.
	fileOut, err := os.Create(path.Join(basepath, filename))
	if err != nil {
		panic("Unable to create tag file!")
	}
	defer fileOut.Close()

	// Write fields and data to the file.
	for key, data := range tf.Data {
		_, err := fmt.Fprintln(fileOut, formatField(key, data))
		if err != nil {
			panic("Unable to write data to tagfile.")
		}
	}
}

/*
Takes a tag field key and data and wraps lines at 79 with indented spaces as
per recommendation in spec.
*/
func formatField(key string, data string) string {
	delimeter := "\n   "
	var buff bytes.Buffer

	// Initiate it by writing the proper key.
	writeLen, err := buff.WriteString(fmt.Sprintf("%s: ", key))
	if err != nil {
		panic("Unable to begin writing field!")
	}
	splitCounter := writeLen

	words := strings.Split(data, " ")

	for word := range words {
		if splitCounter+len(words[word]) > 79 {
			splitCounter, err = buff.WriteString(delimeter)
			if err != nil {
				panic("Unable to write field!")
			}
		}
		writeLen, err = buff.WriteString(strings.Join([]string{" ", words[word]}, ""))
		if err != nil {
			panic("Unable to write field!")
		}
		splitCounter += writeLen

	}
	return buff.String()
}

func main() {
	data := map[string]string{
		"BagIt-Version":                `A metadata element MUST consist of a label, a colon, and a value, each separated by optional whitespace.  It is RECOMMENDED that lines not exceed 79 characters in length.  Long values may be continued onto the next line by inserting a newline (LF), a carriage return (CR), or carriage return plus newline (CRLF) and indenting the next line with linear white space (spaces or tabs).`,
		"Tag-File-Character-Encodeing": "UTF-8",
	}
	tagFile := TagFile{Filepath: "tagfiles/bagit.txt", Data: data}
	tagFile.Create()
}
