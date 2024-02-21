package audio

import (
	"errors"
)

var (
	// ErrInvalidBuffer is a generic error returned when trying to read/write to an invalid buffer.
	ErrInvalidBuffer = errors.New("invalid buffer")
)

// Format is a high level representation of the underlying data.
type Format struct {
	// NumChannels is the number of channels contained in the data
	NumChannels int
	// SampleRate is the sampling rate in Hz
	SampleRate int
}

// Buffer is the representation of an audio buffer.
type Buffer interface {
	// PCMFormat is the format of buffer (describing the buffer content/format).
	PCMFormat() *Format
	// NumFrames returns the number of frames contained in the buffer.
	NumFrames() int
	// AsFloatBuffer returns a float 64 buffer from this buffer.
	AsFloatBuffer() *FloatBuffer
	// AsFloat32Buffer returns a float 32 buffer from this buffer.
	AsFloat32Buffer() *Float32Buffer
	// AsIntBuffer returns an int buffer from this buffer.
	AsIntBuffer() *IntBuffer
	// Clone creates a clean clone that can be modified without
	// changing the source buffer.
	Clone() Buffer
}
