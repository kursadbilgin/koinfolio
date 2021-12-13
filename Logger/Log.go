package Logger

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
)
