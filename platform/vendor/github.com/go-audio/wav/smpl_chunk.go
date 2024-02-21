package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/go-audio/riff"
)

// smpl chunk is documented here:
// https://sites.google.com/site/musicgapi/technical-documents/wav-file-format#smpl

// DecodeSamplerChunk decodes a smpl chunk and put the data in Decoder.Metadata.SamplerInfo
func DecodeSamplerChunk(d *Decoder, ch *riff.Chunk) error {
	if ch == nil {
		return fmt.Errorf("can't decode a nil chunk")
	}
	if d == nil {
		return fmt.Errorf("nil decoder")
	}
	if ch.ID == CIDSmpl {
		// read the entire chunk in memory
		buf := make([]byte, ch.Size)
		var err error
		if _, err = ch.Read(buf); err != nil {
			return fmt.Errorf("failed to read the smpl chunk - %v", err)
		}
		if d.Metadata == nil {
			d.Metadata = &Metadata{}
		}

		d.Metadata.SamplerInfo = &SamplerInfo{}

		r := bytes.NewReader(buf)

		scratch := make([]byte, 4)
		if _, err = r.Read(scratch); err != nil {
			return fmt.Errorf("failed to read the smpl Manufacturer")
		}
		copy(d.Metadata.SamplerInfo.Manufacturer[:], scratch[:4])
		if _, err = r.Read(scratch); err != nil {
			return fmt.Errorf("failed to read the smpl Product")
		}
		copy(d.Metadata.SamplerInfo.Product[:], scratch[:4])

		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.SamplePeriod); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.MIDIUnityNote); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.MIDIPitchFraction); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.SMPTEFormat); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.SMPTEOffset); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &d.Metadata.SamplerInfo.NumSampleLoops); err != nil {
			return err
		}
		var remaining uint32
		// sampler data
		if err := binary.Read(r, binary.BigEndian, &remaining); err != nil {
			return err
		}
		if d.Metadata.SamplerInfo.NumSampleLoops > 0 {
			d.Metadata.SamplerInfo.Loops = []*SampleLoop{}
			for i := uint32(0); i < d.Metadata.SamplerInfo.NumSampleLoops; i++ {
				sl := &SampleLoop{}
				if _, err = r.Read(scratch); err != nil {
					return fmt.Errorf("failed to read the sample loop cue point id")
				}
				copy(sl.CuePointID[:], scratch[:4])
				if err := binary.Read(r, binary.LittleEndian, &sl.Type); err != nil {
					return err
				}
				if err := binary.Read(r, binary.LittleEndian, &sl.Start); err != nil {
					return err
				}
				if err := binary.Read(r, binary.LittleEndian, &sl.End); err != nil {
					return err
				}
				if err := binary.Read(r, binary.LittleEndian, &sl.Fraction); err != nil {
					return err
				}
				if err := binary.Read(r, binary.LittleEndian, &sl.PlayCount); err != nil {
					return err
				}

				d.Metadata.SamplerInfo.Loops = append(d.Metadata.SamplerInfo.Loops, sl)
			}
		}
	}
	ch.Drain()
	return nil
}
