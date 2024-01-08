/*
	Copyright (C) 2016 - 2024, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	g711 encodes 16bit 8kHz LPCM data to 8bit G711 PCM,
	or decodes 8bit G711 PCM to 16bit 8kHz LPCM,
	or converts between A-law and u-law formats.
	Input can be 16bit 8kHz wav or raw LPCM files, or ulaw/alaw encoded PCM data.

	Usage: g711 -in <input format> -out <output format> <input files>
	Valid formats are: alaw, ulaw, lpcm

*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaf/g711"
)

const wavHeader = 44

var (
	in  = flag.String("in", "", "input format: alaw, ulaw, lpcm")
	out = flag.String("out", "", "output format: alaw, ulaw, lpcm")
)

var formats = map[string]int{
	"alaw": g711.Alaw,
	"ulaw": g711.Ulaw,
	"lpcm": g711.Lpcm,
}

func main() {
	flag.Parse()
	var exitCode int
	for _, file := range flag.Args() {
		err := translate(file)
		if err != nil {
			fmt.Println("Error while processing", file, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func translate(file string) error {
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	inExtension := strings.ToLower(filepath.Ext(file))

	inFormat := formats[strings.ToLower(*in)]
	outFormat := formats[strings.ToLower(*out)]
	outName := strings.TrimSuffix(file, filepath.Ext(file)) + "." + strings.ToLower(*out)
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// Create a new translator, it implements io.WriteCloser
	translator, err := g711.NewCoder(outFile, inFormat, outFormat)
	if err != nil {
		os.Remove(outName)
		return err
	}
	defer translator.Close()
	if inExtension == ".wav" {
		input.Seek(wavHeader, 0) // Skip wav header
	}
	// enc/dec/transcode the input file data and write to the output file
	_, err = io.Copy(translator, input)
	outFile.Sync()
	return err
}
