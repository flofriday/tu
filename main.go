package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var rootPath = "/Users/flo/Library/Mobile Documents/com~apple~CloudDocs/TU/"

func main() {
	app := &cli.App{
		Name:  "tu",
		Usage: "Open the folder of the subject",
		Action: func(c *cli.Context) error {
			subject := c.Args().First()
			if subject == "" {
				return openFile(rootPath)
			}

			dir, err := getSubjectPath(subject)
			if err != nil {
				printNotSubject(subject)
				return nil
			}

			return openFile(dir)
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all subjects",
				Action: func(c *cli.Context) error {
					subjects := listSubjects()

					fmt.Printf("Found %d subjects:\n%s\n", len(subjects), strings.Join(subjects, "\n"))
					return nil
				},
			},
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create a new subject",
				Action: func(c *cli.Context) error {
					fmt.Println("Not yet implemented")
					return nil
				},
			},
			{
				Name:    "vowi",
				Aliases: []string{"v"},
				Usage:   "Open the vowi page of the subject",
				Action: func(c *cli.Context) error {
					subject := c.Args().First()
					dir, err := getSubjectPath(subject)
					if err != nil {
						printNotSubject(subject)
						return nil
					}

					file := "vowi.url"
					file = path.Join(dir, file)
					if !fileExists(file) {
						fmt.Printf("Even tough the subject %s exist, there is no vowi.url file in the folder.\n", subject)
						return nil
					}

					return openFile(file)
				},
			},
			{
				Name:  "tuwel",
				Usage: "Open the tuwel page of the subject",
				Action: func(c *cli.Context) error {
					subject := c.Args().First()
					dir, err := getSubjectPath(subject)
					if err != nil {
						printNotSubject(subject)
						return nil
					}

					file := "tuwel.url"
					file = path.Join(dir, file)
					if !fileExists(file) {
						fmt.Printf("Even tough the subject %s exist, there is no vowi.url file in the folder.\n", subject)
						return nil
					}

					return openFile(file)
				},
			},
			{
				Name:  "tiss",
				Usage: "Open the tiss page of the subject",
				Action: func(c *cli.Context) error {
					subject := c.Args().First()
					dir, err := getSubjectPath(subject)
					if err != nil {
						printNotSubject(subject)
						return nil
					}

					file := "tiss.url"
					file = path.Join(dir, file)
					if !fileExists(file) {
						fmt.Printf("Even tough the subject %s exist, there is no vowi.url file in the folder.\n", subject)
						return nil
					}

					return openFile(file)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Open a file or directory with the default application
// Note: Only testsed on macOS for now
func openFile(path string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin": // macOS
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, path)
	return exec.Command(cmd, args...).Start()
}

func listSubjects() []string {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	var subjects []string

	for _, file := range files {
		// Ignore files that do not math the subject
		if !file.IsDir() {
			continue
		}

		// Add the folder to the list
		subjects = append(subjects, file.Name())
	}

	// We didn't find the subject
	return subjects
}

func getSubjectPath(subject string) (string, error) {
	subject = strings.ToLower(subject)
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// Ignore files that do not math the subject
		if !file.IsDir() || strings.ToLower(file.Name()) != subject {
			continue
		}

		// We found the subject so we can return the complete path
		return path.Join(rootPath, file.Name()), nil
	}

	// We didn't find the subject
	return "", errors.New("could not find the subject")
}

func printNotSubject(subject string) {
	if subject == "" {
		fmt.Println("You need to provide a subject.\nType `tu help` to see all options.")
		return
	}

	fmt.Printf("The subject %s does not exist.\n", subject)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
