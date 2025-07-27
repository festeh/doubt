package engine

import (
	"io"
	"time"
)

type Command interface {
	Execute(pty io.Writer) error
}

type SleepCommand struct {
	Duration time.Duration
}

func NewSleepCommand(ms int) *SleepCommand {
	return &SleepCommand{
		Duration: time.Duration(ms) * time.Millisecond,
	}
}

func (s *SleepCommand) Execute(pty io.Writer) error {
	time.Sleep(s.Duration)
	return nil
}

type KeypressCommand struct {
	Key string
}

func NewKeypressCommand(key string) *KeypressCommand {
	return &KeypressCommand{
		Key: key,
	}
}

func (k *KeypressCommand) Execute(pty io.Writer) error {
	var keyBytes []byte
	
	switch k.Key {
	case "ctrl-c":
		keyBytes = []byte{0x03}
	case "enter":
		keyBytes = []byte{'\n'}
	case "escape":
		keyBytes = []byte{0x1b}
	case "tab":
		keyBytes = []byte{'\t'}
	case "space":
		keyBytes = []byte{' '}
	default:
		keyBytes = []byte(k.Key)
	}
	
	_, err := pty.Write(keyBytes)
	return err
}