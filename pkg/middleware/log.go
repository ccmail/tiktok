package middleware

import "log"

func InitLog() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
