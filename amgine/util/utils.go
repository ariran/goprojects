package util

import (
	"bytes"
	"fmt"
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
	ConfigFile      string
	RotorStore      string
	Command         string
	Target          string
	EncryptFilename bool
	DecryptFilename bool
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
