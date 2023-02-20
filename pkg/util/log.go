package util

import "log"

func InitLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[tiktok]")
}
