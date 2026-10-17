/*
	Copyright (C) 2016 - 2024, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	Package g711 implements encoding and decoding of G711 PCM sound data.
	G.711 is an ITU-T standard for audio companding.
*/

package g711

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestNewCoder(t *testing.T) {
	writer := bytes.NewBuffer([]byte{})
	input := Lpcm
	output := Alaw
	// Test case: Valid input
	coder, err := NewCoder(writer, input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if coder == nil {
		t.Error("Expected coder to not be nil")
	}
	// Test case: Nil writer
	_, err = NewCoder(nil, input, output)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Test case: Invalid input format
	_, err = NewCoder(writer, 999, output)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Test case: Invalid output format
	_, err = NewCoder(writer, input, 999)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReset(t *testing.T) {
	// Test case 1: Coder is not nil, writer is not nil
	w, _ := NewCoder(io.Discard, Lpcm, Alaw)
	writer := bytes.NewBufferString("")
	err := w.Reset(writer)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	w.Close()
	// Test case 2: Coder is nil
	nilCoder, _ := NewCoder(io.Discard, Lpcm, Alaw)
	nilCoder.Close()
	err = nilCoder.Reset(writer)
	if err == nil || err.Error() != "coder is uninitialized or closed" {
		t.Errorf("Expected error: 'coder is uninitialized or closed', but got: %v", err)
	}
	// Test case 3: writer is nil
	w, _ = NewCoder(io.Discard, Lpcm, Alaw)
	err = w.Reset(nil)
	if err == nil || err.Error() != "io.Writer is nil" {
		t.Errorf("Expected error: 'io.Writer is nil', but got: %v", err)
	}
}

var EncoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 13},
}

var DecoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0xdc, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 13},
}

var TranscoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 13},
}

