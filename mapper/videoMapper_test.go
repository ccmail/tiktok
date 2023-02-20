package mapper

import (
	"fmt"
	"testing"
)

func TestGetVideoInfo(t *testing.T) {
	InitDBConnectorSupportTest()
	w, err := GetVideoInfo(VideoIDList(1))
	if err != nil {
		t.Error()
	}
	for _, v := range w {
		fmt.Println(v.ID)
		fmt.Println(v.AuthorID)
		fmt.Println(v.FavoriteCount)
		fmt.Println(v.CoverUrl)
		fmt.Println(v.PlayUrl)
		fmt.Println()
		fmt.Println()
	}
}
