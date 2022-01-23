package main

import (
	"fmt"
)

func main() {
	country := map[string]string{"HELLO": "ENGLISH", "HOLA": "SPANISH", "HALLO": "GERMAN", "BONJOUR": "FRENCH", "CIAO": "ITALIAN", "ZDRAVSTVUJTE": "RUSSIAN"}
	msg := ""
	count := 0
	fmt.Scan(&msg)
	for msg != "#" {
		count++
		if lang, ok := country[msg]; ok {
			fmt.Printf("Case %d: %s\n", count, lang)
		} else {
			fmt.Printf("Case %d: UNKNOWN\n", count)
		}
		fmt.Scan(&msg)
	}
}
