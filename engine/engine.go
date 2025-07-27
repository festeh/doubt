package engine

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

type Engine struct {
	executablePath string
	commands       []Command
	outputDir      string
}

func NewEngine(executablePath string) *Engine {
	return &Engine{
		executablePath: executablePath,
		commands:       make([]Command, 0),
	}
}

func (e *Engine) AddCommand(cmd Command) {
	e.commands = append(e.commands, cmd)
}

func (e *Engine) SetOutputDir(dir string) {
	e.outputDir = dir
}

func (e *Engine) Run() error {
	cmd := exec.Command(e.executablePath)
	
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("error starting %s: %w", e.executablePath, err)
	}
	defer ptmx.Close()

	var outputFile *os.File
	var outputWriter io.Writer = os.Stdout

	if e.outputDir != "" {
		if err := os.MkdirAll(e.outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}

		outputPath := filepath.Join(e.outputDir, "output.txt")
		outputFile, err = os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer outputFile.Close()

		outputWriter = io.MultiWriter(os.Stdout, outputFile)
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("error setting raw mode: %w", err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Fprintf(os.Stderr, "Error resizing pty: %v\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	go func() {
		io.Copy(ptmx, os.Stdin)
	}()

	go func() {
		for _, command := range e.commands {
			if err := command.Execute(ptmx); err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			}
		}
	}()
	
	io.Copy(outputWriter, ptmx)
	
	return cmd.Wait()
}