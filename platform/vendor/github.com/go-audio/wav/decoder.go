package wav

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/riff"
)

var (
	// CIDList is the chunk ID for a LIST chunk
	CIDList = [4]byte{'L', 'I', 'S', 'T'}
	// CIDSmpl is the chunk ID for a smpl chunk
	CIDSmpl = [4]byte{'s', 'm', 'p', 'l'}
	// CIDINFO is the chunk ID for an INFO chunk
	CIDInfo = []byte{'I', 'N', 'F', 'O'}
	// CIDCue is the chunk ID for the cue chunk
	CIDCue = [4]byte{'c', 'u', 'e', 0x20}
)

// Decoder handles the decoding of wav files.
type Decoder struct {
	r      io.ReadSeeker
	parser *riff.Parser

	NumChans   uint16
	BitDepth   uint16
	SampleRate uint32

	AvgBytesPerSec uint32
	WavAudioFormat uint16

	err             error
	PCMSize         int
	pcmDataAccessed bool
	// pcmChunk is available so we can use the LimitReader
	PCMChunk *riff.Chunk
	// Metadata for the current file
	Metadata *Metadata
}

// NewDecoder creates a decoder for the passed wav reader.
// Note that the reader doesn't get rewinded as the container is processed.
func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{
		r:      r,
		parser: riff.New(r),
	}
}

// Seek provides access to the cursor position in the PCM data
func (d *Decoder) Seek(offset int64, whence int) (int64, error) {
	return d.r.Seek(offset, whence)
}

// Rewind allows the decoder to be rewound to the beginning of the PCM data.
// This is useful if you want to keep on decoding the same file in a loop.
func (d *Decoder) Rewind() error {

	_, err := d.r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek back to the start %w", err)
	}
	// we have to user a new parser since it's read only and can't be seeked
	d.parser = riff.New(d.r)
	d.pcmDataAccessed = false
	d.PCMChunk = nil
	d.err = nil
	d.NumChans = 0
	err = d.FwdToPCM()
	if err != nil {
		return fmt.Errorf("failed to seek to the PCM data: %w", err)
	}
	return nil
}

// SampleBitDepth returns the bit depth encoding of each sample.
func (d *Decoder) SampleBitDepth() int32 {
	if d == nil {
		return 0
	}
	return int32(d.BitDepth)
}

// PCMLen returns the total number of bytes in the PCM data chunk
func (d *Decoder) PCMLen() int64 {
	if d == nil {
		return 0
	}
	return int64(d.PCMSize)
}

// Err returns the first non-EOF error that was encountered by the Decoder.
func (d *Decoder) Err() error {
	if d.err == io.EOF {
		return nil
	}
	return d.err
}

// EOF returns positively if the underlying reader reached the end of file.
func (d *Decoder) EOF() bool {
	if d == nil || d.err == io.EOF {
		return true
	}
	return false
}

// IsValidFile verifies that the file is valid/readable.
func (d *Decoder) IsValidFile() bool {
	d.err = d.readHeaders()
	if d.err != nil {
		return false
	}
	if d.NumChans < 1 {
		return false
	}
	if d.BitDepth < 8 {
		return false
	}
	if d, err := d.Duration(); err != nil || d <= 0 {
		return false
	}

	return true
}

// ReadInfo reads the underlying reader until the comm header is parsed.
// This method is safe to call multiple times.
func (d *Decoder) ReadInfo() {
	d.err = d.readHeaders()
}

// ReadMetadata parses the file for extra metadata such as the INFO list chunk.
// The entire file will be read and should be rewinded if more data must be
// accessed.
func (d *Decoder) ReadMetadata() {
	if d.Metadata != nil {
		return
	}
	d.ReadInfo()
	if d.Err() != nil || d.Metadata != nil {
		return
	}
	var (
		chunk *riff.Chunk
		err   error
	)
	for err == nil {
		chunk, err = d.parser.NextChunk()
		if err != nil {
			break
		}

		switch chunk.ID {
		case CIDList:
			if err = DecodeListChunk(d, chunk); err != nil {
				if err != io.EOF {
					d.err = err
				}
			}
			if d.Metadata != nil && d.Metadata.SamplerInfo != nil {
				// we got everything we were looking for
				break
			}
		case CIDSmpl:
			if err = DecodeSamplerChunk(d, chunk); err != nil {
				if err != io.EOF {
					d.err = err
				}
			}
		case CIDCue:
			if err = DecodeCueChunk(d, chunk); err != nil {
				if err != io.EOF {
					d.err = err
				}
			}
		default:
			// fmt.Println(string(chunk.ID[:]))
			chunk.Drain()
		}
	}

}

