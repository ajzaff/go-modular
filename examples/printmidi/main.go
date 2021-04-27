package main

import (
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/portmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	drv, err := portmididrv.New()
	must(err)

	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	printInPorts(ins)

	in := ins[1]

	must(in.Open())

	rd := reader.New()
	rd.ListenTo(in)
}

func printPort(port midi.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}
