package main

import (
	&#34;bytes&#34;
	&#34;encoding/binary&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;io/ioutil&#34;
	&#34;os&#34;
	&#34;path/filepath&#34;
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
	if err := binary.Read(r, binary.BigEndian, &amp;w.ChunkID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.ChunkSize); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &amp;w.Format); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &amp;w.Subchunk1ID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.Subchunk1Size); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.AudioFormat); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.NumChannels); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.SampleRate); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.ByteRate); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.BlockAlign); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.BitsPerSample); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &amp;w.Subchunk2ID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &amp;w.Subchunk2Size); err != nil {
		return err
	}

	w.data = make([]byte, w.Subchunk2Size)

	if _, err := io.ReadFull(r, w.data); err != nil {
		return err
	}

	return nil
}

func main() {

	files, err := filepath.Glob(&#34;xa*&#34;)
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
		fmt.Println(&#34;Success! Bytes remaining:&#34;, len(accum))
	}
}