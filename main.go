package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Command string

const (
	GET  Command = "get"
	SET  Command = "set"
	FILE string  = "./.db.txt"
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

func get(args []string, db *os.File) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Must provide a single <key> value.")
	}
	recordsBytes, err := io.ReadAll(db)
	if err != nil {
		return "", err
	}
	recordsString := strings.TrimSpace(string(recordsBytes))
	lines := strings.Split(string(recordsString), "\n")
	dbMap := make(map[string]string)
	for _, l := range lines {
		el := strings.Split(l, ": ")
		dbMap[el[0]] = el[1]
	}

	result, test := dbMap[args[0]]

	if !test {
		err = errors.New("Couldn't find key provided.")
	}

	return result, err
}

func main() {
	file, err := os.OpenFile(FILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Couldn't read database")
	}
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		os.Exit(0)
	}
	command := args[0]
	commandArgs := args[1:]
	switch Command(command) {
	case GET:
		result, err := get(commandArgs, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(result)
	case SET:
		err := set(commandArgs, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	default:
		fmt.Println("Command not recognised.")
		os.Exit(1)
	}
}
