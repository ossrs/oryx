package audio

var _ Buffer = (*FloatBuffer)(nil)
var _ Buffer = (*Float32Buffer)(nil)

// FloatBuffer is an audio buffer with its PCM data formatted as float64.
type FloatBuffer struct {
	// Format is the representation of the underlying data format
	Format *Format
	// Data is the buffer PCM data as floats
	Data []float64
}

// PCMFormat returns the buffer format information.
func (buf *FloatBuffer) PCMFormat() *Format { return buf.Format }

// AsFloatBuffer implements the Buffer interface and returns itself.
func (buf *FloatBuffer) AsFloatBuffer() *FloatBuffer { return buf }

// AsFloat32Buffer implements the Buffer interface and returns a float 32 version of itself.
func (buf *FloatBuffer) AsFloat32Buffer() *Float32Buffer {
	newB := &Float32Buffer{}
	newB.Data = make([]float32, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = float32(buf.Data[i])
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// AsIntBuffer returns a copy of this buffer but with data truncated to Ints.
func (buf *FloatBuffer) AsIntBuffer() *IntBuffer {
	newB := &IntBuffer{}
	newB.Data = make([]int, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = int(buf.Data[i])
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// Clone creates a clean clone that can be modified without
// changing the source buffer.
func (buf *FloatBuffer) Clone() Buffer {
	if buf == nil {
		return nil
	}
	newB := &FloatBuffer{}
	newB.Data = make([]float64, len(buf.Data))
	copy(newB.Data, buf.Data)
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// NumFrames returns the number of frames contained in the buffer.
func (buf *FloatBuffer) NumFrames() int {
	if buf == nil || buf.Format == nil {
		return 0
	}
	numChannels := buf.Format.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}

	return len(buf.Data) / numChannels
}

// Float32Buffer is an audio buffer with its PCM data formatted as float32.
type Float32Buffer struct {
	// Format is the representation of the underlying data format
	Format *Format
	// Data is the buffer PCM data as floats
	Data []float32
	// SourceBitDepth helps us know if the source was encoded on
	// 8, 16, 24, 32, 64 bits.
	SourceBitDepth int
}

// PCMFormat returns the buffer format information.
func (buf *Float32Buffer) PCMFormat() *Format { return buf.Format }

// AsFloatBuffer implements the Buffer interface and returns a float64 version of itself.
func (buf *Float32Buffer) AsFloatBuffer() *FloatBuffer {
	newB := &FloatBuffer{}
	newB.Data = make([]float64, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = float64(buf.Data[i])
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// AsFloat32Buffer implements the Buffer interface and returns itself.
func (buf *Float32Buffer) AsFloat32Buffer() *Float32Buffer { return buf }

// AsIntBuffer returns a copy of this buffer but with data truncated to Ints.
// It is usually recommended to apply a transforms when going from a 24bit source
// to an int (16bit destination). Look at transforms.PCMScaleF32() for instance
func (buf *Float32Buffer) AsIntBuffer() *IntBuffer {
	newB := &IntBuffer{SourceBitDepth: buf.SourceBitDepth}
	if newB.SourceBitDepth == 0 {
		newB.SourceBitDepth = 16
	}
	newB.Data = make([]int, len(buf.Data))
	// TODO: we might want to consider checking the min/max values
	// and if we are in a normalized float range, apply a denormalization.
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = int(buf.Data[i])
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// Clone creates a clean clone that can be modified without
// changing the source buffer.
func (buf *Float32Buffer) Clone() Buffer {
	if buf == nil {
		return nil
	}
	newB := &Float32Buffer{}
	newB.Data = make([]float32, len(buf.Data))
	copy(newB.Data, buf.Data)
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// NumFrames returns the number of frames contained in the buffer.
func (buf *Float32Buffer) NumFrames() int {
	if buf == nil || buf.Format == nil {
		return 0
	}
	numChannels := buf.Format.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}

	return len(buf.Data) / numChannels
}
