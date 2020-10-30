package main

import (
	"amgine/rotor"
	"amgine/util"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// TransformFile ...
func TransformFile(parameters util.Parameters, rotors []rotor.Rotor) string {
	var outputFileName string
	inputFileName := parameters.Target
	inputFile, err := os.Open(inputFileName)
	defer inputFile.Close()

	if err == nil {
		outputFileName, err = getOutputFileName(inputFileName, inputFile, rotors, parameters)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		outputFile, err := os.Create(outputFileName)
		defer outputFile.Close()

		if err == nil {
			fileBytes := make([]byte, 1)
			if parameters.EncryptFilename {
				for _, c := range inputFileName {
					fileBytes[0] = transform(byte(c), rotors)
					outputFile.Write(fileBytes)
				}
				fileBytes[0] = transform(byte('\x00'), rotors)
				outputFile.Write(fileBytes)
			}
			for {
				n, _ := inputFile.Read(fileBytes)
				if n > 0 {
					fileBytes[0] = transform(fileBytes[0], rotors)
					outputFile.Write(fileBytes)
				} else {
					break
				}
			}
		} else {
			fmt.Println("Transform error.")
			return ""
		}
	} else {
		fmt.Println("Input file error: ", err)
		return ""
	}
	return outputFileName
}

func transform(inByte byte, rotors []rotor.Rotor) byte {
	outSlot := int(inByte)
	for i, rotor := range rotors {
		if i < len(rotors)-1 {
			outSlot = rotor.Slots[outSlot]
		}
	}

	outSlot = rotors[len(rotors)-1].Slots[outSlot]

	for i := len(rotors) - 2; i >= 0; i-- {
		for n := range rotors[i].Slots {
			if outSlot == rotors[i].Slots[n] {
				outSlot = n
				break
			}
		}
	}

	rotor.SpinRotors(rotors)
	return byte(outSlot)
}

func getOutputFileName(inputFileName string, inputFile *os.File, rotors []rotor.Rotor, parameters util.Parameters) (string, error) {
	if parameters.EncryptFilename {
		uuidWithHyphen := uuid.New()
		return strings.Replace(uuidWithHyphen.String(), "-", "", -1), nil
	} else if parameters.DecryptFilename {
		var buffer bytes.Buffer
		fileBytes := make([]byte, 1)
		for {
			n, _ := inputFile.Read(fileBytes)
			if n > 0 {
				fileBytes[0] = transform(fileBytes[0], rotors)
				if fileBytes[0] != '\x00' {
					buffer.WriteByte(fileBytes[0])
				} else {
					return buffer.String(), nil
				}
			} else {
				return "", errors.New("could not determine output file name")
			}
		}
	} else {
		amgIndex := strings.Index(inputFileName, ".amg")
		if amgIndex == -1 {
			return inputFileName + ".amg", nil
		}
		return inputFileName[:amgIndex], nil
	}
}
