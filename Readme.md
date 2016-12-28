# g711
--
    import "github.com/zaf/g711"

Package g711 implements encoding and decoding of G711.0 compressed sound data.
G.711 is an ITU-T standard for audio companding.

For usage details please see the code snippets in the cmd folder.

## Usage

```go
const (
	// Input and output formats
	Alaw = iota // Alaw G711 encoded PCM data
	Ulaw        // Ulaw G711  encoded PCM data
	Lpcm        // Lpcm 16bit signed linear data
)
```

#### func  Alaw2Ulaw

```go
func Alaw2Ulaw(alaw []byte) []byte
```
Alaw2Ulaw performs direct A-law to u-law data conversion

#### func  DecodeAlaw

```go
func DecodeAlaw(pcm []byte) []byte
```
DecodeAlaw decodes A-law PCM data to 16bit LPCM

#### func  DecodeUlaw

```go
func DecodeUlaw(pcm []byte) []byte
```
DecodeUlaw decodes u-law PCM data to 16bit LPCM

#### func  EncodeAlaw

```go
func EncodeAlaw(lpcm []byte) []byte
```
EncodeAlaw encodes 16bit LPCM data to G711 A-law PCM

#### func  EncodeUlaw

```go
func EncodeUlaw(lpcm []byte) []byte
```
EncodeUlaw encodes 16bit LPCM data to G711 u-law PCM

#### func  Ulaw2Alaw

```go
func Ulaw2Alaw(ulaw []byte) []byte
```
Ulaw2Alaw performs direct u-law to A-law data conversion

#### type Decoder

```go
type Decoder struct {
}
```

Decoder implements an io.Reader interface. It reads G711 PCM data and decodes it
to 16bit LPCM

#### func  NewAlawDecoder

```go
func NewAlawDecoder(reader io.Reader) (*Decoder, error)
```
NewAlawDecoder returns a pointer to a Decoder. It takes as input the source data
Reader.

#### func  NewUlawDecoder

```go
func NewUlawDecoder(reader io.Reader) (*Decoder, error)
```
NewUlawDecoder returns a pointer to a Decoder It takes as input the source data
Reader.

#### func (*Decoder) Read

```go
func (r *Decoder) Read(p []byte) (int, error)
```
Read decodes G711 data. Reads up to len(p) bytes into p, returns the number of
bytes read and any error encountered.

#### func (*Decoder) Reset

```go
func (r *Decoder) Reset(reader io.Reader)
```
Reset discards the Decoder state. This permits reusing a Decoder rather than
allocating a new one.

#### type Encoder

```go
type Encoder struct {
}
```

Encoder implements an io.Writer interface. It encodes 16bit LPCM data to G711
PCM or directly transcodes between A-law and u-law

#### func  NewAlawEncoder

```go
func NewAlawEncoder(writer io.Writer, input int) (*Encoder, error)
```
NewAlawEncoder returns a pointer to an Encoder. It takes as input the
destination data Writer and the input encoding format.

#### func  NewUlawEncoder

```go
func NewUlawEncoder(writer io.Writer, input int) (*Encoder, error)
```
NewUlawEncoder returns a pointer to an Encoder. It takes as input the
destination data Writer and the input encoding format.

#### func (*Encoder) Reset

```go
func (w *Encoder) Reset(writer io.Writer)
```
Reset discards the Encoder state. This permits reusing an Encoder rather than
allocating a new one.

#### func (*Encoder) Write

```go
func (w *Encoder) Write(p []byte) (int, error)
```
Write encodes G711 Data. Writes len(p) bytes from p to the underlying data
stream, returns the number of bytes written from p (0 <= n <= len(p)) and any
error encountered that caused the write to stop early.
