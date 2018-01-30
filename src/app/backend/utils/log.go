package utils

import (
	"fmt"
	"log"
	//"time"
)

// log info
func LogI(format string, v ...interface{}) {
	log.Printf("[INFO]"+format, v...)
}

// log debug
func LogD(format string, v ...interface{}) {
	//log.Printf(fmt.Sprintf("[%s][DEBUG]%s", time.Now().Format(time.RFC3339), format), v...)
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
