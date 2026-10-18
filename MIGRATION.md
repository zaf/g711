# Changes in v2
This version reworks the public API of the g711 package, adds zero-allocation
functions. **The changes are breaking** for code that used the streaming API (`Encoder`/`Decoder`) or the
`g711enc`/`g711dec` binaries. The low level functions are unaffected.

## New features

- **One unified type for everything.** `g711.NewCoder(writer, input, output)`
  handles encoding (LPCM -> A-law/u-law), decoding (A-law/u-law -> LPCM) *and*
  direct A-law/u-law transcoding in a single `io.WriteCloser`.
- **Zero allocation functions** for performance critical code. The caller
  provides the output buffer:
  - `EncodeAlawTo(lpcm, alaw []byte)` / `EncodeUlawTo(lpcm, ulaw []byte)`
  - `DecodeAlawTo(pcm, lpcm []byte)` / `DecodeUlawTo(pcm, lpcm []byte)`
  - `Alaw2UlawTo(alaw, ulaw []byte)` / `Ulaw2AlawTo(ulaw, alaw []byte)`
- **Speed optimizations**: encoding now uses precomputed lookup tables
  (instead of computing each frame), the hot loops are unrolled to keep
  multiple table lookups in flight, and decoding packs 8 samples into `uint64`
  values for wide stores.
- **Robust streaming**: a 16-bit LPCM sample split across two `Write` calls is
  buffered and completed by the next call. The old API returned an error on an
  odd number of LPCM bytes, now the leftover byte is carried over and only
  discarded at `Close()`.

## API changes

| Master (old)                            | This branch (new)                              |
|-----------------------------------------|-------------------------------------------------|
| `NewAlawDecoder(reader)` / `NewUlawDecoder(reader)` -> `Decoder` (an `io.Reader`) | `NewCoder(writer, Alaw, Lpcm)` / `NewCoder(writer, Ulaw, Lpcm)` (an `io.WriteCloser`) |
| `NewAlawEncoder(writer, input)`         | `NewCoder(writer, input, Alaw)`                 |
| `NewUlawEncoder(writer, input)`         | `NewCoder(writer, input, Ulaw)`                 |
| `(*Encoder).Reset(writer)`              | `(*Coder).Reset(writer)`                        |
| `(*Decoder).Reset(reader)`              | gone - create a new `Coder` instead            |
| (none)                                 | `(*Coder).Close()` (new)                        |

The decoding direction is inverted: it used to *pull* from an `io.Reader`
(`Decoder.Read`), now you *push* into a `Coder` with `io.Copy(coder, src)`.

The low level API is unchanged: `EncodeAlaw`, `EncodeUlaw`, `DecodeAlaw`,
`DecodeUlaw`, `Alaw2Ulaw`, `Ulaw2Alaw` and the `*Frame` functions keep their
signatures. `min` supported Go version is now **1.21** (was 1.20).

## How to update existing v1 code

**Decoding** (G711 -> LPCM), old:

```go
decoder, err := g711.NewAlawDecoder(src) // or NewUlawDecoder
...
io.Copy(dst, decoder)
```

new - arguments are `(writer, input, output)`:

```go
coder, err := g711.NewCoder(dst, g711.Alaw, g711.Lpcm) // or g711.Ulaw
...
io.Copy(coder, src)
defer coder.Close()
```

**Encoding** (LPCM -> G711), old:

```go
encoder, err := g711.NewAlawEncoder(dst, g711.Lpcm) // or NewUlawEncoder
...
io.Copy(encoder, src)
```

new:

```go
coder, err := g711.NewCoder(dst, g711.Lpcm, g711.Alaw) // or g711.Ulaw
...
io.Copy(coder, src)
defer coder.Close()
```

**Transcoding** (A-law -> u-law), old:

```go
encoder, err := g711.NewUlawEncoder(dst, g711.Alaw) // u-law output, A-law input
...
io.Copy(encoder, src)
```

new, mind the order, input format first:

```go
coder, err := g711.NewCoder(dst, g711.Alaw, g711.Ulaw) // A-law in, u-law out
...
io.Copy(coder, src)
defer coder.Close()
```

Notes:

- Always call `Close()` when done; it discards any unpaired LPCM byte. A
  `Close`d (or zero value) `Coder` cannot be reused, `Reset` will error, so
  allocate a fresh one.
- `Write` returns an error only when the underlying writer fails or the coder
  is uninitialized/closed. Odd-length LPCM input no longer errors.
- Use the `*To` functions with a `sync.Pool` or a reused buffer to avoid
  per-call allocations on the hot path.
