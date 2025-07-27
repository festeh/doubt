package main

import (
	"fmt"
	"os"

	"github.com/festeh/doubt/engine"
)

func main() {
	eng := engine.NewEngine("/home/dima/projects/bro/bro")
	
	eng.AddCommand(engine.NewSleepCommand(100))
	eng.AddCommand(engine.NewKeypressCommand("ctrl-c"))
	
	if err := eng.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}