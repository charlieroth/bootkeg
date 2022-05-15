package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Editor    string
    RootDir   string
	ZetDir    string
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

func CreateContainer(location, isosec string) error {
	return os.Mkdir(fmt.Sprintf("%s/%s", location, isosec), 0755)
}

func CreateFile(location, isosec string) (string, error) {
	path := fmt.Sprintf("%s/%s/README.md", location, isosec)
	if err := os.WriteFile(path, []byte(""), 0755); err != nil {
		return path, err
	}

	return path, nil
}

func New(location, contentType string) error {
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

    if err := os.Chdir(config.RootDir); err != nil {
        return err
    }

    fmt.Print("Commit [y/n]?: ")
    var shouldCommit string
    fmt.Scanln(&shouldCommit)

    if shouldCommit == "n" || shouldCommit == "N" {
        return nil
    }

    if shouldCommit != "y" && shouldCommit != "Y" {
        return nil
    }

    cmd = exec.Command("git", "add", ".")
    if err := cmd.Run(); err != nil {
		return err
    }
	
    cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("[new]: %s", isosec))
    if err := cmd.Run(); err != nil {
		return err
    }
    
    out, err := exec.Command("git", "push").Output()
    if err != nil {
		return err
    }
    fmt.Println(string(out))

	return nil
}

func NewCmd(contentType string) error {
	switch contentType {
	case "zet":
		return New(config.ZetDir, contentType)
	case "post":
		return New(config.PostDir, contentType)
	case "note":
		return New(config.NoteDir, contentType)
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

func LoadConfig() error {
    config = Config{
        Editor: os.Getenv("EDITOR"),
        RootDir: os.Getenv("KEG_ROOT"),
        ZetDir: os.Getenv("KEG_ZET"),
        PostDir: os.Getenv("KEG_POST"),
        NoteDir: os.Getenv("KEG_NOTE"),
    }

    if config.Editor == "" {
        return fmt.Errorf("no EDITOR environment variable set")
    }
    
    if config.RootDir == "" {
        return fmt.Errorf("no KEG_ROOT environment variable set")
    }
    
    if config.ZetDir == "" {
        return fmt.Errorf("no KEG_ZET environment variable set")
    }
    
    if config.NoteDir == "" {
        return fmt.Errorf("no KEG_NOTE environment variable set")
    }
    
    if config.PostDir == "" {
        return fmt.Errorf("no KEG_POST environment variable set")
    }

    return nil
}

func main() {
    if err := LoadConfig(); err != nil {
		log.Fatal(err)
    }

	args := os.Args[1:]
	cmd := args[0]
	contentType := args[1]

	if err := Cmd(cmd, contentType); err != nil {
		log.Fatal(err)
	}
}
