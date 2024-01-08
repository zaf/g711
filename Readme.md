# g711

Package g711 implements encoding and decoding of G711 PCM sound data.
G.711 is an ITU-T standard for audio companding.

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

#### func [NewCoder](/g711.go#L40)

`func NewCoder(writer io.Writer, input, output int) (*Coder, error)`

NewCoder returns a pointer to a Coder that implements an io.WriteCloser.
It takes as input the destination data Writer and the input/output encoding formats.

#### func (*Coder) [Close](/g711.go#L90)

`func (w *Coder) Close() error`

Close closes the Encoder, it implements the io.Closer interface.

#### func (*Coder) [Reset](/g711.go#L96)

`func (w *Coder) Reset(writer io.Writer) error`

Reset discards the Encoder state. This permits reusing an Encoder rather than allocating a new one.

#### func (*Coder) [Write](/g711.go#L107)

`func (w *Coder) Write(p []byte) (int, error)`

Write encodes/decodes/transcodes sound data. Writes len(p) bytes from p to the underlying data stream,
returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
that caused the write to stop early.

## Functions

### func [Alaw2Ulaw](/alaw.go#L115)

`func Alaw2Ulaw(alaw []byte) []byte`

Alaw2Ulaw performs direct A-law to u-law data conversion

### func [Alaw2UlawFrame](/alaw.go#L124)

`func Alaw2UlawFrame(frame uint8) uint8`

Alaw2UlawFrame directly converts an A-law frame to u-law

### func [DecodeAlaw](/alaw.go#L99)

`func DecodeAlaw(pcm []byte) []byte`

DecodeAlaw decodes A-law PCM data to 16bit LPCM

### func [DecodeAlawFrame](/alaw.go#L110)

`func DecodeAlawFrame(frame uint8) int16`

DecodeAlawFrame decodes an A-law PCM frame to 16bit LPCM

### func [DecodeUlaw](/ulaw.go#L103)

`func DecodeUlaw(pcm []byte) []byte`

DecodeUlaw decodes u-law PCM data to 16bit LPCM

### func [DecodeUlawFrame](/ulaw.go#L114)

`func DecodeUlawFrame(frame uint8) int16`

DecodeUlawFrame decodes a u-law PCM frame to 16bit LPCM

### func [EncodeAlaw](/alaw.go#L74)

`func EncodeAlaw(lpcm []byte) []byte`

EncodeAlaw encodes 16bit LPCM data to G711 A-law PCM

### func [EncodeAlawFrame](/alaw.go#L83)

`func EncodeAlawFrame(frame int16) uint8`

EncodeAlawFrame encodes a 16bit LPCM frame to G711 A-law PCM

### func [EncodeUlaw](/ulaw.go#L79)

`func EncodeUlaw(lpcm []byte) []byte`

EncodeUlaw encodes 16bit LPCM data to G711 u-law PCM

### func [EncodeUlawFrame](/ulaw.go#L88)

`func EncodeUlawFrame(frame int16) uint8`

EncodeUlawFrame encodes a 16bit LPCM frame to G711 u-law PCM

### func [Ulaw2Alaw](/ulaw.go#L119)

`func Ulaw2Alaw(ulaw []byte) []byte`

Ulaw2Alaw performs direct u-law to A-law data conversion

### func [Ulaw2AlawFrame](/ulaw.go#L128)

`func Ulaw2AlawFrame(frame uint8) uint8`

Ulaw2AlawFrame directly converts a u-law frame to A-law

---
