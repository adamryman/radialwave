// Taken blindly from here https://play.golang.org/p/hTb7CsUjuZ
package wave

import (
	"encoding/binary"
	"io"
)

// Format : .wav data structure
type Format struct {
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
	Data          []byte
}

// Decode : decode wav data
func (w *Format) Decode(r io.Reader) error {
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

	w.Data = make([]byte, w.Subchunk2Size)

	if _, err := io.ReadFull(r, w.Data); err != nil {
		return err
	}

	return nil
}
