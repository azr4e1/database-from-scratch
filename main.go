package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Command string

const (
	GET    Command = "get"
	SET    Command = "set"
	DELETE Command = "delete"
	FILE   string  = "./.db.txt"
)

func help() {
	fmt.Println("This is a help function")
}

func set(args []string, db *os.File) error {
	if len(args) != 2 {
		return errors.New("Must provide a <key,value> pair.")
	}
	key, value := args[0], args[1]

	record := fmt.Sprintf("%s: %s\n", key, value)
	_, err := db.Write([]byte(record))
	return err
}

func get(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("must provide a single <key> value")
	}

	f, err := os.Open(FILE)
	if err != nil {
		return "", err
	}
	defer f.Close()

	prefix := args[0] + ": "
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("key not found")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		os.Exit(0)
	}
	command := args[0]
	commandArgs := args[1:]
	switch Command(command) {
	case GET:
		result, err := get(commandArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(result)
	case SET:
		err := set(commandArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	default:
		fmt.Println("Command not recognised.")
		os.Exit(1)
	}
}
