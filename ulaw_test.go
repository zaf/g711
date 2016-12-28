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
	"io/ioutil"
	"log"
	"testing"
)

// Benchmark EncodeUlaw
func BenchmarkEncodeUlaw(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i <= len(rawData)-2; i = i + 2 {
			EncodeUlaw(int16(rawData[i]) | int16(rawData[i+1])<<8)
		}
	}
}

// Benchmark DecodeUlaw
func BenchmarkDecodeUlaw(b *testing.B) {
	aData, err := ioutil.ReadFile("testing/speech.ulaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, frame := range aData {
			DecodeUlaw(frame)
		}
	}
}

// Benchmark Ulaw2Alaw
func BenchmarkUlaw2Alaw(b *testing.B) {
	aData, err := ioutil.ReadFile("testing/speech.ulaw")
	if err != nil {
		log.Printf("Failed to read test data: %s\n", err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, frame := range aData {
			Ulaw2Alaw(frame)
		}
	}
}
