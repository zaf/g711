/*
	Copyright (C) 2016 - 2026, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	Package g711 implements encoding and decoding of G711 PCM sound data.
	G.711 is an ITU-T standard for audio companding.
*/

package g711

import (
	"encoding/binary"
	"math/bits"
)

var (
	// A-law to LPCM conversion lookup table
	alaw2lpcm = [256]int16{
		-5504, -5248, -6016, -5760, -4480, -4224, -4992, -4736,
		-7552, -7296, -8064, -7808, -6528, -6272, -7040, -6784,
		-2752, -2624, -3008, -2880, -2240, -2112, -2496, -2368,
		-3776, -3648, -4032, -3904, -3264, -3136, -3520, -3392,
		-22016, -20992, -24064, -23040, -17920, -16896, -19968, -18944,
		-30208, -29184, -32256, -31232, -26112, -25088, -28160, -27136,
		-11008, -10496, -12032, -11520, -8960, -8448, -9984, -9472,
		-15104, -14592, -16128, -15616, -13056, -12544, -14080, -13568,
		-344, -328, -376, -360, -280, -264, -312, -296,
		-472, -456, -504, -488, -408, -392, -440, -424,
		-88, -72, -120, -104, -24, -8, -56, -40,
		-216, -200, -248, -232, -152, -136, -184, -168,
		-1376, -1312, -1504, -1440, -1120, -1056, -1248, -1184,
		-1888, -1824, -2016, -1952, -1632, -1568, -1760, -1696,
		-688, -656, -752, -720, -560, -528, -624, -592,
		-944, -912, -1008, -976, -816, -784, -880, -848,
		5504, 5248, 6016, 5760, 4480, 4224, 4992, 4736,
		7552, 7296, 8064, 7808, 6528, 6272, 7040, 6784,
		2752, 2624, 3008, 2880, 2240, 2112, 2496, 2368,
		3776, 3648, 4032, 3904, 3264, 3136, 3520, 3392,
		22016, 20992, 24064, 23040, 17920, 16896, 19968, 18944,
		30208, 29184, 32256, 31232, 26112, 25088, 28160, 27136,
		11008, 10496, 12032, 11520, 8960, 8448, 9984, 9472,
		15104, 14592, 16128, 15616, 13056, 12544, 14080, 13568,
		344, 328, 376, 360, 280, 264, 312, 296,
		472, 456, 504, 488, 408, 392, 440, 424,
		88, 72, 120, 104, 24, 8, 56, 40,
		216, 200, 248, 232, 152, 136, 184, 168,
		1376, 1312, 1504, 1440, 1120, 1056, 1248, 1184,
		1888, 1824, 2016, 1952, 1632, 1568, 1760, 1696,
		688, 656, 752, 720, 560, 528, 624, 592,
		944, 912, 1008, 976, 816, 784, 880, 848,
	}
	// A-law to u-law conversion lookup table based on the ITU-T G.711 specification
	alaw2ulaw = [256]uint8{
		41, 42, 39, 40, 45, 46, 43, 44, 33, 34, 31, 32, 37, 38, 35, 36,
		57, 58, 55, 56, 61, 62, 59, 60, 49, 50, 47, 48, 53, 54, 51, 52,
		10, 11, 8, 9, 14, 15, 12, 13, 2, 3, 0, 1, 6, 7, 4, 5,
		26, 27, 24, 25, 30, 31, 28, 29, 18, 19, 16, 17, 22, 23, 20, 21,
		98, 99, 96, 97, 102, 103, 100, 101, 93, 93, 92, 92, 95, 95, 94, 94,
		116, 118, 112, 114, 124, 126, 120, 122, 106, 107, 104, 105, 110, 111, 108, 109,
		72, 73, 70, 71, 76, 77, 74, 75, 64, 65, 63, 63, 68, 69, 66, 67,
		86, 87, 84, 85, 90, 91, 88, 89, 79, 79, 78, 78, 82, 83, 80, 81,
		169, 170, 167, 168, 173, 174, 171, 172, 161, 162, 159, 160, 165, 166, 163, 164,
		185, 186, 183, 184, 189, 190, 187, 188, 177, 178, 175, 176, 181, 182, 179, 180,
		138, 139, 136, 137, 142, 143, 140, 141, 130, 131, 128, 129, 134, 135, 132, 133,
		154, 155, 152, 153, 158, 159, 156, 157, 146, 147, 144, 145, 150, 151, 148, 149,
		226, 227, 224, 225, 230, 231, 228, 229, 221, 221, 220, 220, 223, 223, 222, 222,
		244, 246, 240, 242, 252, 254, 248, 250, 234, 235, 232, 233, 238, 239, 236, 237,
		200, 201, 198, 199, 204, 205, 202, 203, 192, 193, 191, 191, 196, 197, 194, 195,
		214, 215, 212, 213, 218, 219, 216, 217, 207, 207, 206, 206, 210, 211, 208, 209,
	}
)

