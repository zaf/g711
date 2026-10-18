/*
	Copyright (C) 2016 - 2026, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.
*/

/*
Package g711 implements encoding and decoding of G711 PCM sound data.
G.711 is an ITU-T standard for audio companding.

The package exposes a high level API for encoding and decoding through
an io.WriteCloser. But also a low level API using preallocated buffers
for cases where performance and memory handling are critical.
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
	written     int8                // reported written size
	carry       byte                // first byte of a sample split between Write calls
	pending     bool                // true when a byte from a previous Write is carried
}

const (
	// written parameters
	double = iota + 1
	half
)

// NewCoder returns a pointer to a Coder that implements an io.WriteCloser.
// It takes as input the destination data Writer and the input/output encoding formats.
func NewCoder(writer io.Writer, input, output int) (*Coder, error) {
	if writer == nil {
		return nil, errors.New("io.Writer is nil")
	}
	var translate func([]byte) []byte
	var written int8
	switch input {
	case Lpcm:
		switch output {
		case Alaw:
			translate = EncodeAlaw
			written = double
		case Ulaw:
			translate = EncodeUlaw
			written = double
		default:
			return nil, errors.New("invalid output format")
		}
	case Alaw:
		switch output {
		case Lpcm:
			translate = DecodeAlaw
			written = half
		case Ulaw:
			translate = Alaw2Ulaw
		default:
			return nil, errors.New("invalid output format")
		}
	case Ulaw:
		switch output {
		case Lpcm:
			translate = DecodeUlaw
			written = half
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
		written:     written,
	}
	return &w, nil
}

// Close closes the Encoder, it implements the io.Closer interface.
// An unpaired byte remaining from a sample split between Write calls is
// discarded, as only complete samples can be encoded.
func (w *Coder) Close() error {
	w.destination = nil
	w.translate = nil
	w.written = 0
	w.carry = 0
	w.pending = false
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
	w.carry = 0
	w.pending = false
	return nil
}

// Write encodes/decodes/transcodes sound data. Writes len(p) bytes from p to the underlying data stream,
// returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
// that caused the write to stop early.
// When encoding Lpcm data, a sample split between two Write calls is buffered and
// completed with the next Write. An unpaired byte remaining at the end of the
// stream is discarded when the Coder is Closed, just like the low level API
// that only encodes complete samples.
func (w *Coder) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if w == nil || w.translate == nil || w.destination == nil {
		return 0, errors.New("coder is uninitialized or closed")
	}
	input := p
	carried := false
	if w.pending {
		// The previous Write ended in the middle of a sample, prepend the
		// carried byte so that the sample gets completed.
		input = make([]byte, 0, len(p)+1)
		input = append(input, w.carry)
		input = append(input, p...)
		w.pending = false
		carried = true
	}
	out := w.translate(input)
	i, err := w.destination.Write(out)
	if err == nil && i < len(out) {
		err = io.ErrShortWrite
	}
	if err != nil {
		// The write stopped early, report how many bytes from p made it through.
		// If we are encoding to g711 we need to multiply the number of bytes written by 2
		// this happens because 2 bytes of input data are encoded to 1 byte of output data.
		// In a similar manner if we are decoding from g711 we need to divide the number of bytes written by 2.
		switch w.written {
		case double:
			i <<= 1
			if carried {
				if i > 0 {
					// The carried byte was already reported by the previous Write.
					i--
				} else {
					// Nothing made it through, keep the carried byte pending
					// so that the sample can be completed on the next Write.
					w.pending = true
				}
			}
		case half:
			i >>= 1
		}
		return i, err
	}
	// If we are encoding to g711 and the input ended in the middle of a sample,
	// carry the last byte so that the sample gets completed by the next Write.
	if w.written == double && len(input)&1 == 1 {
		w.carry = input[len(input)-1]
		w.pending = true
	}
	return len(p), nil
}
