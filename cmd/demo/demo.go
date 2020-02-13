package main

import (
	"fmt"
	"log"
	"os"
)

// Так можно делать в тесте, если вы передаёте io.Reader, io.Writer
func main() {
	file, err := os.Open("commands.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var login string
	fmt.Fscan(file, &login)
	fmt.Println(login)
	var password string
	fmt.Fscan(file, &password)
	fmt.Println(password)
}
