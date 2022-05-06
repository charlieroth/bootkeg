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
)

func NewEntryTitle() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	min := now.Minute()
	sec := now.Second()
	return fmt.Sprintf("%d%d%d%d%d%d", year, month, day, hour, min, sec)
}

func NewEntry(title string) error {
    return os.Mkdir(title, 0755)
}

func NewEntryNote(title string) (string, error) {
    path := fmt.Sprintf("%s/README.md", title)
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
