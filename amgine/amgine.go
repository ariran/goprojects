package main

import (
	"amgine/rotor"
	"bytes"
	"encoding/json"
	"errors"
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

type Parameters struct {
	ConfigFile string
	Command    string
	Target     string
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

func (p Parameters) String() string {
	return fmt.Sprintf("ConfigFile: %v, Command: %v, Target: %v", p.ConfigFile, p.Command, p.Target)
}

func getParameters() (Parameters, error) {
	var parameters Parameters
	index := 1
	if len(os.Args) > 1 {
		if os.Args[index] == "-f" && len(os.Args) > 2 {
			index++
			parameters.ConfigFile = os.Args[index]
			index++
		}
		if len(os.Args) > index {
			if os.Args[index] == "help" ||
				os.Args[index] == "showconfig" ||
				os.Args[index] == "showstore" ||
				os.Args[index] == "newrotor" ||
				os.Args[index] == "newreturnrotor" {
				parameters.Command = os.Args[index]
				return parameters, nil
			} else if len(os.Args) > index+1 {
				if os.Args[index] == "transform" {
					parameters.Command = os.Args[index]
					index++
					parameters.Target = os.Args[index]
					return parameters, nil
				}
			}
		}
	}
	return parameters, errors.New("invalid command line")
}

func loadConfiguration() AppConfig {
	// TODO: Use override config file (-f switch)
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

func createRotors(configuration AppConfig, rotorStore rotor.RotorStore) []rotor.Rotor {
	rotors := make([]rotor.Rotor, len(configuration.Rotors)+1)

	for i, r := range configuration.Rotors {
		r1 := rotorStore.Rotors[r.Rid]
		r2 := rotor.Rotor{r1.Slots, r1.Notch, r1.Current}
		rotors[i] = r2
	}

	rr1 := rotorStore.ReturnRotors[configuration.ReturnRotor.Rid]
	rr := rotor.Rotor{rr1.Slots, rr1.Notch, rr1.Current}
	rotors[len(rotors)-1] = rr

	return rotors
}

func createTestFile() {
	outputData, _ := os.Create("C:\\Temp\\testinput.dat")
	defer outputData.Close()
	fileBytes := []byte{5, 8, 9}
	outputData.Write(fileBytes)
}

func printHelp() {
	fmt.Println("Help not implemented.")
}

func main() {
	parameters, err := getParameters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	configuration := loadConfiguration()
	rotorStore := rotor.LoadRotorStore()
	rotors := createRotors(configuration, rotorStore)

	if parameters.Command == "help" {
		printHelp()
	} else if parameters.Command == "showconfiguration" {
		fmt.Println(configuration)
	} else if parameters.Command == "showstore" {
		fmt.Println(rotorStore)
	} else if parameters.Command == "newrotor" {
		rotor.CreateNewRotor()
		fmt.Println("Rotor created.")
	} else if parameters.Command == "newreturnrotor" {
		rotor.CreateNewReturnRotor()
		fmt.Println("Return rotor created.")
	} else if parameters.Command == "transform" {
		TransformFile(parameters.Target, rotors)
		fmt.Println("Transformation done.")
	}
}
