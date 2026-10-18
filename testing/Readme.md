# Test data

## Sound files
```
dtmf-1234.raw     : Signed 16 bit Little Endian, Rate 8000 Hz, Mono
piano.raw         : Signed 16 bit Little Endian, Rate 8000 Hz, Mono
silence-1s.raw    : Signed 16 bit Little Endian, Rate 8000 Hz, Mono
sine-440Hz-1s.raw : Signed 16 bit Little Endian, Rate 8000 Hz, Mono
speech.alaw       : A-Law, Rate 8000 Hz, Mono
speech.raw        : Signed 16 bit Little Endian, Rate 8000 Hz, Mono
speech.ulaw       : Mu-Law, Rate 8000 Hz, Mono
````

### To playback raw data:
```
aplay -f S16_LE -r 8000 [filename]
```

### To playback G711 data:
```
aplay -f [MU_LAW|A_LAW] [filename]
```
