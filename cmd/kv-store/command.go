package main

import (
	"fmt"
	"strings"

	"github.com/ritik6559/kv-store/internal/store"
)

type Command struct {
	Op    string
	Key   string
	Value string // for RENAME this is the new key
}

func dispatch(s store.Store, command Command) error {
	switch strings.ToUpper(command.Op) {
	case "GET":
		value, err := s.Get(command.Key)
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	case "SET":
		return s.Set(command.Key, command.Value)
	case "INCR":
		value, err := s.Incr(command.Key)
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	case "KEYS":
		for _, key := range s.Keys() {
			fmt.Println(key)
		}
		return nil
	case "DELETE":
		if s.Delete(command.Key) {
			fmt.Println(1)
		} else {
			fmt.Println(0)
		}
		return nil
	case "LEN":
		fmt.Println(s.Len())
		return nil
	case "RENAME":
		return s.Rename(command.Key, command.Value)
	case "POP":
		value, err := s.Pop(command.Key)
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	default:
		return fmt.Errorf("unsupported command: %s", command.Op)
	}
}

func runCommands(s store.Store, commands []Command) {
	for _, command := range commands {
		if err := dispatch(s, command); err != nil {
			fmt.Printf("%s %s: %v\n", strings.ToUpper(command.Op), command.Key, err)
		}
	}
}
