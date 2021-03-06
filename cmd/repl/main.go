package main

import (
	"bufio"
	"fmt"
	"github.com/gypsydave5/gobel/pkg/gobel"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if isPipe(os.Stdin) {
		c, _ := ioutil.ReadAll(os.Stdin)
		result := gobel.Eval(gobel.Read(string(c)), gobel.GlobalEnv())
		fmt.Println(result)
	} else {
		repl()
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := gobel.GlobalEnv()

	for {
		fmt.Print("> ")
		expression, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		ts := gobel.Read(expression)
		result := gobel.Eval(ts, env)
		fmt.Println(result)
	}

	fmt.Println("\nHave a nice day!")
}

func isPipe(f *os.File) bool {
	fi, _ := f.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
