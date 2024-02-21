package audio

import "math"

// PCMDataFormat is an enum type to indicate the underlying data format used.
type PCMDataFormat uint8

const (
	// DataTypeUnknown refers to an unknown format
	DataTypeUnknown PCMDataFormat = iota
	// DataTypeI8 indicates that the content of the audio buffer made of 8-bit integers.
	DataTypeI8
	// DataTypeI16 indicates that the content of the audio buffer made of 16-bit integers.
	DataTypeI16
	// DataTypeI32 indicates that the content of the audio buffer made of 32-bit integers.
	DataTypeI32
	// DataTypeF32 indicates that the content of the audio buffer made of 32-bit floats.
	DataTypeF32
	// DataTypeF64 indicates that the content of the audio buffer made of 64-bit floats.
	DataTypeF64
)

var _ Buffer = (*PCMBuffer)(nil)

// PCMBuffer encapsulates uncompressed audio data
// and provides useful methods to read/manipulate this PCM data.
// It's a more flexible buffer type allowing the developer to handle
// different kind of buffer data formats and convert between underlying
// types.
type PCMBuffer struct {
	// Format describes the format of the buffer data.
	Format *Format
	// I8 is a store for audio sample data as integers.
	I8 []int8
	// I16 is a store for audio sample data as integers.
	I16 []int16
	// I32 is a store for audio sample data as integers.
	I32 []int32
	// F32 is a store for audio samples data as float64.
	F32 []float32
	// F64 is a store for audio samples data as float64.
	F64 []float64
	// DataType indicates the primary format used for the underlying data.
	// The consumer of the buffer might want to look at this value to know what store
	// to use to optimaly retrieve data.
	DataType PCMDataFormat
	// SourceBitDepth helps us know if the source was encoded on
	// 1 (int8), 2 (int16), 3(int24), 4(int32), 8(int64) bytes.
	SourceBitDepth uint8
}

// Len returns the length of the underlying data.
func (b *PCMBuffer) Len() int {
	if b == nil {
		return 0
	}

	switch b.DataType {
	case DataTypeI8:
		return len(b.I8)
	case DataTypeI16:
		return len(b.I16)
	case DataTypeI32:
		return len(b.I32)
	case DataTypeF32:
		return len(b.F32)
	case DataTypeF64:
		return len(b.F64)
	default:
		return 0
	}
}

// PCMFormat returns the buffer format information.
func (b *PCMBuffer) PCMFormat() *Format {
	if b == nil {
		return nil
	}
	return b.Format
}

// NumFrames returns the number of frames contained in the buffer.
func (b *PCMBuffer) NumFrames() int {
	if b == nil || b.Format == nil {
		return 0
	}
	numChannels := b.Format.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}

	return b.Len() / numChannels
}

// AsFloatBuffer returns a copy of this buffer but with data converted to floats.
func (b *PCMBuffer) AsFloatBuffer() *FloatBuffer {
	newB := &FloatBuffer{}
	newB.Data = b.AsF64()
	if b.Format != nil {
		newB.Format = &Format{
			NumChannels: b.Format.NumChannels,
			SampleRate:  b.Format.SampleRate,
		}
	}
	return newB
}

// AsFloat32Buffer implements the Buffer interface and returns a float 32 version of itself.
func (b *PCMBuffer) AsFloat32Buffer() *Float32Buffer {
	newB := &Float32Buffer{}
	newB.Data = b.AsF32()
	if b.Format != nil {
		newB.Format = &Format{
			NumChannels: b.Format.NumChannels,
			SampleRate:  b.Format.SampleRate,
		}
	}
	return newB
}

// AsIntBuffer returns a copy of this buffer but with data truncated to Ints.
func (b *PCMBuffer) AsIntBuffer() *IntBuffer {
	newB := &IntBuffer{}
	newB.Data = b.AsInt()
	if b.Format != nil {
		newB.Format = &Format{
			NumChannels: b.Format.NumChannels,
			SampleRate:  b.Format.SampleRate,
		}
	}
	return newB
}

