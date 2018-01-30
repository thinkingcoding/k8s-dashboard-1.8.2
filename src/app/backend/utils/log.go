package utils

import (
	"fmt"
	"log"
)

// log info
func LogI(format string, v ...interface{}) {
	log.Printf("[INFO]"+format, v...)
}

// log debug
func LogD(format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[DEBUG]%s", format), v...)
}

// log warn
func LogW(format string, v ...interface{}) {
	log.Printf("[WARN]"+format, v...)
}

// log error
func LogE(format string, v ...interface{}) {
	log.Printf("[ERROR]"+format, v...)
}
