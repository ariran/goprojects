package main

import (
	"amgine/rotor"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type AppConfig struct {
	Rotors      []RotorsConfig
	ReturnRotor RotorsConfig
}

type RotorsConfig struct {
	Seq  int
	Rid  string
	Curr int
}

func (c AppConfig) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Rotors:\n"))
	for _, r := range c.Rotors {
		buffer.WriteString(fmt.Sprintf("%v", r.String()))
	}
	buffer.WriteString(fmt.Sprintf("ReturnRotor:\n"))
	buffer.WriteString(fmt.Sprintf("%v", c.ReturnRotor.String()))
	return buffer.String()
}

func (c RotorsConfig) String() string {
	return fmt.Sprintf("Seq: %v, Rid: %v, Curr: %v\n", c.Seq, c.Rid, c.Curr)
}

func loadConfiguration() AppConfig {
	configFile, err := os.Open(os.Getenv("AMG_CONFIGFILE"))
	if err != nil {
		fmt.Println("Could not open configuration file!")
		panic(err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	configuration := AppConfig{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Could not read configuration!")
		panic(err)
	}

	sort.Slice(configuration.Rotors, func(i, j int) bool {
		return configuration.Rotors[i].Seq < configuration.Rotors[j].Seq
	})

	return configuration
}

func createTestFile() {
	outputData, _ := os.Create("C:\\Temp\\testinput.dat")
	defer outputData.Close()
	fileBytes := []byte{5, 8, 9}
	outputData.Write(fileBytes)
}

func main() {
	// createTestFile()

	configuration := loadConfiguration()
	// fmt.Println(configuration)

	// rotor.CreateNewRotor()
	// rotor.CreateNewReturnRotor()

	rotorStore := rotor.LoadRotorStore()
	// fmt.Println(rotorStore)

	// TransformFile("C:\\Temp\\input.txt", ENCRYPT, configuration, rotorStore)
	TransformFile("C:\\Temp\\input.txt.amg", DECRYPT, configuration, rotorStore)
}
