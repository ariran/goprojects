package main

import (
	"amgine/rotor"
	"amgine/util"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

func getParameters() (util.Parameters, error) {
	var parameters util.Parameters
	index := 1
	if len(os.Args) > 1 {
		for os.Args[index][0:1] == "-" {
			if os.Args[index] == "-f" && len(os.Args) > 2 {
				index++
				parameters.ConfigFile = os.Args[index]
				index++
			} else if os.Args[index] == "-s" && len(os.Args) > 2 {
				index++
				parameters.RotorStore = os.Args[index]
				index++
			} else if os.Args[index] == "-n" {
				parameters.EncryptFilename = true
				index++
			} else if os.Args[index] == "-N" {
				parameters.DecryptFilename = true
				index++
			}
		}
		if parameters.EncryptFilename && parameters.DecryptFilename {
			return parameters, errors.New("invalid command line")
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

func loadConfiguration(parameters util.Parameters) util.AppConfig {
	configFileName := os.Getenv("AMG_CONFIGFILE")
	if len(parameters.ConfigFile) > 0 {
		configFileName = parameters.ConfigFile
	}
	configFile, err := os.Open(configFileName)
	if err != nil {
		fmt.Println("Could not open configuration file!")
		panic(err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	configuration := util.AppConfig{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Could not read configuration!")
		panic(err)
	}
	configuration.ConfigFile = configFileName

	sort.Slice(configuration.Rotors, func(i, j int) bool {
		return configuration.Rotors[i].Seq < configuration.Rotors[j].Seq
	})

	return configuration
}

func createRotors(configuration util.AppConfig, rotorStore rotor.RotorStore) []rotor.Rotor {
	rotors := make([]rotor.Rotor, len(configuration.Rotors)+1)

	for i, r := range configuration.Rotors {
		r1 := rotorStore.Rotors[r.Rid]
		r2 := rotor.Rotor{r1.Slots, r1.Notch, r.Curr}
		r2.InitializeRotorPosition()
		rotors[i] = r2
	}

	rr1 := rotorStore.ReturnRotors[configuration.ReturnRotor.Rid]
	rr := rotor.Rotor{rr1.Slots, rr1.Notch, configuration.ReturnRotor.Curr}
	rr.InitializeRotorPosition()
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
	helpString := "Amgine encrypts or decrypts the given file.\n\n" +
		"  Usage: amgine [options] {command} [file]\n\n" +
		"  Options:\n" +
		"     -f <optionsfile>  Makes amgine use the given options file instead of the one\n" +
		"                       defined by the environment variable AMG_CONFIGFILE.\n" +
		"     -s <rotorstore>   Makes amgine use the given file as active rotor store instead\n" +
		"                       of the one defined by the environment variable AMG_ROTORSTORE.\n" +
		"     -n                Encrypts also the original file name and includes it in the output file.\n" +
		"                       Generates an UUID for output file name.\n" +
		"     -N                Decrypts the original file name from the input file.\n\n" +
		"  Commands:\n" +
		"     help              Show this help.\n" +
		"     showconfig        Show configuration as defined in the given configuration file\n" +
		"                       or in the one defined by the environment variable AMG_CONFIGFILE.\n" +
		"     showstore         Shows the rotors in the currently active rotor store.\n" +
		"     newrotor          Creates a new rotor and stores it in the rotor store.\n" +
		"     newreturnrotor    Creates a new return rotor and stores it in the rotor store.\n" +
		"     transform <file>  Encrypts the given file if -n is specified or if the file name\n" +
		"                       does not end with .amg.\n" +
		"                       Decrypts the given file if -N is specified or if the file name\n" +
		"                       ends with .amg.\n\n" +
		"  Environment variables:\n" +
		"     AMG_CONFIGFILE    Specifies the default configuration file.\n" +
		"     AMG_ROTORSTORE    Specifies the default active rotor store.\n"
	fmt.Println(helpString)
}

func main() {
	parameters, err := getParameters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if parameters.Command == "help" {
		printHelp()
		os.Exit(1)
	}

	configuration := loadConfiguration(parameters)
	rotorStore := rotor.LoadRotorStore(parameters)
	rotors := createRotors(configuration, rotorStore)

	if parameters.Command == "showconfig" {
		fmt.Println("Config file: ", configuration.ConfigFile)
		fmt.Println("Rotor store: ", rotorStore.FileName)
		fmt.Println(configuration)
	} else if parameters.Command == "showstore" {
		fmt.Println(rotorStore)
	} else if parameters.Command == "newrotor" {
		rotor.CreateNewRotor(parameters)
		fmt.Println("Rotor created.")
	} else if parameters.Command == "newreturnrotor" {
		rotor.CreateNewReturnRotor(parameters)
		fmt.Println("Return rotor created.")
	} else if parameters.Command == "transform" {
		t0 := time.Now()
		outFile := TransformFile(parameters, rotors)
		if len(outFile) > 0 {
			fmt.Printf("Transformation done %v seconds. (%v)\n", time.Since(t0).Seconds(), outFile)
		}
	}
}
