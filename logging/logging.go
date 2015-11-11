package logging

import (
	"log"
)

func Error(err error) {
	log.Printf("[ERROR] %s\n", err)
}

func Info(str string) {
	log.Printf("[INFO] %s\n", str)
}
