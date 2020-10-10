package main

import (
	"amgine/rotor"
	"fmt"
	"os"
	"strings"
)

const (
	// ENCRYPT ...
	ENCRYPT = iota
	// DECRYPT ...
	DECRYPT = iota
)

// TransformFile ...
func TransformFile(inputFileName string, mode int, configuration AppConfig, rotorStore rotor.RotorStore) {
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
					if mode == DECRYPT {
						fileBytes[0] = decrypt(fileBytes[0], configuration, rotorStore)
					} else {
						fileBytes[0] = encrypt(fileBytes[0], configuration, rotorStore)
					}
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

func encrypt(inByte byte, configuration AppConfig, rotorStore rotor.RotorStore) byte {
	outSlot := int(inByte)
	for i := range configuration.Rotors {
		outSlot = rotorStore.Rotors[configuration.Rotors[i].Rid].Slots[outSlot]
	}

	outSlot = rotorStore.ReturnRotors[configuration.ReturnRotor.Rid].Slots[outSlot]

	for i := len(configuration.Rotors) - 1; i >= 0; i-- {
		for n := range rotorStore.Rotors[configuration.Rotors[i].Rid].Slots {
			if outSlot == rotorStore.Rotors[configuration.Rotors[i].Rid].Slots[n] {
				outSlot = n
				break
			}
		}
	}
	return byte(outSlot)
}

func decrypt(inByte byte, configuration AppConfig, rotorStore rotor.RotorStore) byte {
	outSlot := int(inByte)
	for i := range configuration.Rotors {
		outSlot = rotorStore.Rotors[configuration.Rotors[i].Rid].Slots[outSlot]
	}

	for n := range rotorStore.ReturnRotors[configuration.ReturnRotor.Rid].Slots {
		if outSlot == rotorStore.ReturnRotors[configuration.ReturnRotor.Rid].Slots[n] {
			outSlot = n
			break
		}
	}

	for i := len(configuration.Rotors) - 1; i >= 0; i-- {
		for n := range rotorStore.Rotors[configuration.Rotors[i].Rid].Slots {
			if outSlot == rotorStore.Rotors[configuration.Rotors[i].Rid].Slots[n] {
				outSlot = n
				break
			}
		}
	}
	return byte(outSlot)
}

func getOutputFileName(inputFileName string) string {
	amgIndex := strings.Index(inputFileName, ".amg")
	if amgIndex == -1 {
		return inputFileName + ".amg"
	}
	return inputFileName[:amgIndex]
}
