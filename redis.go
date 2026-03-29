package main

import (
	"fmt"
	"strings"

	store "github.com/yaash45/redis/internal/store"
)

func main() {
	fmt.Println("Redis clone starting...")
	fmt.Printf("Ready. Enter commands in format [SET/GET/DEL] [KEY] [VALUE]. Enter 'exit' to stop:\n\n")

	kvs := store.NewStore()

	for {
		var cmd string
		var key string
		var value string

		fmt.Printf("> ")
		fmt.Scanln(&cmd, &key, &value)

		if cmd == "GET" {

			lookup_result, err := kvs.Get(key)

			if err != nil {
				fmt.Printf("key '%s' does not exist", key)
			} else {
				fmt.Printf("%s", lookup_result)
			}

		} else if cmd == "SET" {

			if key == "" || value == "" {
				fmt.Printf("[error] 'SET' command requires non-empty key and value.")

			} else {
				kvs.Set(key, value)
				fmt.Printf("[ok]")
			}

		} else if cmd == "DEL" {
			kvs.Delete(key)
			fmt.Printf("[ok]")

		} else if strings.ToLower(cmd) == "exit" {
			break

		} else {
			fmt.Printf("Unrecognized command, try again.")

		}

		fmt.Printf("\n\n")
	}
}