// LPCM to A-law conversion lookup table, indexed by the 12 Most Significant
// Bits of the input sample, as A-law uses only the 12 MSBs of the 16bit input
var lpcm2alaw [1 << 12]uint8

func init() {
	for i := range lpcm2alaw {
		lpcm2alaw[i] = EncodeAlawFrame(int16(i << 4))
	}
}

// EncodeAlaw encodes 16bit LPCM data to G711 A-law PCM
func EncodeAlaw(lpcm []byte) []byte {
	alaw := make([]byte, len(lpcm)>>1)
	EncodeAlawTo(lpcm, alaw)
	return alaw
}

// EncodeAlawTo encodes 16bit LPCM data to G711 A-law PCM
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (half the size of the LPCM data).
func EncodeAlawTo(lpcm, alaw []byte) {
	n := len(lpcm)
	i := 0
	// Process 4 frames per iteration, so that 4 independent table lookups
	// are in flight at any time, keeping the CPU pipeline busy. The 8 input
	// bytes are read with a single wide load. Writing the loop guard as
	// i < n-7 (instead of i+8 <= n) lets the compiler prove that all loads
	// are in bounds, so it elides their bounds checks.
	for ; i < n-7; i += 8 {
		w := binary.LittleEndian.Uint64(lpcm[i:])
		a := lpcm2alaw[uint16(w)>>4]
		b := lpcm2alaw[uint16(w>>16)>>4]
		c := lpcm2alaw[uint16(w>>32)>>4]
		d := lpcm2alaw[uint16(w>>48)>>4]
		alaw[i>>1] = a
		alaw[(i>>1)+1] = b
		alaw[(i>>1)+2] = c
		alaw[(i>>1)+3] = d
	}
	// Handle trailing frames
	for ; i < n-1; i += 2 {
		alaw[i>>1] = lpcm2alaw[(uint16(lpcm[i])|uint16(lpcm[i+1])<<8)>>4]
	}
}

// EncodeAlawFrame encodes a 16bit LPCM frame to G711 A-law PCM
func EncodeAlawFrame(frame int16) uint8 {
	sign := ((^frame) >> 8) & 0x80
	if sign == 0 {
		frame = ^frame
	}
	compressedByte := frame >> 4
	if compressedByte > 15 {
		seg := int16(12 - bits.LeadingZeros16(uint16(compressedByte)))
		compressedByte >>= seg - 1
		compressedByte -= 16
		compressedByte += seg << 4
	}
	return uint8((sign | compressedByte) ^ 0x0055)
}

// DecodeAlaw decodes A-law PCM data to 16bit LPCM
func DecodeAlaw(pcm []byte) []byte {
	lpcm := make([]byte, len(pcm)*2)
	DecodeAlawTo(pcm, lpcm)
	return lpcm
}

