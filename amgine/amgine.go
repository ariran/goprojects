package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type AppConfig struct {
	Rotors []RotorsConfig
}

type RotorsConfig struct {
	Seq  int
	Rid  string
	Curr int
}

func (c AppConfig) String() string {
	var buffer bytes.Buffer
	for _, r := range c.Rotors {
		buffer.WriteString(fmt.Sprintf("%v", r.String()))
	}
	return buffer.String()
}

func (c RotorsConfig) String() string {
	return fmt.Sprintf("Seq: %v, Rid: %v, Curr: %v\n", c.Seq, c.Rid, c.Curr)
}

func loadConfiguration() {
	configFile, err := os.Open("c:\\temp\\amgine.json")
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

	fmt.Println(configuration)
}

func main() {
	loadConfiguration()
	// rotor.CreateNewRotor()
	// rotors := rotor.LoadRotors()
	// fmt.Println(rotors)
}
