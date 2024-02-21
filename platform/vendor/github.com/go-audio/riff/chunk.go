package riff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

// Chunk represents the header and containt of a sub block
// See https://tech.ebu.ch/docs/tech/tech3285.pdf to see how
// audio content is stored in a BWF/WAVE file.
type Chunk struct {
	ID     [4]byte
	Size   int
	Pos    int
	R      io.Reader
	okChan chan bool
	Wg     *sync.WaitGroup
}

func (ch *Chunk) DecodeWavHeader(p *Parser) error {
	if ch == nil {
		return fmt.Errorf("can't decode a nil chunk")
	}
	if ch.ID == FmtID {
		p.wavHeaderSize = uint32(ch.Size)
		if err := ch.ReadLE(&p.WavAudioFormat); err != nil {
			return err
		}
		if err := ch.ReadLE(&p.NumChannels); err != nil {
			return err
		}
		if err := ch.ReadLE(&p.SampleRate); err != nil {
			return err
		}
		if err := ch.ReadLE(&p.AvgBytesPerSec); err != nil {
			return err
		}
		if err := ch.ReadLE(&p.BlockAlign); err != nil {
			return err
		}
		if err := ch.ReadLE(&p.BitsPerSample); err != nil {
			return err
		}

		// if we aren't dealing with a PCM file, we advance to reader to the
		// end of the chunck.
		if ch.Size > 16 {
			extra := make([]byte, ch.Size-16)
			ch.ReadLE(&extra)
		}
	}
	return nil
}

// Done signals the parent parser that we are done reading the chunk
// if the chunk isn't fully read, this code will do so before signaling.
func (ch *Chunk) Done() {
	if !ch.IsFullyRead() {
		ch.Drain()
	}
	if ch.Wg != nil {
		ch.Wg.Done()
	}
}

// IsFullyRead checks if we're finished reading the chunk
func (ch *Chunk) IsFullyRead() bool {
	if ch == nil || ch.R == nil {
		return true
	}
	return ch.Size <= ch.Pos
}

// Read implements the reader interface
func (ch *Chunk) Read(p []byte) (n int, err error) {
	if ch == nil || ch.R == nil {
		return 0, errors.New("nil chunk/reader pointer")
	}
	n, err = ch.R.Read(p)
	ch.Pos += n
	return n, err
}

// ReadLE reads the Little Endian chunk data into the passed struct
func (ch *Chunk) ReadLE(dst interface{}) error {
	if ch == nil || ch.R == nil {
		return errors.New("nil chunk/reader pointer")
	}
	if ch.IsFullyRead() {
		return io.EOF
	}
	ch.Pos += binary.Size(dst)
	return binary.Read(ch.R, binary.LittleEndian, dst)
}

// ReadBE reads the Big Endian chunk data into the passed struct
func (ch *Chunk) ReadBE(dst interface{}) error {
	if ch.IsFullyRead() {
		return io.EOF
	}
	ch.Pos += binary.Size(dst)
	return binary.Read(ch.R, binary.LittleEndian, dst)
}

// ReadByte reads and returns a single byte
func (ch *Chunk) ReadByte() (byte, error) {
	if ch.IsFullyRead() {
		return 0, io.EOF
	}
	var r byte
	err := ch.ReadLE(&r)
	return r, err
}

// Drain discards the rest of the chunk
func (ch *Chunk) Drain() {
	bytesAhead := ch.Size - ch.Pos
	for bytesAhead > 0 {
		readSize := int64(bytesAhead)

		if _, err := io.CopyN(ioutil.Discard, ch.R, readSize); err != nil {
			return
		}
		bytesAhead -= int(readSize)
	}
}