// DecodeAlawTo decodes A-law PCM data to 16bit LPCM
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (double the size of the PCM data).
func DecodeAlawTo(pcm, lpcm []byte) {
	n := len(pcm)
	i, j := 0, 0
	// Process 8 samples per iteration, so that 8 independent table lookups
	// are in flight at any time, keeping the CPU pipeline busy. The decoded
	// 16bit samples are packed into uint64 values so that each pair of them
	// is stored with a single wide store.
	// NB: the samples are converted to uint16 first, converting an int16
	// directly to uint64 would sign-extend and corrupt the packed value.
	// As in EncodeAlawTo, the i < n-7 guard lets the compiler elide the
	// input-side bounds checks.
	for ; i < n-7; i, j = i+8, j+16 {
		w := binary.LittleEndian.Uint64(pcm[i:])
		lo := uint64(uint16(alaw2lpcm[uint8(w)])) |
			uint64(uint16(alaw2lpcm[uint8(w>>8)]))<<16 |
			uint64(uint16(alaw2lpcm[uint8(w>>16)]))<<32 |
			uint64(uint16(alaw2lpcm[uint8(w>>24)]))<<48
		hi := uint64(uint16(alaw2lpcm[uint8(w>>32)])) |
			uint64(uint16(alaw2lpcm[uint8(w>>40)]))<<16 |
			uint64(uint16(alaw2lpcm[uint8(w>>48)]))<<32 |
			uint64(uint16(alaw2lpcm[uint8(w>>56)]))<<48
		binary.LittleEndian.PutUint64(lpcm[j:], lo)
		binary.LittleEndian.PutUint64(lpcm[j+8:], hi)
	}
	// Handle trailing samples
	for ; i < n; i, j = i+1, j+2 {
		binary.LittleEndian.PutUint16(lpcm[j:], uint16(alaw2lpcm[pcm[i]]))
	}
}

// DecodeAlawFrame decodes an A-law PCM frame to 16bit LPCM
func DecodeAlawFrame(frame uint8) int16 {
	return alaw2lpcm[frame]
}

// Alaw2Ulaw performs direct A-law to u-law data conversion
func Alaw2Ulaw(alaw []byte) []byte {
	ulaw := make([]byte, len(alaw))
	Alaw2UlawTo(alaw, ulaw)
	return ulaw
}

// Alaw2UlawTo performs direct A-law to u-law data conversion
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (the size of the A-law data).
func Alaw2UlawTo(alaw, ulaw []byte) {
	n := len(alaw)
	i := 0
	// Process 8 samples per iteration, so that 8 independent table lookups
	// are in flight at any time, keeping the CPU pipeline busy. The results
	// are packed into a uint64 value so that all 8 samples are stored with a
	// single wide store. The i < n-7 loop guard lets the compiler elide the
	// input-side bounds checks.
	for ; i < n-7; i += 8 {
		w := binary.LittleEndian.Uint64(alaw[i:])
		o := uint64(alaw2ulaw[uint8(w)]) |
			uint64(alaw2ulaw[uint8(w>>8)])<<8 |
			uint64(alaw2ulaw[uint8(w>>16)])<<16 |
			uint64(alaw2ulaw[uint8(w>>24)])<<24 |
			uint64(alaw2ulaw[uint8(w>>32)])<<32 |
			uint64(alaw2ulaw[uint8(w>>40)])<<40 |
			uint64(alaw2ulaw[uint8(w>>48)])<<48 |
			uint64(alaw2ulaw[uint8(w>>56)])<<56
		binary.LittleEndian.PutUint64(ulaw[i:], o)
	}
	// Handle trailing samples
	for ; i < n; i++ {
		ulaw[i] = alaw2ulaw[alaw[i]]
	}
}

// Alaw2UlawFrame directly converts an A-law frame to u-law
func Alaw2UlawFrame(frame uint8) uint8 {
	return alaw2ulaw[frame]
}
