package riff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// Parser is a struct containing the overall container information.
type Parser struct {
	r io.Reader
	// Chan is an Optional channel of chunks that is used to parse chunks
	Chan chan *Chunk
	// ChunkParserTimeout is the duration after which the main parser keeps going
	// if the dev hasn't reported the chunk parsing to be done.
	// By default: 2s
	ChunkParserTimeout time.Duration
	// The waitgroup is used to let the parser that it's ok to continue
	// after a chunk was passed to the optional parser channel.
	Wg *sync.WaitGroup

	// Must match RIFF
	ID [4]byte
	// This size is the size of the block
	// controlled by the RIFF header. Normally this equals the file size.
	Size uint32
	// Format name.
	// The representation of data in <wave-data>, and the content of the <format-specific-fields>
	// of the ‘fmt’ chunk, depend on the format category.
	// 0001h => Microsoft Pulse Code Modulation (PCM) format
	// 0050h => MPEG-1 Audio (audio only)
	Format [4]byte

	// WAV stuff
	// size of the wav specific fmt header
	wavHeaderSize uint32
	// A number indicating the WAVE format category of the file. The content of the
	// <format-specific-fields> portion of the ‘fmt’ chunk, and the interpretation of
	// the waveform data, depend on this value.
	// PCM = 1 (i.e. Linear quantization) Values other than 1 indicate some form of compression.
	WavAudioFormat uint16
	// The number of channels represented in the waveform data: 1 for mono or 2 for stereo.
	// Audio: Mono = 1, Stereo = 2, etc.
	// The EBU has defined the Multi-channel Broadcast Wave
	// Format [4] where more than two channels of audio are required.
	NumChannels uint16
	// The sampling rate (in sample per second) at which each channel should be played.
	// 8000, 44100, etc.
	SampleRate uint32
	// The average number of bytes per second at which the waveform data should be
	// transferred. Playback software can estimate the buffer size using this value.
	// SampleRate * NumChannels * BitsPerSample/8
	AvgBytesPerSec uint32
	// BlockAlign = SignificantBitsPerSample / 8 * NumChannels
	// It is the number of bytes per sample slice. This value is not affected by the number of channels and can be calculated with the formula:
	// NumChannels * BitsPerSample/8 The number of bytes for one sample including
	// all channels.
	// The block alignment (in bytes) of the waveform data. Playback software needs
	// to process a multiple of <nBlockAlign> bytes of data at a time, so the value of
	// <BlockAlign> can be used for buffer alignment.
	BlockAlign uint16
	// BitsPerSample 8, 16, 24...
	// Only available for PCM
	// This value specifies the number of bits used to define each sample. This value is usually 8, 16, 24 or 32.
	// If the number of bits is not byte aligned (a multiple of 8) then the number of bytes used per sample is
	// rounded up to the nearest byte size and the unused bytes are set to 0 and ignored.
	// The <nBitsPerSample> field specifies the number of bits of data used to represent each sample of
	// each channel. If there are multiple channels, the sample size is the same for each channel.
	BitsPerSample uint16
}

// ParseHeaders reads the header of the passed container and populat the container with parsed info.
// Note that this code advances the container reader.
func (c *Parser) ParseHeaders() error {
	id, size, err := c.IDnSize()
	if err != nil {
		return err
	}
	c.ID = id
	if c.ID != RiffID {
		return fmt.Errorf("%s - %s", c.ID, ErrFmtNotSupported)
	}
	c.Size = size
	if err := binary.Read(c.r, binary.BigEndian, &c.Format); err != nil {
		return err
	}

	return nil
}

// Duration returns the time duration for the current RIFF container
// based on the sub format (wav etc...)
func (c *Parser) Duration() (time.Duration, error) {
	if c == nil {
		return 0, errors.New("can't calculate the duration of a nil pointer")
	}
	if c.ID == [4]byte{} || c.AvgBytesPerSec == 0 {
		err := c.Parse()
		if err != nil {
			return 0, nil
		}
	}
	switch c.Format {
	case WavFormatID:
		return c.wavDuration()
	default:
		return 0, ErrFmtNotSupported
	}
}

