package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ritik6559/kv-store/internal/store"
)

type Command struct {
	Op    string
	Key   string
	Value string
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
		s.Delete(command.Key)
		return nil
	case "LEN":
		fmt.Println(s.Len())
		return nil
	default:
		return fmt.Errorf("unsupported command: %s", command.Op)
	}
}

func runCommand(s store.Store, commands []Command) {
	var wg sync.WaitGroup
	errs := make(chan error, len(commands))

	for _, command := range commands {
		wg.Go(func() {
			errs <- dispatch(s, command)
		})
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			fmt.Println(err)
		}
	}
}
