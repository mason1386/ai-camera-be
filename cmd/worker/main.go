package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Starting worker...")
	for {
		log.Println("Worker running...")
		time.Sleep(5 * time.Second)
	}
}
