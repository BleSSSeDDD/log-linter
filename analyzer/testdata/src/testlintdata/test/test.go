package todo

import (
	"fmt"
	"log"
)

func TestLogs() {
	log.Println("hello") // want "нашел log функцию"
	log.Println("world") // want "нашел log функцию"

	println("not log")
	fmt.Println("fmt")
}