// AsI8 returns the buffer's samples as int8 sample values.
// If the buffer isn't in this format, a copy is created and converted.
// Note that converting might result in loss of resolution.
func (b *PCMBuffer) AsI8() (out []int8) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case DataTypeI8:
		return b.I8
	case DataTypeI16:
		out = make([]int8, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = int8(b.I16[i])
		}
	case DataTypeI32:
		out = make([]int8, len(b.I32))
		for i := 0; i < len(b.I32); i++ {
			out[i] = int8(b.I32[i])
		}
	case DataTypeF32:
		out = make([]int8, len(b.F32))
		for i := 0; i < len(b.F32); i++ {
			out[i] = int8(b.F32[i])
		}
	case DataTypeF64:
		out = make([]int8, len(b.F64))
		for i := 0; i < len(b.F64); i++ {
			out[i] = int8(b.F64[i])
		}
	}
	return out
}

// AsI16 returns the buffer's samples as int16 sample values.
// If the buffer isn't in this format, a copy is created and converted.
// Note that converting might result in loss of resolution.
func (b *PCMBuffer) AsI16() (out []int16) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case DataTypeI8:
		out = make([]int16, len(b.I8))
		for i := 0; i < len(b.I8); i++ {
			out[i] = int16(b.I8[i])
		}
	case DataTypeI16:
		return b.I16
	case DataTypeI32:
		out = make([]int16, len(b.I32))
		for i := 0; i < len(b.I32); i++ {
			out[i] = int16(b.I32[i])
		}
	case DataTypeF32:
		out = make([]int16, len(b.F32))
		for i := 0; i < len(b.F32); i++ {
			out[i] = int16(b.F32[i])
		}
	case DataTypeF64:
		out = make([]int16, len(b.F64))
		for i := 0; i < len(b.F64); i++ {
			out[i] = int16(b.F64[i])
		}
	}
	return out
}

// AsI32 returns the buffer's samples as int32 sample values.
// If the buffer isn't in this format, a copy is created and converted.
// Note that converting a float to an int might result in unexpected truncations.
func (b *PCMBuffer) AsI32() (out []int32) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case DataTypeI8:
		out = make([]int32, len(b.I8))
		for i := 0; i < len(b.I8); i++ {
			out[i] = int32(b.I8[i])
		}
	case DataTypeI16:
		out = make([]int32, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = int32(b.I16[i])
		}
	case DataTypeI32:
		return b.I32
	case DataTypeF32:
		out = make([]int32, len(b.F32))
		for i := 0; i < len(b.F32); i++ {
			out[i] = int32(b.F32[i])
		}
	case DataTypeF64:
		out = make([]int32, len(b.F64))
		for i := 0; i < len(b.F64); i++ {
			out[i] = int32(b.F64[i])
		}
	}
	return out
}

// AsInt returns the buffer content as integers (int32s).
// It's recommended to avoid this method since it creates
// an extra copy of the buffer content.
func (b *PCMBuffer) AsInt() (out []int) {
	int32s := b.AsI32()
	out = make([]int, len(int32s))
	for i := 0; i < len(int32s); i++ {
		out[i] = int(int32s[i])
	}
	return out
}

// AsF32 returns the buffer's samples as float32 sample values.
// If the buffer isn't in this format, a copy is created and converted.
// Note that converting might result in unexpected truncations.
func (b *PCMBuffer) AsF32() (out []float32) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case DataTypeI8:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float32, len(b.I8))
		for i := 0; i < len(b.I8); i++ {
			out[i] = float32(float64(int64(b.I8[i])) / factor)
		}
	case DataTypeI16:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float32, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = float32(float64(int64(b.I16[i])) / factor)
		}
	case DataTypeI32:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float32, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = float32(float64(int64(b.I16[i])) / factor)
		}
	case DataTypeF32:
		return b.F32
	case DataTypeF64:
		out = make([]float32, len(b.F64))
		for i := 0; i < len(b.F64); i++ {
			out[i] = float32(b.F64[i])
		}
	}
	return out
}

