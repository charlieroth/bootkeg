package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
    EDITOR = "nvim"
    HOME = "/Users/charlie/github.com/charlieroth/zet"
)

func NewEntryTitle() string {
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

func NewEntry(title string) error {
    return os.Mkdir(fmt.Sprintf("%s/%s", HOME, title), 0755)
}

func NewEntryNote(title string) (string, error) {
    path := fmt.Sprintf("%s/%s/README.md", HOME, title)
    if err := os.WriteFile(path, []byte(""), 0755); err != nil {
        return path, err
    }

    return path, nil
}

func NewZet() error {
    title := NewEntryTitle()
    if err := NewEntry(title); err != nil {
        return err
    }

    path, err := NewEntryNote(title)
    if err != nil {
        return err
    }

    cmd := exec.Command(EDITOR, path)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    if err := cmd.Run(); err != nil {
        return err
    }

    return nil
}

func RunCmd(cmd string) error {
	switch cmd {
	case "new":
		return NewZet()
    default:
        return nil
	}
}

func main() {
	args := os.Args[1:]
	cmd := args[0]

    if err := RunCmd(cmd); err != nil {
        log.Fatal(err)
    }
}
