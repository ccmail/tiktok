package util

import "tiktok/config"

// RemoveIllegalChar 去除文件名字中的非法字符, 以免存储到oss时报错
func RemoveIllegalChar(str string) string {
	mp := map[byte]interface{}{}
	for i := range config.IllegalChar {
		mp[config.IllegalChar[i]] = struct{}{}
	}
	res := make([]byte, 0, len(str)>>2)
	for i := range str {
		if _, ok := mp[str[i]]; ok {
			continue
		}
		res = append(res, str[i])
	}
	return string(res)
}
