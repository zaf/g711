/*
	Copyright (C) 2016 - 2017, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.
*/

/*
	Package g711 implements encoding and decoding of G711.0 compressed sound data.
	G.711 is an ITU-T standard for audio companding.


*/

package g711

import (
	"bytes"
	"errors"
	"io"
)

const (
	// Input and output formats
	Alaw = iota // Alaw G711 encoded PCM data
	Ulaw        // Ulaw G711  encoded PCM data
	Lpcm        // Lpcm 16bit signed linear data
)

type encoder func(int16) uint8
type decoder func(uint8) int16
type transcoder func(uint8) uint8

// Reader reads G711 PCM data and decodes it to 16bit LPCM or directly transcodes between A-law and u-law
type Reader struct {
	input  int           // source format
	output int           // output format
	r      io.Reader     // source data
	buf    *bytes.Buffer // local buffer
}

// Writer encodes 16bit LPCM data to G711 PCM or directly transcodes between A-law and u-law
type Writer struct {
	input  int           // source format
	output int           // output format
	w      io.Writer     // output data
	buf    *bytes.Buffer //local buffer
}

// NewAlawDecoder returns a pointer to a Reader that decodes or trans-codes A-law data.
// It takes as input the source data Reader and the output encoding fomrat.
func NewAlawDecoder(reader io.Reader, output int) (*Reader, error) {
	if output != Ulaw && output != Lpcm {
		return nil, errors.New("Invalid output format")
	}
	b := new(bytes.Buffer)
	return &Reader{input: Alaw, output: output, r: reader, buf: b}, nil
}

// NewUlawDecoder returns a pointer to a Reader that decodes or trans-codes u-law data.
// It takes as input the source data Reader and the output encoding fomrat.
func NewUlawDecoder(reader io.Reader, output int) (*Reader, error) {
	if output != Alaw && output != Lpcm {
		return nil, errors.New("Invalid output format")
	}
	b := new(bytes.Buffer)
	return &Reader{input: Ulaw, output: output, r: reader, buf: b}, nil
}

// NewAlawEncoder returns a pointer to a Writer that encodes data to A-law.
// It takes as input the destination data Writer and the input encoding fomrat.
func NewAlawEncoder(writer io.Writer, input int) (*Writer, error) {
	if input != Ulaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	b := new(bytes.Buffer)
	return &Writer{input: input, output: Alaw, w: writer, buf: b}, nil
}

// NewUlawEncoder returns a pointer to a Writer that encodes data to u-law.
// It takes as input the destination data Writer and the input encoding fomrat.
func NewUlawEncoder(writer io.Writer, input int) (*Writer, error) {
	if input != Alaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	b := new(bytes.Buffer)
	return &Writer{input: input, output: Ulaw, w: writer, buf: b}, nil
}

// Reset discards the Reader state. This permits reusing a Reader rather than allocating a new one.
func (r *Reader) Reset(reader io.Reader) {
	r.buf.Reset()
	r.r = reader
}

// Reset discards the Writer state. This permits reusing a Writer rather than allocating a new one.
func (w *Writer) Reset(writer io.Writer) {
	w.buf.Reset()
	w.w = writer
}

// Read decodes G711 data. Reads up to len(p) bytes into p, returns the number
// of bytes read and any error encountered.
func (r *Reader) Read(p []byte) (int, error) {
	var dec decoder
	var tr transcoder
	if r.input == Alaw {
		dec = DecodeAlaw
		tr = Alaw2Ulaw
	} else {
		dec = DecodeUlaw
		tr = Ulaw2Alaw
	}
	b := make([]byte, 4096)
	n, rErr := r.r.Read(b)
	var wrErr error
	for _, data := range b[0:n] {
		if r.output == Lpcm {
			lpcm := dec(data) // Decode G711 data to LPCM
			_, wrErr = r.buf.Write([]byte{byte(lpcm), byte(lpcm >> 8)})
		} else {
			wrErr = r.buf.WriteByte(tr(data)) // Trans-code
		}
		if wrErr != nil {
			break
		}
	}
	i, err := r.buf.Read(p)
	if err == nil {
		if wrErr != nil {
			err = wrErr
		} else {
			err = rErr
		}
	}
	return i, err
}

// Write encodes G711 Data. Writes len(p) bytes from p to the underlying data stream,
// returns the number of bytes written from p (0 <= n <= len(p)/2 due to compression)
// and any error encountered that caused the write to stop early.
func (w *Writer) Write(p []byte) (int, error) {
	var err, wrErr error
	var enc encoder
	var tr transcoder
	if w.output == Alaw {
		enc = EncodeAlaw
		tr = Ulaw2Alaw
	} else {
		enc = EncodeUlaw
		tr = Alaw2Ulaw
	}
	if w.input == Lpcm { // Encode LPCM data to G711
		for i := 0; i <= len(p)-2; i = i + 2 {
			wrErr = w.buf.WriteByte(enc(int16(p[i]) | int16(p[i+1])<<8))
			if wrErr != nil {
				break
			}
		}
	} else { // Trans-code
		for _, data := range p {
			wrErr = w.buf.WriteByte(tr(data))
			if wrErr != nil {
				break
			}
		}
	}
	i, err := w.w.Write(w.buf.Bytes())
	if err == nil {
		err = wrErr
	}
	return i, err
}

// Flush flushes any pending data to the underlying writer.
func (w *Writer) Flush() (err error) {
	if w.buf.Len() > 0 {
		_, err = w.buf.WriteTo(w.w)
	}
	return
}
