package rotor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const RotorLength = 256

type Rotor struct {
	Rotor   [RotorLength]int
	Notch   int
	Current int
}

func (r Rotor) String() string {
	return fmt.Sprintf("Rotor: %v", r.Rotor)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func CreateNewRotor() {
	fmt.Println("Creating new rotor...")

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rotor Rotor
	for i := range rotor.Rotor {
		rotor.Rotor[i] = -1
	}
	current := 0

	for current < RotorLength {
		var n int
		for rotor.Rotor[current] == -1 {
			n = r1.Intn(RotorLength)
			found := false
			for i := 0; i <= current; i++ {
				if n == rotor.Rotor[i] {
					found = true
					break
				}
			}
			if !found {
				rotor.Rotor[current] = n
				current++
				break
			}
		}
	}

	for n := range rotor.Rotor {
		fmt.Print(rotor.Rotor[n], " ")
	}
	fmt.Println()
	WriteRotor(rotor)
}

func WriteRotor(rotor Rotor) {
	var outbuffer bytes.Buffer
	encoder := gob.NewEncoder(&outbuffer)

	err := encoder.Encode(rotor)
	check(err)

	f, err := os.Create("c:\\temp\\rotors.dat")
	check(err)
	defer f.Close()

	n, err := f.Write(outbuffer.Bytes())
	check(err)
	fmt.Printf("Wrote Rotor with %d bytes\n", n)
}

func LoadRotor() {
	rotorData, err := os.Open("c:\\temp\\rotors.dat")
	check(err)

	decoder := gob.NewDecoder(rotorData)
	var rotor Rotor
	err = decoder.Decode(&rotor)

	fmt.Println(rotor)
}
