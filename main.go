package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Editor string
	ZetDir   string
	NoteDir   string
	PostDir   string
}

var config Config

func CreateIsosec() string {
	now := time.Now().UTC()
	year := now.Year()

	monthInt := int(now.Month())
	var month string
	if monthInt < 10 {
		month = fmt.Sprintf("0%d", monthInt)
	} else {
		month = fmt.Sprintf("%d", monthInt)
	}

	dayInt := now.Day()
	var day string
	if dayInt < 10 {
		day = fmt.Sprintf("0%d", dayInt)
	} else {
		day = fmt.Sprintf("%d", dayInt)
	}

	hourInt := now.Hour()
	var hour string
	if hourInt < 10 {
		hour = fmt.Sprintf("0%d", hourInt)
	} else {
		hour = fmt.Sprintf("%d", hourInt)
	}

	minInt := now.Minute()
	var min string
	if minInt < 10 {
		min = fmt.Sprintf("0%d", minInt)
	} else {
		min = fmt.Sprintf("%d", minInt)
	}

	secInt := now.Second()
	var sec string
	if secInt < 10 {
		sec = fmt.Sprintf("0%d", secInt)
	} else {
		sec = fmt.Sprintf("%d", secInt)
	}

	return fmt.Sprintf("%d%s%s%s%s%s", year, month, day, hour, min, sec)
}

func CreateContainer(isosec, location string) error {
	return os.Mkdir(fmt.Sprintf("%s/%s", location, isosec), 0755)
}

func CreateFile(isosec, location string) (string, error) {
	path := fmt.Sprintf("%s/%s/README.md", location, isosec)
	if err := os.WriteFile(path, []byte(""), 0755); err != nil {
		return path, err
	}

	return path, nil
}

func New(location string) error {
	isosec := CreateIsosec()
	if err := CreateContainer(location, isosec); err != nil {
		return err
	}

	path, err := CreateFile(location, isosec)
	if err != nil {
		return err
	}

	cmd := exec.Command(config.Editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func NewCmd(contentType string) error {
	switch contentType {
	case "zet":
		return New(config.ZetDir)
	case "post":
		return New(config.PostDir)
	case "note":
		return New(config.NoteDir)
	default:
        // TODO(charlieroth): Replace with Usage()
		return fmt.Errorf("content type '%s' is not a valid", contentType)
	}
}

func Cmd(cmd, contentType string) error {
	switch cmd {
	case "new":
		return NewCmd(contentType)
	default:
        // TODO(charlieroth): Replace with Usage()
		return fmt.Errorf("command '%s' is not a valid", cmd)
	}
}

func main() {
    config = Config{
        Editor: os.Getenv("EDITOR"),
        ZetDir: os.Getenv("KEG_ZET"),
        PostDir: os.Getenv("KEG_POST"),
        NoteDir: os.Getenv("KEG_NOTE"),
    }

	args := os.Args[1:]
	cmd := args[0]
	contentType := args[1]

	if err := Cmd(cmd, contentType); err != nil {
		log.Fatal(err)
	}
}
