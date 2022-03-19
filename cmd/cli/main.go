package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/CristianCurteanu/url-shortener/pkg/mappings/client"
	"github.com/urfave/cli/v2"
)

func main() {
	cl := client.NewClient(os.Getenv("API_HOST"))

	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "create-mapping",
				Usage: "Creates a new url key mapping",
				Action: func(c *cli.Context) error {
					mapping, err := cl.CreateMapping(client.CreateMappingRequest{
						URL: c.String("url"),
					})
					if err != nil {
						return err
					}

					prettyPrint(map[string]string{
						"Key": mapping.Key,
						"URL": c.String("URL"),
					})

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Usage:    "URL that will be mapped to a key",
						Required: true,
					},
				},
			},
			{
				Name:  "get-mapping",
				Usage: "Fetches url key mapping",
				Action: func(c *cli.Context) error {
					key := c.String("key")

					mapping, err := cl.GetMapping(client.MappingRequest{Key: key})
					if err != nil {
						return err
					}

					counter, err := cl.GetMappingCounter(client.MappingRequest{Key: key})
					if err != nil {
						return err
					}

					prettyPrint(map[string]string{
						"Key":     mapping.Key,
						"URL":     mapping.URL,
						"Counter": strconv.Itoa(int(counter.Counter)),
					})

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Usage:    "Key of the stored URL",
						Required: true,
					},
				},
			},
			{
				Name:  "delete-mapping",
				Usage: "Deletes url key mapping",
				Action: func(c *cli.Context) error {
					key := c.String("key")

					resp, err := cl.DeleteMapping(client.DeleteMappingRequest{Key: key})
					if err != nil {
						return err
					}

					prettyPrint(map[string]string{
						"Deleted": resp.Deleted,
					})

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Usage:    "Key of the stored URL",
						Required: true,
					},
				},
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(m map[string]string) {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	for k, v := range m {
		fmt.Println(k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + v)
	}
}
