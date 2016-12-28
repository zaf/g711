/*
	Copyright (C) 2016 - 2017, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

	g711dec decodes 8bit G711 PCM data to 16 Bit signed LPCM raw data

*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaf/g711"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" || os.Args[1] == "--help" {
		fmt.Printf("%s Decodes 8bit G711 PCM data to raw 16 Bit signed LPCM\n", os.Args[0])
		fmt.Println("The program takes as input a list A-law or u-law encoded files")
		fmt.Println("decodes them to LPCM and saves the files with a \".raw\" extension.")
		fmt.Printf("\nUsage: %s [files]\n", os.Args[0])
		os.Exit(1)
	}
	var exitCode int
	for _, file := range os.Args[1:] {
		err := decodeG711(file)
		if err != nil {
			fmt.Println(err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func decodeG711(file string) error {
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	extension := strings.ToLower(filepath.Ext(file))
	decoder := new(g711.Decoder)
	if extension == ".alaw" || extension == ".al" {
		decoder, err = g711.NewAlawDecoder(input)
		if err != nil {
			return err
		}
	} else if extension == ".ulaw" || extension == ".ul" {
		decoder, err = g711.NewUlawDecoder(input)
		if err != nil {
			return err
		}
	} else {
		err = fmt.Errorf("Unrecognised format for file: %s", file)
		return err
	}
	outName := strings.TrimSuffix(file, filepath.Ext(file)) + ".raw"
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, decoder)
	return err
}
