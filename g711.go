/*
	Copyright (C) 2016 - 2024, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.
*/

/*
Package g711 implements encoding and decoding of G711 PCM sound data.
G.711 is an ITU-T standard for audio companding.

For usage details please see the code snippet in the cmd folder.
*/
package g711

import (
	"errors"
	"io"
)

const (
	// Input and output formats
	Alaw = iota + 1 // Alaw G711 encoded PCM data
	Ulaw            // Ulaw G711  encoded PCM data
	Lpcm            // Lpcm 16bit signed linear data
)

// Coder encodes 16bit 8000Hz LPCM data to G711 PCM, or
// decodes G711 PCM data to 16bit 8000Hz LPCM data, or
// directly transcodes between A-law and u-law
type Coder struct {
	translate   func([]byte) []byte // enc/decoding function
	destination io.Writer           // output data
	multiplier  float64
}

// NewCoder returns a pointer to a Coder that implements an io.WriteCloser.
// It takes as input the destination data Writer and the input/output encoding formats.
func NewCoder(writer io.Writer, input, output int) (*Coder, error) {
	if writer == nil {
		return nil, errors.New("io.Writer is nil")
	}
	var translate func([]byte) []byte
	multiplier := 1.0
	switch input {
	case Lpcm:
		switch output {
		case Alaw:
			translate = EncodeAlaw
			multiplier = 2
		case Ulaw:
			translate = EncodeUlaw
			multiplier = 2
		default:
			return nil, errors.New("invalid output format")
		}
	case Alaw:
		switch output {
		case Lpcm:
			translate = DecodeAlaw
			multiplier = 0.5
		case Ulaw:
			translate = Alaw2Ulaw
		default:
			return nil, errors.New("invalid output format")
		}
	case Ulaw:
		switch output {
		case Lpcm:
			translate = DecodeUlaw
			multiplier = 0.5
		case Alaw:
			translate = Ulaw2Alaw
		default:
			return nil, errors.New("invalid output format")
		}
	default:
		return nil, errors.New("invalid input format")
	}
	w := Coder{
		translate:   translate,
		destination: writer,
		multiplier:  multiplier,
	}
	return &w, nil
}

// Close closes the Encoder, it implements the io.Closer interface.
func (w *Coder) Close() error {
	w.destination = nil
	w.translate = nil
	w.multiplier = 0
	w = nil
	return nil
}

// Reset discards the Encoder state. This permits reusing an Encoder rather than allocating a new one.
func (w *Coder) Reset(writer io.Writer) error {
	if w == nil || w.translate == nil || w.destination == nil {
		return errors.New("coder is uninitialized or closed")
	}
	if writer == nil {
		return errors.New("io.Writer is nil")
	}
	w.destination = writer
	return nil
}

// Write encodes/decodes/transcodes sound data. Writes len(p) bytes from p to the underlying data stream,
// returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
// that caused the write to stop early.
func (w *Coder) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	i, err := w.destination.Write(w.translate(p))
	// If we are encoding to g711 we need to multiply the number of bytes written by 2 to avoid reporting short writes
	// this happens because 2 bytes of input data are encoded to 1 byte of output data.
	// In a similar manner if we are decoding from g711 we need to divide the number of bytes written by 2.
	i = int(float64(i) * w.multiplier)
	return i, err
}
