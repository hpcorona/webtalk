package main

import (
	"strings"
	"os"
	"fmt"
)

type Frame struct {
	command string
	data string
}

func NewFrame(msg string) (*Frame, os.Error) {
	if strings.HasPrefix(msg, "[") == false {
		return nil, os.NewError("[ brace not found")
	}
	
	idx := strings.Index(msg, "]")
	if idx < 0 {
		return nil, os.NewError("] brace not found")
	}
	
	cmd := msg[1:idx]
	
	mess := ""
	
	if idx + 2 < len(msg) {
		mess = msg[idx+2:]
	}
	
	return &Frame { command: cmd , data: mess }, nil
}

func NewFrameString(cmd, data string) string {
	total := len(cmd) + len(data) + 3
	
	allCmd := fmt.Sprintf("%05d", total) + "[" + cmd + "] " + data
	
	return allCmd
}
