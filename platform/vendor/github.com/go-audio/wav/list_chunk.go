package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-audio/riff"
)

var (
	// See http://bwfmetaedit.sourceforge.net/listinfo.html
	markerIART    = [4]byte{'I', 'A', 'R', 'T'}
	markerISFT    = [4]byte{'I', 'S', 'F', 'T'}
	markerICRD    = [4]byte{'I', 'C', 'R', 'D'}
	markerICOP    = [4]byte{'I', 'C', 'O', 'P'}
	markerIARL    = [4]byte{'I', 'A', 'R', 'L'}
	markerINAM    = [4]byte{'I', 'N', 'A', 'M'}
	markerIENG    = [4]byte{'I', 'E', 'N', 'G'}
	markerIGNR    = [4]byte{'I', 'G', 'N', 'R'}
	markerIPRD    = [4]byte{'I', 'P', 'R', 'D'}
	markerISRC    = [4]byte{'I', 'S', 'R', 'C'}
	markerISBJ    = [4]byte{'I', 'S', 'B', 'J'}
	markerICMT    = [4]byte{'I', 'C', 'M', 'T'}
	markerITRK    = [4]byte{'I', 'T', 'R', 'K'}
	markerITRKBug = [4]byte{'i', 't', 'r', 'k'}
	markerITCH    = [4]byte{'I', 'T', 'C', 'H'}
	markerIKEY    = [4]byte{'I', 'K', 'E', 'Y'}
	markerIMED    = [4]byte{'I', 'M', 'E', 'D'}
)

// DecodeListChunk decodes a LIST chunk
func DecodeListChunk(d *Decoder, ch *riff.Chunk) error {
	if ch == nil {
		return fmt.Errorf("can't decode a nil chunk")
	}
	if d == nil {
		return fmt.Errorf("nil decoder")
	}
	if ch.ID == CIDList {
		// read the entire chunk in memory
		buf := make([]byte, ch.Size)
		var err error
		if _, err = ch.Read(buf); err != nil {
			return fmt.Errorf("failed to read the LIST chunk - %v", err)
		}
		r := bytes.NewReader(buf)
		// INFO subchunk
		scratch := make([]byte, 4)
		if _, err = r.Read(scratch); err != nil {
			return fmt.Errorf("failed to read the INFO subchunk - %v", err)
		}
		if !bytes.Equal(scratch, CIDInfo[:]) {
			// "expected an INFO subchunk but got %s", string(scratch)
			// TODO: support adtl subchunks
			ch.Drain()
			return nil
		}
		if d.Metadata == nil {
			d.Metadata = &Metadata{}
		}

		// the rest is a list of string entries
		var (
			id   [4]byte
			size uint32
		)
		readSubHeader := func() error {
			if err := binary.Read(r, binary.BigEndian, &id); err != nil {
				return err
			}
			return binary.Read(r, binary.LittleEndian, &size)
		}

		for err == nil {
			err = readSubHeader()
			if err != nil {
				break
			}
			scratch = make([]byte, size)
			r.Read(scratch)
			switch id {
			case markerIARL:
				d.Metadata.Location = nullTermStr(scratch)
			case markerIART:
				d.Metadata.Artist = nullTermStr(scratch)
			case markerISFT:
				d.Metadata.Software = nullTermStr(scratch)
			case markerICRD:
				d.Metadata.CreationDate = nullTermStr(scratch)
			case markerICOP:
				d.Metadata.Copyright = nullTermStr(scratch)
			case markerINAM:
				d.Metadata.Title = nullTermStr(scratch)
			case markerIENG:
				d.Metadata.Engineer = nullTermStr(scratch)
			case markerIGNR:
				d.Metadata.Genre = nullTermStr(scratch)
			case markerIPRD:
				d.Metadata.Product = nullTermStr(scratch)
			case markerISRC:
				d.Metadata.Source = nullTermStr(scratch)
			case markerISBJ:
				d.Metadata.Subject = nullTermStr(scratch)
			case markerICMT:
				d.Metadata.Comments = nullTermStr(scratch)
			case markerITRK, markerITRKBug:
				d.Metadata.TrackNbr = nullTermStr(scratch)
			case markerITCH:
				d.Metadata.Technician = nullTermStr(scratch)
			case markerIKEY:
				d.Metadata.Keywords = nullTermStr(scratch)
			case markerIMED:
				d.Metadata.Medium = nullTermStr(scratch)
			}

			if size%2 == 1 {
				r.Seek(1, io.SeekCurrent)
			}
		}
	}
	ch.Drain()
	return nil
}

func encodeInfoChunk(e *Encoder) []byte {
	if e == nil || e.Metadata == nil {
		return nil
	}
	buf := bytes.NewBuffer(nil)

	writeSection := func(id [4]byte, val string) {
		buf.Write(id[:])
		binary.Write(buf, binary.LittleEndian, uint32(len(val)+1))
		buf.Write(append([]byte(val), 0x00))
	}
	if e.Metadata.Artist != "" {
		writeSection(markerIART, e.Metadata.Artist)
	}
	if e.Metadata.Comments != "" {
		writeSection(markerICMT, e.Metadata.Comments)
	}
	if e.Metadata.Copyright != "" {
		writeSection(markerICOP, e.Metadata.Copyright)
	}
	if e.Metadata.CreationDate != "" {
		writeSection(markerICRD, e.Metadata.CreationDate)
	}
	if e.Metadata.Engineer != "" {
		writeSection(markerIENG, e.Metadata.Engineer)
	}
	if e.Metadata.Technician != "" {
		writeSection(markerITCH, e.Metadata.Technician)
	}
	if e.Metadata.Genre != "" {
		writeSection(markerIGNR, e.Metadata.Genre)
	}
	if e.Metadata.Keywords != "" {
		writeSection(markerIKEY, e.Metadata.Keywords)
	}
	if e.Metadata.Medium != "" {
		writeSection(markerIMED, e.Metadata.Medium)
	}
	if e.Metadata.Title != "" {
		writeSection(markerINAM, e.Metadata.Title)
	}
	if e.Metadata.Product != "" {
		writeSection(markerIPRD, e.Metadata.Product)
	}
	if e.Metadata.Subject != "" {
		writeSection(markerISBJ, e.Metadata.Subject)
	}
	if e.Metadata.Software != "" {
		writeSection(markerISFT, e.Metadata.Software)
	}
	if e.Metadata.Source != "" {
		writeSection(markerISRC, e.Metadata.Source)
	}
	if e.Metadata.Location != "" {
		writeSection(markerIARL, e.Metadata.Location)
	}
	if e.Metadata.TrackNbr != "" {
		writeSection(markerITRK, e.Metadata.TrackNbr)
	}

	return append(CIDInfo, buf.Bytes()...)
}
