package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/festeh/doubt/config"
	"github.com/festeh/doubt/engine"
)

func main() {
	var configPath string
	var outputDir string

	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.StringVar(&outputDir, "o", "", "Output directory for artifacts")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -c <config_file_path> [-o <output_directory>]\n", os.Args[0])
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	eng := engine.NewEngine(cfg.ExecutablePath)

	if outputDir != "" {
		eng.SetOutputDir(outputDir)
	}

	commands, err := cfg.ToEngineCommands()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing commands: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands {
		eng.AddCommand(cmd)
	}

	if err := eng.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}