/*
	Copyright (C) 2016 - 2024, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	Package g711 implements encoding and decoding of G711 PCM sound data.
	G.711 is an ITU-T standard for audio companding.
*/

package g711

import "math/bits"

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

// EncodeUlaw encodes 16bit LPCM data to G711 u-law PCM
func EncodeUlaw(lpcm []byte) []byte {
	ulaw := make([]byte, len(lpcm)>>1)
	for i := 0; i < len(lpcm)-1; i += 2 {
		ulaw[i>>1] = EncodeUlawFrame(int16(lpcm[i]) | int16(lpcm[i+1])<<8)
	}
	return ulaw
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
	for i := 0; i < len(pcm); i++ {
		frame := ulaw2lpcm[pcm[i]]
		lpcm[i*2] = byte(frame)
		lpcm[i*2+1] = byte(frame >> 8)
	}
	return lpcm
}

// DecodeUlawFrame decodes a u-law PCM frame to 16bit LPCM
func DecodeUlawFrame(frame uint8) int16 {
	return ulaw2lpcm[frame]
}

// Ulaw2Alaw performs direct u-law to A-law data conversion
func Ulaw2Alaw(ulaw []byte) []byte {
	alaw := make([]byte, len(ulaw))
	for i := 0; i < len(alaw); i++ {
		alaw[i] = ulaw2alaw[ulaw[i]]
	}
	return alaw
}

// Ulaw2AlawFrame directly converts a u-law frame to A-law
func Ulaw2AlawFrame(frame uint8) uint8 {
	return ulaw2alaw[frame]
}
