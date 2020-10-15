package main

import (
	"amgine/rotor"
	"amgine/util"
	"fmt"
	"os"
	"strings"
)

// TransformFile ...
func TransformFile(parameters util.Parameters, rotors []rotor.Rotor) {
	inputFileName := parameters.Target
	inputData, err := os.Open(inputFileName)
	defer inputData.Close()

	if err == nil {
		outputFileName := getOutputFileName(inputFileName)
		outputData, err := os.Create(outputFileName)
		defer outputData.Close()

		if err == nil {
			for {
				fileBytes := make([]byte, 1)
				n, _ := inputData.Read(fileBytes)
				if n > 0 {
					fileBytes[0] = transform(fileBytes[0], rotors)
					outputData.Write(fileBytes)
				} else {
					break
				}
			}
		}
	}
	if err != nil {
		fmt.Println("File error: ", err)
	}
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

	spinRotors(rotors)
	return byte(outSlot)
}

func spinRotors(rotors []rotor.Rotor) {
	rotorIndex := 0
	rotors[rotorIndex].IncrementCurrent()
	rotors[rotorIndex].Rotate()

	for rotors[rotorIndex].Current == rotors[rotorIndex].Notch && rotorIndex < len(rotors)-1 {
		rotorIndex++
		rotors[rotorIndex].IncrementCurrent()
		rotors[rotorIndex].Rotate()
	}
}

func getOutputFileName(inputFileName string) string {
	amgIndex := strings.Index(inputFileName, ".amg")
	if amgIndex == -1 {
		return inputFileName + ".amg"
	}
	return inputFileName[:amgIndex]
}
