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

const (
	ulawBias = 33
	ulawClip = 0x1FFF
)

var (
	// u-law to LPCM conversion lookup table
	ulaw2lpcm = [256]int16{
		-32124, -31100, -30076, -29052, -28028, -27004, -25980, -24956,
		-23932, -22908, -21884, -20860, -19836, -18812, -17788, -16764,
		-15996, -15484, -14972, -14460, -13948, -13436, -12924, -12412,
		-11900, -11388, -10876, -10364, -9852, -9340, -8828, -8316,
		-7932, -7676, -7420, -7164, -6908, -6652, -6396, -6140,
		-5884, -5628, -5372, -5116, -4860, -4604, -4348, -4092,
		-3900, -3772, -3644, -3516, -3388, -3260, -3132, -3004,
		-2876, -2748, -2620, -2492, -2364, -2236, -2108, -1980,
		-1884, -1820, -1756, -1692, -1628, -1564, -1500, -1436,
		-1372, -1308, -1244, -1180, -1116, -1052, -988, -924,
		-876, -844, -812, -780, -748, -716, -684, -652,
		-620, -588, -556, -524, -492, -460, -428, -396,
		-372, -356, -340, -324, -308, -292, -276, -260,
		-244, -228, -212, -196, -180, -164, -148, -132,
		-120, -112, -104, -96, -88, -80, -72, -64,
		-56, -48, -40, -32, -24, -16, -8, 0,
		32124, 31100, 30076, 29052, 28028, 27004, 25980, 24956,
		23932, 22908, 21884, 20860, 19836, 18812, 17788, 16764,
		15996, 15484, 14972, 14460, 13948, 13436, 12924, 12412,
		11900, 11388, 10876, 10364, 9852, 9340, 8828, 8316,
		7932, 7676, 7420, 7164, 6908, 6652, 6396, 6140,
		5884, 5628, 5372, 5116, 4860, 4604, 4348, 4092,
		3900, 3772, 3644, 3516, 3388, 3260, 3132, 3004,
		2876, 2748, 2620, 2492, 2364, 2236, 2108, 1980,
		1884, 1820, 1756, 1692, 1628, 1564, 1500, 1436,
		1372, 1308, 1244, 1180, 1116, 1052, 988, 924,
		876, 844, 812, 780, 748, 716, 684, 652,
		620, 588, 556, 524, 492, 460, 428, 396,
		372, 356, 340, 324, 308, 292, 276, 260,
		244, 228, 212, 196, 180, 164, 148, 132,
		120, 112, 104, 96, 88, 80, 72, 64,
		56, 48, 40, 32, 24, 16, 8, 0,
	}
	// u-law to A-law conversion lookup table based on the ITU-T G.711 specification
	ulaw2alaw = [256]uint8{
		42, 43, 40, 41, 46, 47, 44, 45, 34, 35, 32, 33, 38, 39, 36, 37,
		58, 59, 56, 57, 62, 63, 60, 61, 50, 51, 48, 49, 54, 55, 52, 53,
		11, 8, 9, 14, 15, 12, 13, 2, 3, 0, 1, 6, 7, 4, 5, 26,
		27, 24, 25, 30, 31, 28, 29, 18, 19, 16, 17, 22, 23, 20, 21, 107,
		104, 105, 110, 111, 108, 109, 98, 99, 96, 97, 102, 103, 100, 101, 123, 121,
		126, 127, 124, 125, 114, 115, 112, 113, 118, 119, 116, 117, 75, 73, 79, 77,
		66, 67, 64, 65, 70, 71, 68, 69, 90, 91, 88, 89, 94, 95, 92, 93,
		82, 83, 83, 80, 80, 81, 81, 86, 86, 87, 87, 84, 84, 85, 85, 213,
		170, 171, 168, 169, 174, 175, 172, 173, 162, 163, 160, 161, 166, 167, 164, 165,
		186, 187, 184, 185, 190, 191, 188, 189, 178, 179, 176, 177, 182, 183, 180, 181,
		139, 136, 137, 142, 143, 140, 141, 130, 131, 128, 129, 134, 135, 132, 133, 154,
		155, 152, 153, 158, 159, 156, 157, 146, 147, 144, 145, 150, 151, 148, 149, 235,
		232, 233, 238, 239, 236, 237, 226, 227, 224, 225, 230, 231, 228, 229, 251, 249,
		254, 255, 252, 253, 242, 243, 240, 241, 246, 247, 244, 245, 203, 201, 207, 205,
		194, 195, 192, 193, 198, 199, 196, 197, 218, 219, 216, 217, 222, 223, 220, 221,
		210, 210, 211, 211, 208, 208, 209, 209, 214, 214, 215, 215, 212, 212, 213, 213,
	}
)

