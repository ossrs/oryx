package audio

import "math"

// IntMaxSignedValue returns the max value of an integer
// based on its memory size
func IntMaxSignedValue(b int) int {
	switch b {
	case 8:
		return 255 / 2
	case 16:
		return 65535 / 2
	case 24:
		return 16777215 / 2
	case 32:
		return 4294967295 / 2
	default:
		return 0
	}
}

// IEEEFloatToInt converts a 10 byte IEEE float into an int.
func IEEEFloatToInt(b [10]byte) int {
	var i uint32
	// Negative number
	if (b[0] & 0x80) == 1 {
		return 0
	}

	// Less than 1
	if b[0] <= 0x3F {
		return 1
	}

	// Too big
	if b[0] > 0x40 {
		return 67108864
	}

	// Still too big
	if b[0] == 0x40 && b[1] > 0x1C {
		return 800000000
	}

	i = (uint32(b[2]) << 23) | (uint32(b[3]) << 15) | (uint32(b[4]) << 7) | (uint32(b[5]) >> 1)
	i >>= (29 - uint32(b[1]))

	return int(i)
}

// IntToIEEEFloat converts an int into a 10 byte IEEE float.
func IntToIEEEFloat(i int) [10]byte {
	b := [10]byte{}
	num := float64(i)

	var sign int
	var expon int
	var fMant, fsMant float64
	var hiMant, loMant uint

	if num < 0 {
		sign = 0x8000
	} else {
		sign = 0
	}

	if num == 0 {
		expon = 0
		hiMant = 0
		loMant = 0
	} else {
		fMant, expon = math.Frexp(num)
		if (expon > 16384) || !(fMant < 1) { /* Infinity or NaN */
			expon = sign | 0x7FFF
			hiMant = 0
			loMant = 0 /* infinity */
		} else { /* Finite */
			expon += 16382
			if expon < 0 { /* denormalized */
				fMant = math.Ldexp(fMant, expon)
				expon = 0
			}
			expon |= sign
			fMant = math.Ldexp(fMant, 32)
			fsMant = math.Floor(fMant)
			hiMant = uint(fsMant)
			fMant = math.Ldexp(fMant-fsMant, 32)
			fsMant = math.Floor(fMant)
			loMant = uint(fsMant)
		}
	}

	b[0] = byte(expon >> 8)
	b[1] = byte(expon)
	b[2] = byte(hiMant >> 24)
	b[3] = byte(hiMant >> 16)
	b[4] = byte(hiMant >> 8)
	b[5] = byte(hiMant)
	b[6] = byte(loMant >> 24)
	b[7] = byte(loMant >> 16)
	b[8] = byte(loMant >> 8)
	b[9] = byte(loMant)

	return b
}

// Uint24to32 converts a 3 byte uint23 into a uint32
// BigEndian!
func Uint24to32(bytes []byte) uint32 {
	var output uint32
	output |= uint32(bytes[2]) << 0
	output |= uint32(bytes[1]) << 8
	output |= uint32(bytes[0]) << 16

	return output
}

// Int24BETo32 converts an int24 value from 3 bytes into an int32 value
func Int24BETo32(bytes []byte) int32 {
	if len(bytes) < 3 {
		return 0
	}
	ss := int32(0xFF&bytes[0])<<16 | int32(0xFF&bytes[1])<<8 | int32(0xFF&bytes[2])
	if (ss & 0x800000) > 0 {
		ss |= ^0xffffff
	}

	return ss
}

// Int24LETo32 converts an int24 value from 3 bytes into an int32 value
func Int24LETo32(bytes []byte) int32 {
	if len(bytes) < 3 {
		return 0
	}
	ss := int32(bytes[0]) | int32(bytes[1])<<8 | int32(bytes[2])<<16
	if (ss & 0x800000) > 0 {
		ss |= ^0xffffff
	}

	return ss
}

// Uint32toUint24Bytes converts a uint32 into a 3 byte uint24 representation
func Uint32toUint24Bytes(n uint32) []byte {
	bytes := make([]byte, 3)
	bytes[0] = byte(n >> 16)
	bytes[1] = byte(n >> 8)
	bytes[2] = byte(n >> 0)

	return bytes
}

// Int32toInt24LEBytes converts an int32 into a little endian 3 byte int24 representation
func Int32toInt24LEBytes(n int32) []byte {
	bytes := make([]byte, 3)
	if (n & 0x800000) > 0 {
		n |= ^0xffffff
	}
	bytes[2] = byte(n >> 16)
	bytes[1] = byte(n >> 8)
	bytes[0] = byte(n >> 0)
	return bytes
}

// Int32toInt24BEBytes converts an int32 into a big endian 3 byte int24 representation
func Int32toInt24BEBytes(n int32) []byte {
	bytes := make([]byte, 3)
	if (n & 0x800000) > 0 {
		n |= ^0xffffff
	}
	bytes[0] = byte(n >> 16)
	bytes[1] = byte(n >> 8)
	bytes[2] = byte(n >> 0)

	return bytes
}
