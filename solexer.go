// solexer.go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
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

	app := cli.NewApp()
	app.Name = "solexer"
	app.Usage = "Lexical Analyzer for solidity files"
	app.Action = func(c *cli.Context) error {
		solexerEntry(c)
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input-file, i",
			Usage: "Perform lex on `FILE`",
		},
		cli.StringFlag{
			Name:  "out, o",
			Usage: "Output to `FILE`, if not specified stdout is used",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "Choose format for output file",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "about",
			Usage: "More info about this project",
			Action: func(c *cli.Context) error {
				fmt.Println("Find out more about this project at github.com/joeb000/solexer")
				return nil
			},
		},
	}

	app.Run(os.Args)

}

func solexerEntry(c *cli.Context) {

	b, err := ioutil.ReadFile(c.String("input-file"))
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
	retRx := regexp.MustCompile("returns .*?\\)")

	for i, f := range strSlice {

		fnslice[i].name = strings.TrimSpace(strings.Replace(nameRx.FindString(f), "(", " ", 1))
		args := removeSpaceAndParens(argsRx.FindString(f))
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
		if strings.Contains(f, "returns") {
			fnslice[i].returnType = removeSpaceAndParens(strings.Trim(retRx.FindString(f), "returns"))
		}
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

func removeSpaceAndParens(s string) string {
	s = strings.Replace(s, "(", " ", -1)
	s = strings.Replace(s, ")", " ", -1)
	s = strings.TrimSpace(s)
	return s
}