// AsF64 returns the buffer's samples as float64 sample values.
// If the buffer isn't in this format, a copy is created and converted.
// Note that converting might result in unexpected truncations.
func (b *PCMBuffer) AsF64() (out []float64) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case DataTypeI8:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float64, len(b.I8))
		for i := 0; i < len(b.I8); i++ {
			out[i] = float64(int64(b.I8[i])) / factor
		}
	case DataTypeI16:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float64, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = float64(int64(b.I16[i])) / factor
		}
	case DataTypeI32:
		bitDepth := b.calculateIntBitDepth()
		factor := math.Pow(2, 8*float64(bitDepth/8)-1)
		out = make([]float64, len(b.I16))
		for i := 0; i < len(b.I16); i++ {
			out[i] = float64(int64(b.I16[i])) / factor
		}
	case DataTypeF32:
		out = make([]float64, len(b.F32))
		for i := 0; i < len(b.F32); i++ {
			out[i] = float64(b.F32[i])
		}
	case DataTypeF64:
		return b.F64
	}
	return out
}

// Clone creates a clean clone that can be modified without
// changing the source buffer.
func (b *PCMBuffer) Clone() Buffer {
	if b == nil {
		return nil
	}
	newB := &PCMBuffer{DataType: b.DataType}
	switch b.DataType {
	case DataTypeI8:
		newB.I8 = make([]int8, len(b.I8))
		copy(newB.I8, b.I8)
	case DataTypeI16:
		newB.I16 = make([]int16, len(b.I16))
		copy(newB.I16, b.I16)
	case DataTypeI32:
		newB.I32 = make([]int32, len(b.I32))
		copy(newB.I32, b.I32)
	case DataTypeF32:
		newB.F32 = make([]float32, len(b.F32))
		copy(newB.F32, b.F32)
	case DataTypeF64:
		newB.F64 = make([]float64, len(b.F64))
		copy(newB.F64, b.F64)
	}

	newB.Format = &Format{
		NumChannels: b.Format.NumChannels,
		SampleRate:  b.Format.SampleRate,
	}
	return newB
}

// SwitchPrimaryType is a convenience method to switch the primary data type.
// Use this if you process/swap a different type than the original type.
// Notes that conversion might be lossy if you switch to a lower resolution format.
func (b *PCMBuffer) SwitchPrimaryType(t PCMDataFormat) {
	if b == nil || t == b.DataType {
		return
	}
	switch t {
	case DataTypeI8:
		b.I8 = b.AsI8()
		b.I16 = nil
		b.I32 = nil
		b.F32 = nil
		b.F64 = nil
	case DataTypeI16:
		b.I8 = nil
		b.I16 = b.AsI16()
		b.I32 = nil
		b.F32 = nil
		b.F64 = nil
	case DataTypeI32:
		b.I8 = nil
		b.I16 = nil
		b.I32 = b.AsI32()
		b.F32 = nil
		b.F64 = nil
	case DataTypeF32:
		b.I8 = nil
		b.I16 = nil
		b.I32 = nil
		b.F32 = b.AsF32()
		b.F64 = nil
	case DataTypeF64:
		b.I8 = nil
		b.I16 = nil
		b.I32 = nil
		b.F32 = nil
		b.F64 = b.AsF64()
	}

	b.DataType = t
}

// calculateIntBithDepth looks at the int values in the buffer and returns
// the required lowest bit depth.
func (b *PCMBuffer) calculateIntBitDepth() uint8 {
	if b == nil {
		return 0
	}
	bitDepth := b.SourceBitDepth
	if bitDepth != 0 {
		return bitDepth
	}
	var max int64
	switch b.DataType {
	case DataTypeI8:
		var i8max int8
		for _, s := range b.I8 {
			if s > i8max {
				i8max = s
			}
		}
		max = int64(i8max)
	case DataTypeI16:
		var i16max int16
		for _, s := range b.I16 {
			if s > i16max {
				i16max = s
			}
		}
		max = int64(i16max)
	case DataTypeI32:
		var i32max int32
		for _, s := range b.I32 {
			if s > i32max {
				i32max = s
			}
		}
		max = int64(i32max)
	default:
		// This method is only meant to be used on int buffers.
		return bitDepth
	}
	bitDepth = 8
	if max > 127 {
		bitDepth = 16
	}
	// greater than int16, expecting int24
	if max > 32767 {
		bitDepth = 24
	}
	// int 32
	if max > 8388607 {
		bitDepth = 32
	}
	// int 64
	if max > 4294967295 {
		bitDepth = 64
	}

	return bitDepth
}
