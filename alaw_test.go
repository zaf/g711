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
	"os"
	"testing"
)

// Benchmark EncodeAlaw
func BenchmarkEncodeAlaw(b *testing.B) {
	rawData, err := os.ReadFile("testing/speech.raw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)
	}
	b.SetBytes(int64(len(rawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeAlaw(rawData)
	}
}

// Benchmark DecodeAlaw
func BenchmarkDecodeAlaw(b *testing.B) {
	aData, err := os.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)
	}
	b.SetBytes(int64(len(aData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeAlaw(aData)
	}
}

// Benchmark Alaw2Ulaw
func BenchmarkAlaw2Ulaw(b *testing.B) {
	aData, err := os.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)
	}
	b.SetBytes(int64(len(aData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Alaw2Ulaw(aData)
	}
}
