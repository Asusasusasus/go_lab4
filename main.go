package main

import (
	"fmt"
	"github.com/drunyaD/go_lab4/engine"
        "github.com/drunyaD/go_lab4/commands"
	"bufio"
	"strings"
	"os"
	"errors"
)


func main() {
	eventLoop := new(engine.EventLoop)
	eventLoop.Start()
	if input, err := os.Open("./commands.txt"); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := parse(commandLine)
			eventLoop.Post(cmd)
		}
	}
	eventLoop.AwaitFinish()
}