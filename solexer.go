// solexer.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

type arg struct {
	Name    string `json:"name"`
	ArgType string `json:"type"`
}
type function struct {
	Name       string `json:"name"`
	Args       []arg  `json:"args"`
	ReturnType string `json:"returnType"`
}

const (
	version = "0.0.1"
)

var (
	iow io.Writer
)

func main() {
	iow = io.Writer(os.Stdout)
	fmt.Fprintln(iow, "Testing testing")
	app := cli.NewApp()
	app.Name = "solexer"
	app.Usage = "Lexical Analyzer for solidity files"
	app.Version = version
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
				fmt.Fprintln(iow, "Find out more about this project at github.com/joeb000/solexer")
				return nil
			},
		},
	}

	app.Run(os.Args)

}

func solexerEntry(c *cli.Context) {

	if c.String("out") != "" {
		file, err := os.OpenFile(c.String("out"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Error opening file: %v", err)
		}
		defer file.Close()
		iow = io.Writer(file)
	}

	b, err := ioutil.ReadFile(c.String("input-file"))
	if err != nil {
		fmt.Fprint(iow, err)
	}
	str := string(b) // convert content to a 'string'
	str = stripLineComments(str)
	//TODO - strip block comments
	str = strings.Replace(str, "\n", " ", -1) //strip newlines

	fSlice := createFuncSlice(str)
	jsonOut, err := json.Marshal(fSlice)
	if err != nil {
		fmt.Errorf("Error marshalling json: %v", err)
	}
	fmt.Fprintln(iow, string(jsonOut))
	for _, result := range fSlice {
		fmt.Fprintf(iow, "Function: %s\nArgs: %s\nReturns: %s\n\n", result.Name, result.Args, result.ReturnType)

	}

}

func createFuncSlice(s string) []function {

	funcRx := regexp.MustCompile(`function.*?{`)
	strSlice := funcRx.FindAllString(s, -1)
	fmt.Fprintln(iow, strSlice)

	fnslice := make([]function, len(strSlice))

	nameRx := regexp.MustCompile(` .*?\(`)
	argsRx := regexp.MustCompile(`\(.*?\)`)
	spaceRx := regexp.MustCompile("[^\\s]+")
	retRx := regexp.MustCompile("returns .*?\\)")

	for i, f := range strSlice {

		fnslice[i].Name = strings.TrimSpace(strings.Replace(nameRx.FindString(f), "(", " ", 1))
		args := removeSpaceAndParens(argsRx.FindString(f))
		argSlice := strings.SplitN(args, ",", -1)

		//grep out arguments
		for _, a := range argSlice {
			splitArg := spaceRx.FindAllString(a, -1)
			if len(splitArg) == 2 {
				var anArg arg
				anArg.ArgType = splitArg[0]
				anArg.Name = splitArg[1]
				fnslice[i].Args = append(fnslice[i].Args, anArg)
			}
		}

		//find return
		if strings.Contains(f, "returns") {
			fnslice[i].ReturnType = removeSpaceAndParens(strings.Trim(retRx.FindString(f), "returns"))
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