// FwdToPCM forwards the underlying reader until the start of the PCM chunk.
// If the PCM chunk was already read, no data will be found (you need to rewind).
func (d *Decoder) FwdToPCM() error {
	if d == nil {
		return fmt.Errorf("PCM data not found")
	}
	d.err = d.readHeaders()
	if d.err != nil {
		return nil
	}

	var chunk *riff.Chunk
	for d.err == nil {
		chunk, d.err = d.NextChunk()
		if d.err != nil {
			return d.err
		}
		if chunk.ID == riff.DataFormatID {
			d.PCMSize = chunk.Size
			d.PCMChunk = chunk
			break
		}
		if chunk.ID == CIDList {
			DecodeListChunk(d, chunk)
		}
		chunk.Drain()
	}
	if chunk == nil {
		return fmt.Errorf("PCM data not found")
	}
	d.pcmDataAccessed = true

	return nil
}

// WasPCMAccessed returns positively if the PCM data was previously accessed.
func (d *Decoder) WasPCMAccessed() bool {
	if d == nil {
		return false
	}
	return d.pcmDataAccessed
}

// FullPCMBuffer is an inefficient way to access all the PCM data contained in the
// audio container. The entire PCM data is held in memory.
// Consider using PCMBuffer() instead.
func (d *Decoder) FullPCMBuffer() (*audio.IntBuffer, error) {
	if !d.WasPCMAccessed() {
		err := d.FwdToPCM()
		if err != nil {
			return nil, d.err
		}
	}
	if d.PCMChunk == nil {
		return nil, errors.New("PCM chunk not found")
	}
	format := &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
	}

	buf := &audio.IntBuffer{Data: make([]int, 4096), Format: format, SourceBitDepth: int(d.BitDepth)}
	bytesPerSample := (d.BitDepth-1)/8 + 1
	sampleBufData := make([]byte, bytesPerSample)
	decodeF, err := sampleDecodeFunc(int(d.BitDepth))
	if err != nil {
		return nil, fmt.Errorf("could not get sample decode func %v", err)
	}

	i := 0
	for err == nil {
		buf.Data[i], err = decodeF(d.PCMChunk, sampleBufData)
		if err != nil {
			break
		}
		i++
		// grow the underlying slice if needed
		if i == len(buf.Data) {
			buf.Data = append(buf.Data, make([]int, 4096)...)
		}
	}
	buf.Data = buf.Data[:i]

	if err == io.EOF {
		err = nil
	}

	return buf, err
}

// PCMBuffer populates the passed PCM buffer
func (d *Decoder) PCMBuffer(buf *audio.IntBuffer) (n int, err error) {
	if buf == nil {
		return 0, nil
	}

	if !d.pcmDataAccessed {
		err := d.FwdToPCM()
		if err != nil {
			return 0, d.err
		}
	}
	if d.PCMChunk == nil {
		return 0, ErrPCMChunkNotFound
	}

	format := &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
	}

	buf.SourceBitDepth = int(d.BitDepth)
	decodeF, err := sampleDecodeFunc(int(d.BitDepth))
	if err != nil {
		return 0, fmt.Errorf("could not get sample decode func %v", err)
	}

	bPerSample := bytesPerSample(int(d.BitDepth))
	// populate a file buffer to avoid multiple very small reads
	// we need to cap the buffer size to not be bigger than the pcm chunk.
	size := len(buf.Data) * bPerSample
	tmpBuf := make([]byte, size)
	var m int
	m, err = d.PCMChunk.R.Read(tmpBuf)
	if err != nil {
		if err == io.EOF {
			return m, nil
		}
		return m, err
	}
	if m == 0 {
		return m, nil
	}
	bufR := bytes.NewReader(tmpBuf[:m])
	sampleBuf := make([]byte, bPerSample, bPerSample)
	var misaligned bool
	if m%bPerSample > 0 {
		misaligned = true
	}

	// Note that we populate the buffer even if the
	// size of the buffer doesn't fit an even number of frames.
	for n = 0; n < len(buf.Data); n++ {
		buf.Data[n], err = decodeF(bufR, sampleBuf)
		if err != nil {
			// the last sample isn't a full sample but just padding.
			if misaligned {
				n--
			}
			break
		}
	}
	buf.Format = format
	if err == io.EOF {
		err = nil
	}

	return n, err
}

// Format returns the audio format of the decoded content.
func (d *Decoder) Format() *audio.Format {
	if d == nil {
		return nil
	}
	return &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
	}
}

// NextChunk returns the next available chunk
func (d *Decoder) NextChunk() (*riff.Chunk, error) {
	if d.err = d.readHeaders(); d.err != nil {
		d.err = fmt.Errorf("failed to read header - %v", d.err)
		return nil, d.err
	}

	var (
		id   [4]byte
		size uint32
	)

	id, size, d.err = d.parser.IDnSize()
	if d.err != nil {
		d.err = fmt.Errorf("error reading chunk header - %v", d.err)
		return nil, d.err
	}

	// TODO: any reason we don't use d.parser.NextChunk (riff.NextChunk) here?
	// It correctly handles the misaligned chunk.

	// TODO: copied over from riff.parser.NextChunk
	// all RIFF chunks (including WAVE "data" chunks) must be word aligned.
	// If the data uses an odd number of bytes, a padding byte with a value of zero must be placed at the end of the sample data.
	// The "data" chunk header's size should not include this byte.
	if size%2 == 1 {
		size++
	}

	c := &riff.Chunk{
		ID:   id,
		Size: int(size),
		R:    io.LimitReader(d.r, int64(size)),
	}
	return c, d.err
}

