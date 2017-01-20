// solexer.go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type arg struct {
	name    string
	argType string
}
type function struct {
	name       string
	args       []arg
	returnType string
}

func main() {
	b, err := ioutil.ReadFile("/Users/joeb/Workspaces/eth/misc/contracts/Constructor.sol") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	str = stripLineComments(str)
	//TODO - strip block comments
	str = strings.Replace(str, "\n", " ", -1) //strip newlines

	fSlice := createFuncSlice(str)
	for _, result := range fSlice {
		fmt.Printf("Function: %s\nArgs: %s\nReturns: %s\n\n", result.name, result.args, result.returnType)

	}

}

func createFuncSlice(s string) []function {

	funcRx := regexp.MustCompile(`function.*?{`)
	strSlice := funcRx.FindAllString(s, -1)
	fmt.Println(strSlice)

	fnslice := make([]function, len(strSlice))

	nameRx := regexp.MustCompile(` .*?\(`)
	argsRx := regexp.MustCompile(`\(.*?\)`)
	spaceRx := regexp.MustCompile("[^\\s]+")

	for i, f := range strSlice {

		fnslice[i].name = strings.TrimSpace(strings.Replace(nameRx.FindString(f), "(", " ", 1))
		args := strings.Trim(strings.Trim(argsRx.FindString(f), "("), ")")
		argSlice := strings.SplitN(args, ",", -1)

		//grep out arguments
		for _, a := range argSlice {
			splitArg := spaceRx.FindAllString(a, -1)
			if len(splitArg) == 2 {
				var anArg arg
				anArg.argType = splitArg[0]
				anArg.name = splitArg[1]
				fnslice[i].args = append(fnslice[i].args, anArg)
			}
		}

		//find return

	}

	return fnslice
}

func stripLineComments(s string) string {
	var buffer bytes.Buffer
	lines := strings.Split(s, "\n")
	rx := regexp.MustCompile(`//.*`)
	for _, line := range lines {
		buffer.WriteString(rx.ReplaceAllString(line, " "))
	}
	return buffer.String()
}
