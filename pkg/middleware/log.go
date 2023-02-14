package middleware

import "log"

func InitLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[tiktok]项目中发生的错误")
}
