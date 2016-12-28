# g711

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
func Alaw2Ulaw(pcm uint8) uint8
```
Alaw2Ulaw performs direct A-law to u-law frame conversion

#### func  DecodeAlaw

```go
func DecodeAlaw(pcm uint8) int16
```
DecodeAlaw decodes an A-law PCM frame to 16bit LPCM

#### func  DecodeUlaw

```go
func DecodeUlaw(pcm uint8) int16
```
DecodeUlaw decodes a u-law PCM frame to 16bit LPCM

#### func  EncodeAlaw

```go
func EncodeAlaw(lpcm int16) uint8
```
EncodeAlaw encodes a 16bit LPCM frame to G711 A-law PCM

#### func  EncodeUlaw

```go
func EncodeUlaw(lpcm int16) uint8
```
EncodeUlaw encodes a 16bit LPCM frame to G711 u-law PCM

#### func  Ulaw2Alaw

```go
func Ulaw2Alaw(pcm uint8) uint8
```
Ulaw2Alaw performs direct u-law to A-law frame conversion

#### type Reader

```go
type Reader struct {
}
```

Reader reads G711 PCM data and decodes it to 16bit LPCM or directly transcodes
between A-law and u-law

#### func  NewAlawDecoder

```go
func NewAlawDecoder(reader io.Reader, output int) (*Reader, error)
```
NewAlawDecoder returns a pointer to a Reader that decodes or trans-codes A-law
data. It takes as input the source data Reader and the output encoding fomrat.

#### func  NewUlawDecoder

```go
func NewUlawDecoder(reader io.Reader, output int) (*Reader, error)
```
NewUlawDecoder returns a pointer to a Reader that decodes or trans-codes u-law
data. It takes as input the source data Reader and the output encoding fomrat.

#### func (*Reader) Read

```go
func (r *Reader) Read(p []byte) (int, error)
```
Read decodes G711 data. Reads up to len(p) bytes into p, returns the number of
bytes read and any error encountered.

#### func (*Reader) Reset

```go
func (r *Reader) Reset(reader io.Reader)
```
Reset discards the Reader state. This permits reusing a Reader rather than
allocating a new one.

#### type Writer

```go
type Writer struct {
}
```

Writer encodes 16bit LPCM data to G711 PCM or directly transcodes between A-law
and u-law

#### func  NewAlawEncoder

```go
func NewAlawEncoder(writer io.Writer, input int) (*Writer, error)
```
NewAlawEncoder returns a pointer to a Writer that encodes data to A-law. It
takes as input the destination data Writer and the input encoding fomrat.

#### func  NewUlawEncoder

```go
func NewUlawEncoder(writer io.Writer, input int) (*Writer, error)
```
NewUlawEncoder returns a pointer to a Writer that encodes data to u-law. It
takes as input the destination data Writer and the input encoding fomrat.

#### func (*Writer) Flush

```go
func (w *Writer) Flush() (err error)
```
Flush flushes any pending data to the underlying writer.

#### func (*Writer) Reset

```go
func (w *Writer) Reset(writer io.Writer)
```
Reset discards the Writer state. This permits reusing a Writer rather than
allocating a new one.

#### func (*Writer) Write

```go
func (w *Writer) Write(p []byte) (int, error)
```
Write encodes G711 Data. Writes len(p) bytes from p to the underlying data
stream, returns the number of bytes written from p (0 <= n <= len(p)/2 due to
compression) and any error encountered that caused the write to stop early.
