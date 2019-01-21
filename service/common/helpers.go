package common

import (
	"fmt"
	"log"
	"os"
)

func CreateLogger(name string) *log.Logger {
	prefix := fmt.Sprintf("[%s] ", name)
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
