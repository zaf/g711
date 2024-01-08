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
	"io"
	"os"
	"testing"
)

var EncoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01}, 0},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 12},
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
		i, _ := aenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("Alaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	aenc.Close()
	uenc, _ := NewCoder(io.Discard, Lpcm, Ulaw)
	for _, tc := range EncoderTest {
		i, _ := uenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	uenc.Close()
	utrans, _ := NewCoder(io.Discard, Alaw, Ulaw)
	for _, tc := range TranscoderTest {
		i, _ := utrans.Write(tc.data)
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
		i, _ := adec.Write(tc.data)
		if i != tc.expected {
			t.Errorf("Alaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	adec.Close()
	udec, _ := NewCoder(io.Discard, Ulaw, Lpcm)
	for _, tc := range DecoderTest {
		i, _ := udec.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	udec.Close()
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