// Duration returns the time duration for the current audio container
func (d *Decoder) Duration() (time.Duration, error) {
	if d == nil || d.parser == nil {
		return 0, errors.New("can't calculate the duration of a nil pointer")
	}
	return d.parser.Duration()
}

// String implements the Stringer interface.
func (d *Decoder) String() string {
	return d.parser.String()
}

// readHeaders is safe to call multiple times
func (d *Decoder) readHeaders() error {
	if d == nil || d.NumChans > 0 {
		return nil
	}

	id, size, err := d.parser.IDnSize()
	if err != nil {
		return err
	}
	d.parser.ID = id
	if d.parser.ID != riff.RiffID {
		return fmt.Errorf("%s - %s", d.parser.ID, riff.ErrFmtNotSupported)
	}
	d.parser.Size = size
	if err := binary.Read(d.r, binary.BigEndian, &d.parser.Format); err != nil {
		return err
	}

	var chunk *riff.Chunk
	var rewindBytes int64

	for err == nil {
		chunk, err = d.parser.NextChunk()
		if err != nil {
			break
		}

		if chunk.ID == riff.FmtID {
			chunk.DecodeWavHeader(d.parser)
			d.NumChans = d.parser.NumChannels
			d.BitDepth = d.parser.BitsPerSample
			d.SampleRate = d.parser.SampleRate
			d.WavAudioFormat = d.parser.WavAudioFormat
			d.AvgBytesPerSec = d.parser.AvgBytesPerSec

			if rewindBytes > 0 {
				d.r.Seek(-(rewindBytes + int64(chunk.Size) + 8), 1)
			}
			break
		} else if chunk.ID == CIDList {
			// The list chunk can be in the header or footer
			// because so many players don't support that chunk properly
			// it is recommended to have it at the end of the file.
			DecodeListChunk(d, chunk)
			// unexpected chunk order, might be a bext chunk
			rewindBytes += int64(chunk.Size) + 8
		} else if chunk.ID == CIDSmpl {
			DecodeSamplerChunk(d, chunk)
			rewindBytes += int64(chunk.Size) + 8
		} else {
			// unexpected chunk order, might be a bext chunk
			rewindBytes += int64(chunk.Size) + 8
			// drain the chunk
			io.CopyN(ioutil.Discard, d.r, int64(chunk.Size))
		}
	}

	return d.err
}

func bytesPerSample(bitDepth int) int {
	return bitDepth / 8
}

// sampleDecodeFunc returns a function that can be used to convert
// a byte range into an int value based on the amount of bits used per sample.
// Note that 8bit samples are unsigned, all other values are signed.
func sampleDecodeFunc(bitsPerSample int) (func(io.Reader, []byte) (int, error), error) {
	// NOTE: WAV PCM data is stored using little-endian
	switch bitsPerSample {
	case 8:
		// 8bit values are unsigned
		return func(r io.Reader, buf []byte) (int, error) {
			_, err := r.Read(buf[:1])
			return int(buf[0]), err
		}, nil
	case 16:
		return func(r io.Reader, buf []byte) (int, error) {
			_, err := r.Read(buf[:2])
			return int(int16(binary.LittleEndian.Uint16(buf[:2]))), err
		}, nil
	case 24:
		// -34,359,738,367 (0x7FFFFF) to 34,359,738,368	(0x800000)
		return func(r io.Reader, buf []byte) (int, error) {
			_, err := r.Read(buf[:3])
			if err != nil {
				return 0, err
			}
			return int(audio.Int24LETo32(buf[:3])), nil
		}, nil
	case 32:
		return func(r io.Reader, buf []byte) (int, error) {
			_, err := r.Read(buf[:4])
			return int(int32(binary.LittleEndian.Uint32(buf[:4]))), err
		}, nil
	default:
		return nil, fmt.Errorf("unhandled byte depth:%d", bitsPerSample)
	}
}

// sampleDecodeFloat64Func returns a function that can be used to convert
// a byte range into a float64 value based on the amount of bits used per sample.
func sampleFloat64DecodeFunc(bitsPerSample int) (func([]byte) float64, error) {
	bytesPerSample := bitsPerSample / 8
	switch bytesPerSample {
	case 1:
		// 8bit values are unsigned
		return func(s []byte) float64 {
			return float64(uint8(s[0]))
		}, nil
	case 2:
		return func(s []byte) float64 {
			return float64(int(s[0]) + int(s[1])<<8)
		}, nil
	case 3:
		return func(s []byte) float64 {
			var output int32
			output |= int32(s[2]) << 0
			output |= int32(s[1]) << 8
			output |= int32(s[0]) << 16
			return float64(output)
		}, nil
	case 4:
		// TODO: fix the float64 conversion (current int implementation)
		return func(s []byte) float64 {
			return float64(int(s[0]) + int(s[1])<<8 + int(s[2])<<16 + int(s[3])<<24)
		}, nil
	default:
		return nil, fmt.Errorf("unhandled byte depth:%d", bitsPerSample)
	}
}