// Test Encoding
func TestEncode(t *testing.T) {
	aenc, _ := NewCoder(io.Discard, Lpcm, Alaw)
	for _, tc := range EncoderTest {
		i, err := aenc.Write(tc.data)
		if err != nil {
			t.Errorf("Alaw Encode: unexpected error: %v", err)
		}
		if i != tc.expected {
			t.Errorf("Alaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	aenc.Close()
	uenc, _ := NewCoder(io.Discard, Lpcm, Ulaw)
	for _, tc := range EncoderTest {
		i, err := uenc.Write(tc.data)
		if err != nil {
			t.Errorf("ulaw Encode: unexpected error: %v", err)
		}
		if i != tc.expected {
			t.Errorf("ulaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	uenc.Close()
	utrans, _ := NewCoder(io.Discard, Alaw, Ulaw)
	for _, tc := range TranscoderTest {
		i, err := utrans.Write(tc.data)
		if err != nil {
			t.Errorf("ulaw Transcode: unexpected error: %v", err)
		}
		if i != tc.expected {
			t.Errorf("ulaw Transcode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	utrans.Close()
}

// Test Decoding
func TestDecode(t *testing.T) {
	adec, _ := NewCoder(io.Discard, Alaw, Lpcm)
	for _, tc := range DecoderTest {
		i, err := adec.Write(tc.data)
		if err != nil {
			t.Errorf("Alaw Decode: unexpected error: %v", err)
		}
		if i != tc.expected {
			t.Errorf("Alaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	adec.Close()
	udec, _ := NewCoder(io.Discard, Ulaw, Lpcm)
	for _, tc := range DecoderTest {
		i, err := udec.Write(tc.data)
		if err != nil {
			t.Errorf("ulaw Decode: unexpected error: %v", err)
		}
		if i != tc.expected {
			t.Errorf("ulaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	udec.Close()
}

// Test writing after Close and on a zero value Coder
func TestWriteAfterClose(t *testing.T) {
	writer := bytes.NewBuffer([]byte{})
	coder, _ := NewCoder(writer, Lpcm, Alaw)
	coder.Close()
	_, err := coder.Write([]byte{0x00, 0x00})
	if err == nil {
		t.Error("Expected error when writing after Close, got nil")
	}
	var zeroCoder Coder
	_, err = zeroCoder.Write([]byte{0x00, 0x00})
	if err == nil {
		t.Error("Expected error when writing to a zero value Coder, got nil")
	}
}

// Test that samples split between Write calls are completed and that
// the streamed output matches the output of the low level buffer API
func TestSplitSample(t *testing.T) {
	rawData, err := os.ReadFile("testing/speech.raw")
	if err != nil {
		t.Fatalf("Failed to read test data: %s\n", err)
	}
	for _, chunk := range []int{1, 3, 5, 7, 1023, 4095} {
		for _, length := range []int{len(rawData), len(rawData) - 1} {
			for _, format := range []int{Alaw, Ulaw} {
				data := rawData[:length]
				var buf bytes.Buffer
				coder, _ := NewCoder(&buf, Lpcm, format)
				for off := 0; off < len(data); off += chunk {
					end := off + chunk
					if end > len(data) {
						end = len(data)
					}
					i, err := coder.Write(data[off:end])
					if err != nil {
						t.Fatalf("Write failed: %s\n", err)
					}
					if i != end-off {
						t.Errorf("Chunk size %d: expected %d bytes written, actual: %d\n", chunk, end-off, i)
					}
				}
				coder.Close()
				var ref []byte
				if format == Alaw {
					ref = EncodeAlaw(data)
				} else {
					ref = EncodeUlaw(data)
				}
				if !bytes.Equal(buf.Bytes(), ref) {
					t.Errorf("Chunk size %d: streamed output does not match buffer API output\n", chunk)
				}
			}
		}
	}
}

// Test that Writes that stop early can be resumed without data loss
func TestResumeAfterShortWrite(t *testing.T) {
	rawData, err := os.ReadFile("testing/speech.raw")
	if err != nil {
		t.Fatalf("Failed to read test data: %s\n", err)
	}
	for _, k := range []int{1, 999} {
		writer := &limitedWriter{k: k}
		coder, _ := NewCoder(writer, Lpcm, Alaw)
		data := rawData[:4002]
		for len(data) > 0 {
			i, err := coder.Write(data)
			if err != nil {
				if i < 0 || i > len(data) {
					t.Fatalf("Invalid number of bytes reported: %d\n", i)
				}
				data = data[i:]
				continue
			}
			break
		}
		coder.Close()
		if !bytes.Equal(writer.buf.Bytes(), EncodeAlaw(rawData[:4002])) {
			t.Errorf("Limit %d: resumed output does not match buffer API output\n", k)
		}
	}
}

// limitedWriter accepts up to k bytes per Write, failing anything longer
type limitedWriter struct {
	buf bytes.Buffer
	k   int
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	if len(p) <= l.k {
		return l.buf.Write(p)
	}
	l.buf.Write(p[:l.k])
	return l.k, errors.New("short write")
}

// Benchmark Encoding data to Alaw
func BenchmarkAEncode(b *testing.B) {
	rawData, err := os.ReadFile("testing/speech.raw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)
	}
	b.SetBytes(int64(len(rawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewCoder(io.Discard, Lpcm, Alaw)
		if err != nil {
			b.Fatalf("Failed to create Encoder: %s\n", err)
		}
		_, err = encoder.Write(rawData)
		if err != nil {
			b.Fatalf("Encoding failed: %s\n", err)
		}
		encoder.Close()
	}
}

// Benchmark Encoding data to Ulaw
func BenchmarkUEncode(b *testing.B) {
	rawData, err := os.ReadFile("testing/speech.raw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(rawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewCoder(io.Discard, Lpcm, Ulaw)
		if err != nil {
			b.Fatalf("Failed to create Encoder: %s\n", err)

		}
		_, err = encoder.Write(rawData)
		if err != nil {
			b.Fatalf("Encoding failed: %s\n", err)

		}
		encoder.Close()
	}
}

// Benchmark transcoding g711 data
func BenchmarkTranscode(b *testing.B) {
	alawData, err := os.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(alawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transcoder, err := NewCoder(io.Discard, Ulaw, Alaw)
		if err != nil {
			b.Fatalf("Failed to create Transcoder: %s\n", err)

		}
		_, err = transcoder.Write(alawData)
		if err != nil {
			b.Fatalf("Transcoding failed: %s\n", err)

		}
		transcoder.Close()
	}
}

// Benchmark Decoding Ulaw data
func BenchmarkUDecode(b *testing.B) {
	ulawData, err := os.ReadFile("testing/speech.ulaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(ulawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewCoder(io.Discard, Ulaw, Lpcm)
		if err != nil {
			b.Fatalf("Failed to create Reader: %s\n", err)

		}
		_, err = decoder.Write(ulawData)
		if err != nil {
			b.Fatalf("Decoding failed: %s\n", err)

		}
		decoder.Close()
	}
}

// Benchmark Decoding Alaw data
func BenchmarkADecode(b *testing.B) {
	alawData, err := os.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(alawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewCoder(io.Discard, Alaw, Lpcm)
		if err != nil {
			b.Fatalf("Failed to create Reader: %s\n", err)

		}
		_, err = decoder.Write(alawData)
		if err != nil {
			b.Fatalf("Decoding failed: %s\n", err)

		}
		decoder.Close()
	}
}
