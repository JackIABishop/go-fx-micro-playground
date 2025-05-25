package logging

import (
	"log"
	"os"
)

var Logger *log.Logger

func Init() {
	Logger = log.New(os.Stdout, "[fx-playground] ", log.LstdFlags|log.Lshortfile)
	Logger.Println("Logger initialised.")
}
