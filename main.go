package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pastr",
		Usage: "A simple text file sharing service. Upload text files via stdin or file and get them back via key.",
		Commands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "Get a file by key",
				Action: get,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key",
						Aliases: []string{"k"},
						Usage: "The key of the file to get.",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "out",
						Aliases: []string{"o"},
						Usage: "The output file to save the retrieved file.",
					},
				},
			},
			{
				Name:   "create",
				Usage:  "Create a new file",
				Action: create,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Usage:    "The file to upload.",
					},
				},
			},
		},
		DefaultCommand: "create",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func get(ctx *cli.Context) error {
	key := ctx.String("key")
	outFile := ctx.String("out")

	url := "http://pastr.ritiksahni0203.workers.dev/get/" + url.QueryEscape(key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	var out io.Writer
	if outFile != "" {
		file, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}
		defer file.Close()
		out = file
	} else {
		out = os.Stdout
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	return nil
}

func create(ctx *cli.Context) error {
	url := "http://pastr.ritiksahni0203.workers.dev/create"

	var data []byte
	var err error

	if ctx.String("create") != "" {
		filePath := ctx.String("create")
		data, err = os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}
	} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}
	} else {
		err := fmt.Errorf("no input provided: please specify a file or provide input via stdin")
		cli.ShowAppHelp(ctx)
		return err
	}

	stdinStream := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", url, stdinStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	// TODO: handle response and display user-friendly message.
	// 	- e.g. "File created successfully

	io.Copy(os.Stdout, resp.Body)	
	return nil
}
