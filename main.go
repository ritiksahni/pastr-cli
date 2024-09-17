package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
)

func main()  {
	app := &cli.App{
		Name: "pastr",
		Usage: "A simple text file sharing service.",
		Action: general,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "key",
				Usage: "The key of the file to get.",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func general(ctx *cli.Context) error {
	if ctx.String("key") != "" {
		return get(ctx);
	} 

	create();

	return nil;
}

func get(c *cli.Context) error {
	key := c.String("key");

	url := "http://localhost:8787/get/" + url.QueryEscape(key);
	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err);
		return err;
	}

	client := &http.Client{};
	resp, err := client.Do(req);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err);
		return err;
	}

	defer resp.Body.Close();
	io.Copy(os.Stdout, resp.Body);
	
	return nil;
}

func create() error {
	url := "http://localhost:8787/create";

	data := os.Stdin;

	req, err := http.NewRequest("POST", url, data);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err);
		return err;
	}

	client := &http.Client{};
	resp, err := client.Do(req);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err);
		return err;
	}

	defer resp.Body.Close();
	io.Copy(os.Stdout, resp.Body);
	
	return nil;
}