// LPCM to u-law conversion lookup table, indexed by the 14 Most Significant
// Bits of the input sample, as u-law uses only the 14 MSBs of the 16bit input
var lpcm2ulaw [1 << 14]uint8

func init() {
	for i := range lpcm2ulaw {
		lpcm2ulaw[i] = EncodeUlawFrame(int16(i << 2))
	}
}

// EncodeUlaw encodes 16bit LPCM data to G711 u-law PCM
func EncodeUlaw(lpcm []byte) []byte {
	ulaw := make([]byte, len(lpcm)>>1)
	EncodeUlawTo(lpcm, ulaw)
	return ulaw
}

// EncodeUlawTo encodes 16bit LPCM data to G711 u-law PCM
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (half the size of the LPCM data).
func EncodeUlawTo(lpcm, ulaw []byte) {
	n := len(lpcm)
	i := 0
	// Process 8 frames per iteration, so that 8 independent table lookups
	// are in flight at any time, keeping the CPU pipeline busy. The 16 input
	// bytes are read with two wide loads. Writing the loop  as
	// i < n-15 (instead of i+16 <= n) lets the compiler prove that the
	// input loads are in bounds, so it avoids bounds checks. The
	// output frames are stored through a fixed-size subslice: slicing costs
	// a single bounds check per iteration and the following constant-index
	// stores are then provably in bounds, so they carry no bounds checks.
	for ; i < n-15; i += 16 {
		w0 := binary.LittleEndian.Uint64(lpcm[i:])
		w1 := binary.LittleEndian.Uint64(lpcm[i+8:])
		o := ulaw[i>>1 : (i>>1)+8]
		o[0] = lpcm2ulaw[uint16(w0)>>2]
		o[1] = lpcm2ulaw[uint16(w0>>16)>>2]
		o[2] = lpcm2ulaw[uint16(w0>>32)>>2]
		o[3] = lpcm2ulaw[uint16(w0>>48)>>2]
		o[4] = lpcm2ulaw[uint16(w1)>>2]
		o[5] = lpcm2ulaw[uint16(w1>>16)>>2]
		o[6] = lpcm2ulaw[uint16(w1>>32)>>2]
		o[7] = lpcm2ulaw[uint16(w1>>48)>>2]
	}
	// Handle trailing frames
	for ; i < n-1; i += 2 {
		ulaw[i>>1] = lpcm2ulaw[(uint16(lpcm[i])|uint16(lpcm[i+1])<<8)>>2]
	}
}

// EncodeUlawFrame encodes a 16bit LPCM frame to G711 u-law PCM
func EncodeUlawFrame(frame int16) uint8 {
	sign := ((^frame) >> 8) & 0x80
	if sign == 0 {
		frame = ^frame
	}
	frame = (frame >> 2) + ulawBias
	if frame > ulawClip {
		frame = ulawClip
	}
	seg := int16(16 - bits.LeadingZeros16(uint16(frame>>5)))
	lowNibble := 0x000F - ((frame >> (seg)) & 0x000F)
	return uint8(sign | ((8 - seg) << 4) | lowNibble)
}

// DecodeUlaw decodes u-law PCM data to 16bit LPCM
func DecodeUlaw(pcm []byte) []byte {
	lpcm := make([]byte, len(pcm)*2)
	DecodeUlawTo(pcm, lpcm)
	return lpcm
}

