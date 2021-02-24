package main

import (
  "github.com/asatale/go-lox/cmd"
  "os"
)

func main() {
  cmdArgs := os.Args[1:]
  switch len(cmdArgs) {
  case 1:
    cmd.Script(cmdArgs[0])
  case 0:
    cmd.Prompt()
  default:
    println("Usage: glox [script]")
  }
}
