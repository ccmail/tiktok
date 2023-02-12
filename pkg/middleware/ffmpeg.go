package middleware

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

// ExampleReadFrameAsJpeg 获取封面
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	//if err != nil {
	//	return buf
	//}
	return buf, err
}
