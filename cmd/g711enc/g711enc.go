/*
	Copyright (C) 2016 - 2017, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	g711enc encodes 16bit 8kHz LPCM data to 8bit G711 PCM.
	It works with wav or raw files as input.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaf/g711"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] == "help" || os.Args[1] == "--help" || (os.Args[1] != "ulaw" && os.Args[1] != "alaw") {
		fmt.Printf("%s Encodes 16bit 8kHz LPCM data to 8bit G711 PCM\n", os.Args[0])
		fmt.Println("The program takes as input a list or wav or raw files, encodes them")
		fmt.Println("to G711 PCM and saves them with the proper extension.")
		fmt.Printf("\nUsage: %s [encoding format] [files]\n", os.Args[0])
		fmt.Println("encoding format can be either alaw or ulaw")
		os.Exit(1)
	}
	var exitCode int
	format := os.Args[1]
	for _, file := range os.Args[2:] {
		err := encodeG711(file, format)
		if err != nil {
			fmt.Println(err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func encodeG711(file, format string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	extension := strings.ToLower(filepath.Ext(file))
	if extension != ".wav" && extension != ".raw" && extension != ".sln" {
		err = fmt.Errorf("Unrecognised format for input file: %s", file)
		return err
	}
	outName := strings.TrimSuffix(file, filepath.Ext(file)) + "." + format
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	encoder := new(g711.Encoder)
	if format == "alaw" {
		encoder, err = g711.NewAlawEncoder(outFile, g711.Lpcm)
		if err != nil {
			return err
		}
	} else {
		encoder, err = g711.NewUlawEncoder(outFile, g711.Lpcm)
		if err != nil {
			return err
		}
	}
	if extension == ".wav" {
		_, err = encoder.Write(input[44:]) // Skip WAV header
		return err
	}
	_, err = encoder.Write(input)
	return err
}
