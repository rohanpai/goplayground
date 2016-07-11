// Main package to running BagIns from the commandline.
package main

import (
	&#34;bytes&#34;
	&#34;fmt&#34;
	&#34;os&#34;
	&#34;path&#34;
	&#34;strings&#34;
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
		panic(&#34;Unable to create directory for tagfile!&#34;)
	}

	// Create the tagfile.
	fileOut, err := os.Create(path.Join(basepath, filename))
	if err != nil {
		panic(&#34;Unable to create tag file!&#34;)
	}
	defer fileOut.Close()

	// Write fields and data to the file.
	for key, data := range tf.Data {
		_, err := fmt.Fprintln(fileOut, formatField(key, data))
		if err != nil {
			panic(&#34;Unable to write data to tagfile.&#34;)
		}
	}
}

/*
Takes a tag field key and data and wraps lines at 79 with indented spaces as
per recommendation in spec.
*/
func formatField(key string, data string) string {
	delimeter := &#34;\n   &#34;
	var buff bytes.Buffer

	// Initiate it by writing the proper key.
	writeLen, err := buff.WriteString(fmt.Sprintf(&#34;%s: &#34;, key))
	if err != nil {
		panic(&#34;Unable to begin writing field!&#34;)
	}
	splitCounter := writeLen

	words := strings.Split(data, &#34; &#34;)

	for word := range words {
		if splitCounter&#43;len(words[word]) &gt; 79 {
			splitCounter, err = buff.WriteString(delimeter)
			if err != nil {
				panic(&#34;Unable to write field!&#34;)
			}
		}
		writeLen, err = buff.WriteString(strings.Join([]string{&#34; &#34;, words[word]}, &#34;&#34;))
		if err != nil {
			panic(&#34;Unable to write field!&#34;)
		}
		splitCounter &#43;= writeLen

	}
	return buff.String()
}

func main() {
	data := map[string]string{
		&#34;BagIt-Version&#34;:                `A metadata element MUST consist of a label, a colon, and a value, each separated by optional whitespace.  It is RECOMMENDED that lines not exceed 79 characters in length.  Long values may be continued onto the next line by inserting a newline (LF), a carriage return (CR), or carriage return plus newline (CRLF) and indenting the next line with linear white space (spaces or tabs).`,
		&#34;Tag-File-Character-Encodeing&#34;: &#34;UTF-8&#34;,
	}
	tagFile := TagFile{Filepath: &#34;tagfiles/bagit.txt&#34;, Data: data}
	tagFile.Create()
}
