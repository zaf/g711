/*
	Copyright (C) 2016 - 2017, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	Package g711 implements encoding and decoding of G711.0 compressed sound data.
	G.711 is an ITU-T standard for audio companding.
*/

package g711

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"
)

var EncoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 12},
}

var DecoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 24},
	{[]byte{0x01, 0x00, 0xdc, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 26},
}

var TranscoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 13},
}

// Test Encoding
func TestEncode(t *testing.T) {
	aenc, _ := NewAlawEncoder(ioutil.Discard, Lpcm)
	for _, tc := range EncoderTest {
		i, _ := aenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("Alaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	uenc, _ := NewUlawEncoder(ioutil.Discard, Lpcm)
	for _, tc := range EncoderTest {
		i, _ := uenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	utrans, _ := NewUlawEncoder(ioutil.Discard, Alaw)
	for _, tc := range TranscoderTest {
		i, _ := utrans.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Transcode via encoder: expected: %d , actual: %d", tc.expected, i)
		}
	}
}

// Test Decoding
func TestDecode(t *testing.T) {
	b := new(bytes.Buffer)
	adec, _ := NewAlawDecoder(b, Lpcm)
	for _, tc := range DecoderTest {
		b.Write(tc.data)
		p := make([]byte, 2*tc.expected)
		i, _ := adec.Read(p)
		if i != tc.expected {
			t.Errorf("Alaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	b.Reset()
	udec, _ := NewUlawDecoder(b, Lpcm)
	for _, tc := range DecoderTest {
		b.Write(tc.data)
		p := make([]byte, 2*tc.expected)
		i, _ := udec.Read(p)
		if i != tc.expected {
			t.Errorf("ulaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	b.Reset()
	utrans, _ := NewUlawDecoder(b, Alaw)
	for _, tc := range TranscoderTest {
		b.Write(tc.data)
		p := make([]byte, 2*tc.expected)
		i, _ := utrans.Read(p)
		if i != tc.expected {
			t.Errorf("ulaw Transcode via decoder: expected: %d , actual: %d", tc.expected, i)
		}
	}
}

// Benchmark Encoding data to Alaw
func BenchmarkAEncode(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewAlawEncoder(ioutil.Discard, Lpcm)
		if err != nil {
			log.Printf("Failed to create Writer: %s\n", err)
			b.FailNow()
		}
		_, err = encoder.Write(rawData)
		if err != nil {
			log.Printf("Encoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Encoding data to Ulaw
func BenchmarkUEncode(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewUlawEncoder(ioutil.Discard, Lpcm)
		if err != nil {
			log.Printf("Failed to create Writer: %s\n", err)
			b.FailNow()
		}
		_, err = encoder.Write(rawData)
		if err != nil {
			log.Printf("Encoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark transcoding g711 data via Reader
func BenchmarkTranscodeR(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transcoder, err := NewAlawDecoder(bytes.NewReader(alawData), Ulaw)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(transcoder)
		if err != nil {
			log.Printf("Transcoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark transcoding g711 data via Writer
func BenchmarkTranscodeW(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transcoder, err := NewAlawEncoder(ioutil.Discard, Ulaw)
		if err != nil {
			log.Printf("Failed to create Writer: %s\n", err)
			b.FailNow()
		}
		_, err = transcoder.Write(alawData)
		if err != nil {
			log.Printf("Transcoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Decoding Ulaw data
func BenchmarkUDecode(b *testing.B) {
	ulawData, err := ioutil.ReadFile("testing/speech.ulaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewUlawDecoder(bytes.NewReader(ulawData), Lpcm)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(decoder)
		if err != nil {
			log.Printf("Decoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Decoding Alaw data
func BenchmarkADecode(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewAlawDecoder(bytes.NewReader(alawData), Lpcm)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(decoder)
		if err != nil {
			log.Printf("Decoding failed: %s\n", err)
			b.FailNow()
		}
	}
}
