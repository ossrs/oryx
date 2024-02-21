package riff

import (
	"errors"
	"io"
	"sync"
	"time"
)

var (
	RiffID = [4]byte{'R', 'I', 'F', 'F'}
	FmtID  = [4]byte{'f', 'm', 't', ' '}
	// To align RIFF chunks to certain boundaries (i.e. 2048bytes for CD-ROMs) the RIFF specification includes a JUNK chunk.
	// Its contents are to be skipped when reading. When writing RIFFs, JUNK chunks should not have odd number as Size.
	junkID      = [4]byte{'J', 'U', 'N', 'K'}
	WavFormatID = [4]byte{'W', 'A', 'V', 'E'}
	// DataFormatID is the Wave Data Chunk ID, it contains the digital audio sample data which can be decoded using the format
	// and compression method specified in the Wave Format Chunk. If the Compression Code is 1 (uncompressed PCM), then the Wave Data contains raw sample values.
	DataFormatID = [4]byte{'d', 'a', 't', 'a'}
	rmiFormatID  = [4]byte{'R', 'M', 'I', 'D'}
	aviFormatID  = [4]byte{'A', 'V', 'I', ' '}
	// ErrFmtNotSupported is a generic error reporting an unknown format.
	ErrFmtNotSupported = errors.New("format not supported")
	// ErrUnexpectedData is a generic error reporting that the parser encountered unexpected data.
	ErrUnexpectedData = errors.New("unexpected data content")
)

// New creates a parser wrapper for a reader.
// Note that the reader doesn't get rewinded as the container is processed.
func New(r io.Reader) *Parser {
	return &Parser{r: r, Wg: &sync.WaitGroup{}}
}

// Duration returns the time duration of the passed reader if the sub format is supported.
func Duration(r io.Reader) (time.Duration, error) {
	c := New(r)
	if err := c.ParseHeaders(); err != nil {
		return 0, err
	}
	return c.Duration()
}
