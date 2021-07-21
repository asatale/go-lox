package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/asatale/go-lox/interpreter"
	"io"
	"os"
)

func Prompt() {

	sigs := make(chan os.Signal, 1)
	data := make(chan string, 1)
	control := make(chan struct{}, 1)
	control <- struct{}{}

	go func() {
	Loop1:
		for {
			select {
			case _, ok := <-control:
				if !ok {
					break Loop1
				}
				fmt.Printf("Glox Shell>>> ")
				reader := bufio.NewReader(os.Stdin)
				text, err := reader.ReadString('\n')
				if err == io.EOF {
					close(data)
					break Loop1
				}
				data <- text
			}
		}
	}()

Loop2:
	for {
		select {
		case <-sigs:
			fmt.Println("Got shutdown, exiting")
			break Loop2
		case s, ok := <-data:
			if !ok {
				break Loop2
			}
			source := bytes.NewBufferString(s)
			err := interpreter.Run(source)
			if err != nil && err != interpreter.EOFError {
				fmt.Println("Error in interpreter: ", err)
				break Loop2
			}
			control <- struct{}{}
		}
	}
	close(control)
}