// String implements the Stringer interface.
func (c *Parser) String() string {
	out := fmt.Sprintf("Format: %s - ", c.Format)
	if c.Format == WavFormatID {
		out += fmt.Sprintf("%d channels @ %d / %d bits - ", c.NumChannels, c.SampleRate, c.BitsPerSample)
		d, _ := c.Duration()
		out += fmt.Sprintf("Duration: %f seconds", d.Seconds())
	}
	return out
}

// NextChunk returns a convenient structure to parse the next chunk.
// If the container is fully read, io.EOF is returned as an error.
func (c *Parser) NextChunk() (*Chunk, error) {
	if c == nil {
		return nil, errors.New("can't calculate the duration of a nil pointer")
	}
	id, size, err := c.IDnSize()
	if err != nil {
		return nil, err
	}

	// all RIFF chunks (including WAVE "data" chunks) must be word aligned.
	// If the data uses an odd number of bytes, a padding byte with a value of zero must be placed at the end of the sample data.
	// The "data" chunk header's size should not include this byte.
	if size%2 == 1 {
		size++
	}

	ch := &Chunk{
		ID:   id,
		Size: int(size),
		R:    c.r,
	}
	return ch, nil
}

// IDnSize returns the next ID + block size
func (c *Parser) IDnSize() ([4]byte, uint32, error) {
	var ID [4]byte
	var blockSize uint32
	if err := binary.Read(c.r, binary.BigEndian, &ID); err != nil {
		return ID, blockSize, err
	}
	if err := binary.Read(c.r, binary.LittleEndian, &blockSize); err != err {
		return ID, blockSize, err
	}
	return ID, blockSize, nil
}

// Parse parses the content of the file and populate the useful fields.
// If the parser has a chan set, chunks are sent to the channel.
func (p *Parser) Parse() error {
	if p == nil {
		return errors.New("can't calculate the wav duration of a nil pointer")
	}

	if p.Size == 0 {
		id, size, err := p.IDnSize()
		if err != nil {
			return err
		}
		p.ID = id
		if p.ID != RiffID {
			return fmt.Errorf("%s - %s", p.ID, ErrFmtNotSupported)
		}
		p.Size = size
		if err := binary.Read(p.r, binary.BigEndian, &p.Format); err != nil {
			return err
		}
	}

	var chunk *Chunk
	var err error
	for err == nil {
		chunk, err = p.NextChunk()
		if err != nil {
			break
		}

		if chunk.ID == FmtID {
			chunk.DecodeWavHeader(p)
		} else {
			if p.Chan != nil {
				if chunk.Wg == nil {
					chunk.Wg = p.Wg
				}
				chunk.Wg.Add(1)
				p.Chan <- chunk
				// the channel has to release otherwise the goroutine is locked
				chunk.Wg.Wait()
			}
		}

		// BFW: bext chunk described here
		// https://tech.ebu.ch/docs/tech/tech3285.pdf

		if !chunk.IsFullyRead() {
			chunk.Drain()
		}

	}
	if p.Wg != nil {
		p.Wg.Wait()
	}

	if p.Chan != nil {
		close(p.Chan)
	}

	if err == io.EOF {
		return nil
	}
	return err
}

// WavDuration returns the time duration of a wav container.
func (p *Parser) wavDuration() (time.Duration, error) {
	if p.Size == 0 || p.AvgBytesPerSec == 0 {
		return 0, fmt.Errorf("can't extract the duration due to the file not properly parsed")
	}
	duration := time.Duration((float64(p.Size) / float64(p.AvgBytesPerSec)) * float64(time.Second))
	return duration, nil
}

// jumpTo advances the reader to the amount of bytes provided
func (p *Parser) jumpTo(bytesAhead int) error {
	var err error
	for bytesAhead > 0 {
		readSize := bytesAhead
		if readSize > 4000 {
			readSize = 4000
		}

		buf := make([]byte, readSize)
		err = binary.Read(p.r, binary.LittleEndian, &buf)
		if err != nil {
			return nil
		}
		bytesAhead -= readSize
	}
	return nil
}
