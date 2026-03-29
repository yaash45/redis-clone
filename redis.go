package main

import (
	"fmt"
	"strings"

	store "github.com/yaash45/redis/internal/store"
)

func main() {
	fmt.Println("Redis clone starting...")
	fmt.Printf("Ready. Enter commands:\n\n")

	kvs := store.NewStore()

	for {
		var cmd string
		var key string
		var value string

		fmt.Scanln(&cmd, &key, &value)

		if cmd == "GET" {
			fmt.Printf("> %s\n\n", kvs.Get(key))
		} else if cmd == "SET" {
			kvs.Set(key, value)
			fmt.Printf(">\n\n")
		} else if cmd == "DEL" {
			kvs.Delete(key)
			fmt.Printf(">\n\n")
		} else if strings.ToLower(cmd) == "exit" {
			break
		} else {
			fmt.Printf("Unrecognized command, try again.\n\n")
		}
	}
}
