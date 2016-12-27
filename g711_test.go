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

// Benchmark Encoding data to Alaw
func BenchmarkAWrite(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alawWr, err := NewAlawWriter(ioutil.Discard, Lpcm)
		if err != nil {
			log.Printf("Failed to create Writer: %s\n", err)
			b.FailNow()
		}
		_, err = alawWr.Write(rawData)
		if err != nil {
			log.Printf("Encoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Encoding data to Ulaw
func BenchmarkUWrite(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ulawWr, err := NewUlawWriter(ioutil.Discard, Lpcm)
		if err != nil {
			log.Printf("Failed to create Writer: %s\n", err)
			b.FailNow()
		}
		_, err = ulawWr.Write(rawData)
		if err != nil {
			log.Printf("Encoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark transcoding g711 data
func BenchmarkTranscode(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alawR, err := NewAlawReader(bytes.NewReader(alawData), Ulaw)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(alawR)
		if err != nil {
			log.Printf("Transcoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Decoding Ulaw data
func BenchmarkURead(b *testing.B) {
	ulawData, err := ioutil.ReadFile("testing/speech.ulaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ulawR, err := NewUlawReader(bytes.NewReader(ulawData), Lpcm)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(ulawR)
		if err != nil {
			log.Printf("Decoding failed: %s\n", err)
			b.FailNow()
		}
	}
}

// Benchmark Decoding Alaw data
func BenchmarkARead(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alawR, err := NewAlawReader(bytes.NewReader(alawData), Lpcm)
		if err != nil {
			log.Printf("Failed to create Reader: %s\n", err)
			b.FailNow()
		}
		_, err = ioutil.ReadAll(alawR)
		if err != nil {
			log.Printf("Decoding failed: %s\n", err)
			b.FailNow()
		}
	}
}
