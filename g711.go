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

// Reader reads G711 PCM data and decodes it to 16bit LPCM or directly transcodes between A-law and u-law
type Reader struct {
	output    int               // output format
	decode    func(uint8) int16 // decoding function
	transcode func(uint8) uint8 // transcoding function
	source    io.Reader         // source data
	buf       *bytes.Buffer     // local buffer
}

// Writer encodes 16bit LPCM data to G711 PCM or directly transcodes between A-law and u-law
type Writer struct {
	input       int               // input format
	encode      func(int16) uint8 // encoding function
	transcode   func(uint8) uint8 // transcoding function
	destination io.Writer         // output data
	buf         *bytes.Buffer     //local buffer
}

// NewAlawDecoder returns a pointer to a Reader that decodes or trans-codes A-law data.
// It takes as input the source data Reader and the output encoding fomrat.
func NewAlawDecoder(reader io.Reader, output int) (*Reader, error) {
	if output != Ulaw && output != Lpcm {
		return nil, errors.New("Invalid output format")
	}
	r := Reader{
		output:    output,
		decode:    DecodeAlaw,
		transcode: Alaw2Ulaw,
		source:    reader,
		buf:       new(bytes.Buffer),
	}
	return &r, nil
}

// NewUlawDecoder returns a pointer to a Reader that decodes or trans-codes u-law data.
// It takes as input the source data Reader and the output encoding fomrat.
func NewUlawDecoder(reader io.Reader, output int) (*Reader, error) {
	if output != Alaw && output != Lpcm {
		return nil, errors.New("Invalid output format")
	}
	r := Reader{
		output:    output,
		decode:    DecodeUlaw,
		transcode: Ulaw2Alaw,
		source:    reader,
		buf:       new(bytes.Buffer),
	}
	return &r, nil
}

// NewAlawEncoder returns a pointer to a Writer that encodes data to A-law.
// It takes as input the destination data Writer and the input encoding fomrat.
func NewAlawEncoder(writer io.Writer, input int) (*Writer, error) {
	if input != Ulaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	w := Writer{
		input:       input,
		encode:      EncodeAlaw,
		transcode:   Ulaw2Alaw,
		destination: writer,
		buf:         new(bytes.Buffer),
	}
	return &w, nil
}

// NewUlawEncoder returns a pointer to a Writer that encodes data to u-law.
// It takes as input the destination data Writer and the input encoding fomrat.
func NewUlawEncoder(writer io.Writer, input int) (*Writer, error) {
	if input != Alaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	w := Writer{
		input:       input,
		encode:      EncodeUlaw,
		transcode:   Alaw2Ulaw,
		destination: writer,
		buf:         new(bytes.Buffer),
	}
	return &w, nil
}

// Reset discards the Reader state. This permits reusing a Reader rather than allocating a new one.
func (r *Reader) Reset(reader io.Reader) {
	r.buf.Reset()
	r.source = reader
}

// Reset discards the Writer state. This permits reusing a Writer rather than allocating a new one.
func (w *Writer) Reset(writer io.Writer) {
	w.buf.Reset()
	w.destination = writer
}

// Read decodes G711 data. Reads up to len(p) bytes into p, returns the number
// of bytes read and any error encountered.
func (r *Reader) Read(p []byte) (int, error) {
	var i int
	var err error
	_, err = r.buf.ReadFrom(r.source)
	if err != nil {
		return i, err
	}
	var frame byte
	if r.output == Lpcm { // Decode G711 data to LPCM
		for i = 0; i < len(p)-2; i = i + 2 {
			frame, err = r.buf.ReadByte()
			if err != nil {
				break
			}
			decoded := r.decode(frame)
			p[i] = byte(decoded)
			p[i+1] = byte(decoded >> 8)
		}
	} else { // Trans-code
		for i = 0; i < len(p); i++ {
			frame, err = r.buf.ReadByte()
			if err != nil {
				break
			}
			p[i] = r.transcode(frame)
		}
	}
	return i, err
}

// Write encodes G711 Data. Writes len(p) bytes from p to the underlying data stream,
// returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
// that caused the write to stop early.
func (w *Writer) Write(p []byte) (int, error) {
	var err, wrErr error
	if w.input == Lpcm { // Encode LPCM data to G711
		for i := 0; i <= len(p)-2; i = i + 2 {
			wrErr = w.buf.WriteByte(w.encode(int16(p[i]) | int16(p[i+1])<<8))
			if wrErr != nil {
				break
			}
		}
	} else { // Trans-code
		for _, data := range p {
			wrErr = w.buf.WriteByte(w.transcode(data))
			if wrErr != nil {
				break
			}
		}
	}
	i, err := w.buf.WriteTo(w.destination)
	if err == nil {
		err = wrErr
	}
	if w.input == Lpcm {
		i *= 2 // Report back the correct number of bytes written from p
	}
	return int(i), err
}

// Flush flushes any pending data to the underlying writer.
func (w *Writer) Flush() (err error) {
	if w.buf.Len() > 0 {
		_, err = w.buf.WriteTo(w.destination)
	}
	return
}
