package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WavFormat : data structure
type WavFormat struct {
	ChunkID       uint32
	ChunkSize     uint32
	Format        uint32
	Subchunk1ID   uint32
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   uint32
	Subchunk2Size uint32
	data          []byte
}

// Decode : decode wav data
func (w *WavFormat) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &w.ChunkID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.ChunkSize); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &w.Format); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &w.Subchunk1ID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.Subchunk1Size); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.AudioFormat); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.NumChannels); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.SampleRate); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.ByteRate); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.BlockAlign); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.BitsPerSample); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &w.Subchunk2ID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &w.Subchunk2Size); err != nil {
		return err
	}

	w.data = make([]byte, w.Subchunk2Size)

	if _, err := io.ReadFull(r, w.data); err != nil {
		return err
	}

	return nil
}

func main() {

	files, err := filepath.Glob("xa*")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	accum := []byte{}

	w := WavFormat{}
	for _, file := range files {
		fmt.Println(file)

		b, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		accum = append(accum, b...)

		r := bytes.NewReader(accum)
		err = w.Decode(r)
		if err != nil {
			fmt.Println(err)
			continue
		}
		accum = accum[len(accum)-r.Len():]
		fmt.Println("Success! Bytes remaining:", len(accum))
	}
}