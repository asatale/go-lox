package cmd

import (
  "fmt"
  "github.com/asatale/go-lox/interpreter"
  "os"
)

func Script(filename string) {
  fd, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fd.Close()

  err = interpreter.Run(fd)
  if err != nil {
    fmt.Println("Error in interpreter: ", err)
    return
  }
}
