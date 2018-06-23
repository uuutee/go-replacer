package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
	"fmt"
	"path/filepath"
	"bufio"
	"strings"
	"io/ioutil"
)

func main() {
	app := cli.NewApp()

	app.Name = "replacer"
	app.Usage = "search and replace text file."
	app.Version = "0.0.1"

	// オプションを登録する
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "s, search",
			Usage: "search text",
		},
		cli.StringFlag{
			Name: "r, replace",
			Usage: "replace text",
		},
		cli.BoolFlag{
			Name: "dry-run",
			Usage: "use dry-run mode",
		},
	}

	// 実行内容
	app.Action = doAction

	// 実行
	app.Run(os.Args)
}

func doAction(context *cli.Context) error {
	root := getRoot(context)
	searchText := context.String("search")
	replaceText := context.String("replace")
	isDryRun := context.Bool("dry-run")

	switch {
	case searchText == "":
		cli.ShowCommandHelp(context, "search")
		os.Exit(1)
	case replaceText == "":
		cli.ShowCommandHelp(context, "replace")
		os.Exit(1)
	case isDryRun:
		targets := getTargets(root, searchText)
		for _, path := range targets {
			fmt.Println(path)
		}
		fmt.Println("end dry-run mode.")
		os.Exit(1)
	default:
		targets := getTargets(root, searchText)
		for _, path := range targets {
			err := replaceFile(path, searchText, replaceText)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
		}
	}

	return nil
}

func getRoot(context *cli.Context) string {
	root := context.Args().Get(0)
	if root != "" {
		return root
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return wd
}

func getTargets(root, query string) []string {
	var targets []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if contains(path, query) {
				targets = append(targets, path)
			}
		}
		return nil
	})

	return targets
}

func contains(path string, query string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), query) {
			return true
		}
	}

	return false
}

func replaceFile(path, search, replace string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	replaced := strings.Replace(string(data), search, replace, -1)

	return ioutil.WriteFile(path, []byte(replaced), 0666)
}


