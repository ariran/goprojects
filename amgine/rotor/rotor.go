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
	Slots   [RotorLength]int
	Notch   int
	Current int
}

type RotorStore struct {
	Rotors       map[string]Rotor
	ReturnRotors map[string]Rotor
}

func (r Rotor) String() string {
	return fmt.Sprintf("Slots: %v\nNotch: %v", r.Slots, r.Notch)
}

func (s RotorStore) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Rotors:\n"))
	for _, r := range s.Rotors {
		buffer.WriteString(fmt.Sprintf("%v\n", r))
	}
	buffer.WriteString(fmt.Sprintf("ReturnRotors:\n"))
	for _, r := range s.ReturnRotors {
		buffer.WriteString(fmt.Sprintf("%v\n", r))
	}
	return buffer.String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateNewRotor() {
	fmt.Println("Creating new Rotor...")

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rotor Rotor
	rotor.Notch = r1.Intn(RotorLength)
	for i := range rotor.Slots {
		rotor.Slots[i] = -1
	}
	current := 0

	for current < RotorLength {
		var n int
		for rotor.Slots[current] == -1 {
			n = r1.Intn(RotorLength)
			found := false
			for i := 0; i <= current; i++ {
				if n == rotor.Slots[i] {
					found = true
					break
				}
			}
			if !found {
				rotor.Slots[current] = n
				current++
				break
			}
		}
	}

	for n := range rotor.Slots {
		fmt.Print(rotor.Slots[n], " ")
	}
	fmt.Println()
	AppendRotor(rotor, false)
}

func CreateNewReturnRotor() {
	fmt.Println("Creating new ReturnRotor...")

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rotor Rotor
	rotor.Notch = r1.Intn(RotorLength)
	for i := range rotor.Slots {
		rotor.Slots[i] = -1
	}
	current := 0

	for current < RotorLength/2 {
		n1 := r1.Intn(RotorLength)
		if rotor.Slots[n1] == -1 {
			n2 := r1.Intn(RotorLength)
			if n2 != n1 && rotor.Slots[n2] == -1 {
				rotor.Slots[n1] = n2
				rotor.Slots[n2] = n1
				current++
			}
		}
	}

	for n := range rotor.Slots {
		fmt.Print(rotor.Slots[n], " ")
	}
	fmt.Println()
	AppendRotor(rotor, true)
}

func AppendRotor(rotor Rotor, isReturnRotor bool) {
	rotorStore := LoadRotorStore()
	if isReturnRotor {
		newid := fmt.Sprintf("E%v", len(rotorStore.ReturnRotors))
		rotorStore.ReturnRotors[newid] = rotor
	} else {
		newid := fmt.Sprintf("R%v", len(rotorStore.Rotors))
		rotorStore.Rotors[newid] = rotor
	}

	var outbuffer bytes.Buffer
	encoder := gob.NewEncoder(&outbuffer)

	err := encoder.Encode(rotorStore)
	check(err)

	f, err := os.Create("c:\\temp\\rotorStore.dat")
	check(err)
	defer f.Close()

	_, err = f.Write(outbuffer.Bytes())
	check(err)
}

func LoadRotorStore() RotorStore {
	var rotorStore RotorStore
	rotorData, err := os.Open("c:\\temp\\rotorStore.dat")
	if err == nil {
		decoder := gob.NewDecoder(rotorData)
		err = decoder.Decode(&rotorStore)
	} else {
		rotors := make(map[string]Rotor)
		returnRotors := make(map[string]Rotor)
		rotorStore = RotorStore{rotors, returnRotors}
	}
	return rotorStore
}
