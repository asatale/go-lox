package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/asatale/go-lox/interpreter"
	"os"
)

func Prompt() {
	for {
		fmt.Printf("Glox Shell>>> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		source := bytes.NewBufferString(text)
		err := interpreter.Run(source)
		if err != nil && err != interpreter.EOFError {
			fmt.Println("Error in interpreter: ", err)
			return
		}
	}
}
