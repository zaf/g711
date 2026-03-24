# g711

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/.)

Package g711 implements encoding and decoding of G711 PCM sound data.
G.711 is an ITU-T standard for audio companding.

The package exposes a high level API for encoding and decoding through
an io.WriteCloser. But also a low level API using preallocated buffers
for cases where performance and memory handling are critical.

For usage details please see the code snippet in the cmd folder.

## Constants

```golang
const (
    // Input and output formats
    Alaw = iota + 1 // Alaw G711 encoded PCM data
    Ulaw            // Ulaw G711  encoded PCM data
    Lpcm            // Lpcm 16bit signed linear data
)
```

## Types

### type [Coder](/g711.go#L32)

`type Coder struct { ... }`

Coder encodes 16bit 8000Hz LPCM data to G711 PCM, or
decodes G711 PCM data to 16bit 8000Hz LPCM data, or
directly transcodes between A-law and u-law

#### func (*Coder) [Close](/g711.go#L96)

`func (w *Coder) Close() error`

Close closes the Encoder, it implements the io.Closer interface.

#### func (*Coder) [Reset](/g711.go#L105)

`func (w *Coder) Reset(writer io.Writer) error`

Reset discards the Encoder state. This permits reusing an Encoder rather than allocating a new one.

#### func (*Coder) [Write](/g711.go#L119)

`func (w *Coder) Write(p []byte) (int, error)`

Write encodes/decodes/transcodes sound data. Writes len(p) bytes from p to the underlying data stream,
returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
that caused the write to stop early.


## Functions

### func [Alaw2Ulaw](/alaw.go#L129)

`func Alaw2Ulaw(alaw []byte) []byte`

Alaw2Ulaw performs direct A-law to u-law data conversion

### func [Alaw2UlawFrame](/alaw.go#L145)

`func Alaw2UlawFrame(frame uint8) uint8`

Alaw2UlawFrame directly converts an A-law frame to u-law

### func [Alaw2UlawTo](/alaw.go#L138)

`func Alaw2UlawTo(alaw, ulaw []byte)`

Alaw2UlawTo performs direct A-law to u-law data conversion
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (the size of the A-law data).

### func [DecodeAlaw](/alaw.go#L106)

`func DecodeAlaw(pcm []byte) []byte`

DecodeAlaw decodes A-law PCM data to 16bit LPCM

### func [DecodeAlawFrame](/alaw.go#L124)

`func DecodeAlawFrame(frame uint8) int16`

DecodeAlawFrame decodes an A-law PCM frame to 16bit LPCM

### func [DecodeAlawTo](/alaw.go#L115)

`func DecodeAlawTo(pcm, lpcm []byte)`

DecodeAlawTo decodes A-law PCM data to 16bit LPCM
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (double the size of the PCM data).

### func [DecodeUlaw](/ulaw.go#L110)

`func DecodeUlaw(pcm []byte) []byte`

DecodeUlaw decodes u-law PCM data to 16bit LPCM

### func [DecodeUlawFrame](/ulaw.go#L128)

`func DecodeUlawFrame(frame uint8) int16`

DecodeUlawFrame decodes a u-law PCM frame to 16bit LPCM

### func [DecodeUlawTo](/ulaw.go#L119)

`func DecodeUlawTo(pcm, lpcm []byte)`

DecodeUlawTo decodes u-law PCM data to 16bit LPCM
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (double the size of the PCM data).

### func [EncodeAlaw](/alaw.go#L74)

`func EncodeAlaw(lpcm []byte) []byte`

EncodeAlaw encodes 16bit LPCM data to G711 A-law PCM

### func [EncodeAlawFrame](/alaw.go#L90)

`func EncodeAlawFrame(frame int16) uint8`

EncodeAlawFrame encodes a 16bit LPCM frame to G711 A-law PCM

### func [EncodeAlawTo](/alaw.go#L83)

`func EncodeAlawTo(lpcm, alaw []byte)`

EncodeAlawTo encodes 16bit LPCM data to G711 A-law PCM
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (half the size of the LPCM data).

### func [EncodeUlaw](/ulaw.go#L79)

`func EncodeUlaw(lpcm []byte) []byte`

EncodeUlaw encodes 16bit LPCM data to G711 u-law PCM

### func [EncodeUlawFrame](/ulaw.go#L95)

`func EncodeUlawFrame(frame int16) uint8`

EncodeUlawFrame encodes a 16bit LPCM frame to G711 u-law PCM

### func [EncodeUlawTo](/ulaw.go#L88)

`func EncodeUlawTo(lpcm, ulaw []byte)`

EncodeUlawTo encodes 16bit LPCM data to G711 u-law PCM
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (half the size of the LPCM data).

### func [Ulaw2Alaw](/ulaw.go#L133)

`func Ulaw2Alaw(ulaw []byte) []byte`

Ulaw2Alaw performs direct u-law to A-law data conversion

### func [Ulaw2AlawFrame](/ulaw.go#L149)

`func Ulaw2AlawFrame(frame uint8) uint8`

Ulaw2AlawFrame directly converts a u-law frame to A-law

### func [Ulaw2AlawTo](/ulaw.go#L142)

`func Ulaw2AlawTo(ulaw, alaw []byte)`

Ulaw2AlawTo performs direct u-law to A-law data conversion
using an already allocated buffer provided by the user.
The user is responsible for ensuring that the buffer is large enough (the size of the A-law data).

---
