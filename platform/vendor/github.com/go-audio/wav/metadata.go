package wav

// Metadata represents optional metadata added to the wav file.
type Metadata struct {
	SamplerInfo *SamplerInfo
	// Artist of the original subject of the file. For example, Michaelangelo.
	Artist string
	// Comments provides general comments about the file or the subject of the
	// file. If the comment is several sentences long, end each sentence with a
	// period. Do not include newline characters.
	Comments string
	// Copyright records the copyright information for the file.
	Copyright string
	// CreationDate specifies the date the subject of the file was created. List
	// dates in year-month-day format, padding one-digit months and days with a
	// zero on the left. For example: 1553-05-03 for May 3, 1553. The year
	// should always be given using four digits.
	CreationDate string
	// Engineer stores the name of the engineer who worked on the file. If there
	// are multiple engineers, separate the names by a semicolon and a blank.
	// For example: Smith, John; Adams, Joe.
	Engineer string
	// Technician identifies the technician who sampled the subject file. For
	// example: Smith, John.
	Technician string
	// Genre describes the original work, such as jazz, classical, rock, etc.
	Genre string
	// Keywords provides a list of keywords that refer to the file or subject of
	// the file. Separate multiple keywords with a semicolon and a blank. For
	// example, Seattle; zoology; The Civil War.
	Keywords string
	// Medium describes the original subject of the file, such as record, CD and so forth.
	Medium string
	// Title stores the title of the subject of the file, such as bohemian rhapsody.
	Title string
	// Product AKA album specifies the name of the title the file was originally
	// intended for: A Night at the Opera
	Product string
	// Subject describes the contents of the file, such as Metadata Management.
	Subject string
	// Software identifies the name of the software package used to create the
	// file, such as go-audio.
	Software string
	// Source identifies the name of the person or organization who supplied the
	// original subject of the file. For example: Splice.
	Source string
	// Location or Archival Location - Indicates where the subject of the file is archived.
	Location string
	// TrackNbr is the track number
	TrackNbr string
	// CuePoints is a list of cue points in the wav file.
	CuePoints []*CuePoint
}

// SamplerInfo is extra metadata pertinent to a sampler type usage.
type SamplerInfo struct {
	// Manufacturer field specifies the MIDI Manufacturer's Association
	// (MMA) Manufacturer code for the sampler intended to receive this file's
	// waveform. Each manufacturer of a MIDI product is assigned a unique ID
	// which identifies the company. If no particular manufacturer is to be
	// specified, a value of 0 should be used. The value is stored with some
	// extra information to enable translation to the value used in a MIDI
	// System Exclusive transmission to the sampler. The high byte indicates the
	// number of low order bytes (1 or 3) that are valid for the manufacturer
	// code. For example, the value for Digidesign will be 0x01000013 (0x13) and
	// the value for Microsoft will be 0x30000041 (0x00, 0x00, 0x41).
	Manufacturer [4]byte
	// Product field specifies the MIDI model ID defined by the manufacturer
	// corresponding to the Manufacturer field. Contact the manufacturer of the
	// sampler to get the model ID. If no particular manufacturer's product is
	// to be specified, a value of 0 should be used.
	Product [4]byte
	// SamplePeriod The sample period specifies the duration of time that passes
	// during the playback of one sample in nanoseconds (normally equal to 1 /
	// Samplers Per Second, where Samples Per Second is the value found in the
	// format chunk).
	SamplePeriod uint32
	// MIDIUnityNote The MIDI unity note value has the same meaning as the instrument chunk's
	// MIDI Unshifted Note field which specifies the musical note at which the
	// sample will be played at it's original sample rate (the sample rate
	// specified in the format chunk).
	MIDIUnityNote uint32
	// MIDIPitchFraction The MIDI pitch fraction specifies the fraction of a
	// semitone up from the specified MIDI unity note field. A value of
	// 0x80000000 means 1/2 semitone (50 cents) and a value of 0x00000000 means
	// no fine tuning between semitones.
	MIDIPitchFraction uint32
	// SMPTEFormat The SMPTE format specifies the Society of Motion Pictures and
	// Television E time format used in the following SMPTE Offset field. If a
	// value of 0 is set, SMPTE Offset should also be set to 0. (0, 24, 25, 29, 30)
	SMPTEFormat uint32
	// SMPTEOffset The SMPTE Offset value specifies the time offset to be used
	// for the synchronization / calibration to the first sample in the
	// waveform. This value uses a format of 0xhhmmssff where hh is a signed
	// value that specifies the number of hours (-23 to 23), mm is an unsigned
	// value that specifies the number of minutes (0 to 59), ss is an unsigned
	// value that specifies the number of seconds (0 to 59) and ff is an
	// unsigned value that specifies the number of frames (0 to -1).
	SMPTEOffset uint32
	// NumSampleLoops The sample loops field specifies the number Sample Loop
	// definitions in the following list. This value may be set to 0 meaning
	// that no sample loops follow.
	NumSampleLoops uint32
	// Loops A list of sample loops is simply a set of consecutive loop
	// descriptions. The sample loops do not have to be in any particular order
	// because each sample loop associated cue point position is used to
	// determine the play order.
	Loops []*SampleLoop
}

// SampleLoop indicates a loop and its properties within the audio file
type SampleLoop struct {
	// CuePointID - The Cue Point ID specifies the unique ID that corresponds to one of the
	// defined cue points in the cue point list. Furthermore, this ID
	// corresponds to any labels defined in the associated data list chunk which
	// allows text labels to be assigned to the various sample loops.
	CuePointID [4]byte
	// Type - The type field defines how the waveform samples will be looped.
	// 0 Loop forward (normal)
	// 1 Alternating loop (forward/backward, also known as Ping Pong)
	// 2 Loop backward (reverse)
	// 3 Reserved for future standard types
	// 32 - 0xFFFFFFFF Sampler specific types (defined by manufacturer)
	Type uint32
	// Start - The start value specifies the byte offset into the waveform data
	// of the first sample to be played in the loop.
	Start uint32
	// End - The end value specifies the byte offset into the waveform data of
	// the last sample to be played in the loop.
	End uint32
	// Fraction - The fractional value specifies a fraction of a sample at which
	// to loop. This allows a loop to be fine tuned at a resolution greater than
	// one sample. The value can range from 0x00000000 to 0xFFFFFFFF. A value of
	// 0 means no fraction, a value of 0x80000000 means 1/2 of a sample length.
	// 0xFFFFFFFF is the smallest fraction of a sample that can be represented.
	Fraction uint32
	// PlayCount - The play count value determines the number of times to play
	// the loop. A value of 0 specifies an infinite sustain loop. An infinite
	// sustain loop will continue looping until some external force interrupts
	// playback, such as the musician releasing the key that triggered the
	// wave's playback. All other values specify an absolute number of times to
	// loop.
	PlayCount uint32
}
