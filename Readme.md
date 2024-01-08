# g711

Package g711 implements encoding and decoding of G711 PCM sound data.
G.711 is an ITU-T standard for audio companding.

## Usage

See code examples in `cmd/` folder.

## Constants

```golang
const (
    // Input and output formats
    Alaw = iota // Alaw G711 encoded PCM data
    Ulaw        // Ulaw G711  encoded PCM data
    Lpcm        // Lpcm 16bit signed linear data
)
```
## Types

### type [Decoder](/g711.go#L30)

`type Decoder struct { ... }`

Decoder reads G711 PCM data and decodes it to 16bit 8000Hz LPCM

#### func [NewAlawDecoder](/g711.go#L46)

`func NewAlawDecoder(reader io.Reader) (*Decoder, error)`

NewAlawDecoder returns a pointer to a Decoder that implements an io.Reader.
It takes as input the source data Reader.

#### func [NewUlawDecoder](/g711.go#L59)

`func NewUlawDecoder(reader io.Reader) (*Decoder, error)`

NewUlawDecoder returns a pointer to a Decoder that implements an io.Reader.
It takes as input the source data Reader.

#### func (*Decoder) [Close](/g711.go#L107)

`func (r *Decoder) Close() error`

Close closes the Decoder, it implements the io.Closer interface.

#### func (*Decoder) [Read](/g711.go#L138)

`func (r *Decoder) Read(p []byte) (i int, err error)`

Read decodes G711 data. Reads up to len(p) bytes into p, returns the number
of bytes read and any error encountered.

#### func (*Decoder) [Reset](/g711.go#L119)

`func (r *Decoder) Reset(reader io.Reader) error`

Reset discards the Decoder state. This permits reusing a Decoder rather than allocating a new one.

### type [Encoder](/g711.go#L37)

`type Encoder struct { ... }`

Encoder encodes 16bit 8000Hz LPCM data to G711 PCM or
directly transcodes between A-law and u-law

#### func [NewAlawEncoder](/g711.go#L72)

`func NewAlawEncoder(writer io.Writer, input int) (*Encoder, error)`

NewAlawEncoder returns a pointer to an Encoder that implements an io.Writer.
It takes as input the destination data Writer and the input encoding format.

#### func [NewUlawEncoder](/g711.go#L90)

`func NewUlawEncoder(writer io.Writer, input int) (*Encoder, error)`

NewUlawEncoder returns a pointer to an Encoder that implements an io.Writer.
It takes as input the destination data Writer and the input encoding format.

#### func (*Encoder) [Close](/g711.go#L113)

`func (w *Encoder) Close() error`

Close closes the Encoder, it implements the io.Closer interface.

#### func (*Encoder) [Reset](/g711.go#L128)

`func (w *Encoder) Reset(writer io.Writer) error`

Reset discards the Encoder state. This permits reusing an Encoder rather than allocating a new one.

#### func (*Encoder) [Write](/g711.go#L152)

`func (w *Encoder) Write(p []byte) (i int, err error)`

Write encodes G711 Data. Writes len(p) bytes from p to the underlying data stream,
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