// DecodeUlawTo decodes u-law PCM data to 16bit LPCM
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (double the size of the PCM data).
func DecodeUlawTo(pcm, lpcm []byte) {
	n := len(pcm)
	i, j := 0, 0
	// Process 16 samples per iteration, so that 16 independent table
	// lookups are in flight at any time, keeping the CPU pipeline busy.
	// The decoded 16bit samples are packed into uint64 values so that each
	// 4 of them are stored with a single wide store.
	// NB: the samples are converted to uint16 first, converting an int16
	// directly to uint64 would sign-extend and corrupt the packed value.
	for ; i < n-15; i, j = i+16, j+32 {
		w0 := binary.LittleEndian.Uint64(pcm[i:])
		w1 := binary.LittleEndian.Uint64(pcm[i+8:])
		q0 := uint64(uint16(ulaw2lpcm[uint8(w0)])) |
			uint64(uint16(ulaw2lpcm[uint8(w0>>8)]))<<16 |
			uint64(uint16(ulaw2lpcm[uint8(w0>>16)]))<<32 |
			uint64(uint16(ulaw2lpcm[uint8(w0>>24)]))<<48
		q1 := uint64(uint16(ulaw2lpcm[uint8(w0>>32)])) |
			uint64(uint16(ulaw2lpcm[uint8(w0>>40)]))<<16 |
			uint64(uint16(ulaw2lpcm[uint8(w0>>48)]))<<32 |
			uint64(uint16(ulaw2lpcm[uint8(w0>>56)]))<<48
		q2 := uint64(uint16(ulaw2lpcm[uint8(w1)])) |
			uint64(uint16(ulaw2lpcm[uint8(w1>>8)]))<<16 |
			uint64(uint16(ulaw2lpcm[uint8(w1>>16)]))<<32 |
			uint64(uint16(ulaw2lpcm[uint8(w1>>24)]))<<48
		q3 := uint64(uint16(ulaw2lpcm[uint8(w1>>32)])) |
			uint64(uint16(ulaw2lpcm[uint8(w1>>40)]))<<16 |
			uint64(uint16(ulaw2lpcm[uint8(w1>>48)]))<<32 |
			uint64(uint16(ulaw2lpcm[uint8(w1>>56)]))<<48
		o := lpcm[j : j+32 : j+32]
		binary.LittleEndian.PutUint64(o[0:], q0)
		binary.LittleEndian.PutUint64(o[8:], q1)
		binary.LittleEndian.PutUint64(o[16:], q2)
		binary.LittleEndian.PutUint64(o[24:], q3)
	}
	// Handle trailing samples
	for ; i < n; i, j = i+1, j+2 {
		binary.LittleEndian.PutUint16(lpcm[j:], uint16(ulaw2lpcm[pcm[i]]))
	}
}

// DecodeUlawFrame decodes a u-law PCM frame to 16bit LPCM
func DecodeUlawFrame(frame uint8) int16 {
	return ulaw2lpcm[frame]
}

// Ulaw2Alaw performs direct u-law to A-law data conversion
func Ulaw2Alaw(ulaw []byte) []byte {
	alaw := make([]byte, len(ulaw))
	Ulaw2AlawTo(ulaw, alaw)
	return alaw
}

// Ulaw2AlawTo performs direct u-law to A-law data conversion
// using an already allocated buffer provided by the user.
// The user is responsible for ensuring that the buffer is large enough (the size of the u-law data).
func Ulaw2AlawTo(ulaw, alaw []byte) {
	n := len(ulaw)
	i := 0
	// Process 16 samples per iteration, so that 16 independent table
	// lookups are in flight at any time, keeping the CPU pipeline busy.
	// The input bytes are read with two wide loads.
	for ; i < n-15; i += 16 {
		w0 := binary.LittleEndian.Uint64(ulaw[i:])
		w1 := binary.LittleEndian.Uint64(ulaw[i+8:])
		o := alaw[i : i+16]
		o[0] = ulaw2alaw[uint8(w0)]
		o[1] = ulaw2alaw[uint8(w0>>8)]
		o[2] = ulaw2alaw[uint8(w0>>16)]
		o[3] = ulaw2alaw[uint8(w0>>24)]
		o[4] = ulaw2alaw[uint8(w0>>32)]
		o[5] = ulaw2alaw[uint8(w0>>40)]
		o[6] = ulaw2alaw[uint8(w0>>48)]
		o[7] = ulaw2alaw[uint8(w0>>56)]
		o[8] = ulaw2alaw[uint8(w1)]
		o[9] = ulaw2alaw[uint8(w1>>8)]
		o[10] = ulaw2alaw[uint8(w1>>16)]
		o[11] = ulaw2alaw[uint8(w1>>24)]
		o[12] = ulaw2alaw[uint8(w1>>32)]
		o[13] = ulaw2alaw[uint8(w1>>40)]
		o[14] = ulaw2alaw[uint8(w1>>48)]
		o[15] = ulaw2alaw[uint8(w1>>56)]
	}
	// Handle trailing samples
	for ; i < n; i++ {
		alaw[i] = ulaw2alaw[ulaw[i]]
	}
}

// Ulaw2AlawFrame directly converts a u-law frame to A-law
func Ulaw2AlawFrame(frame uint8) uint8 {
	return ulaw2alaw[frame]
}
