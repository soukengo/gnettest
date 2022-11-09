package strings

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var charStr = "zxcvbnmlkjhgfdsaqwertyuiopQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func RandStr(length int) string {
	strArr := make([]string, 0)
	//长度为几就循环几次
	for i := 0; i < length; i++ {
		//产生0-61的数字
		number := rand.Intn(62)
		//将产生的数字通过length次承载到sb中
		strArr = append(strArr, string(charStr[number]))
	}
	//将承载的字符转换成字符串
	return strings.Join(strArr, "")
}

func DecodeUnicode(str string) string {
	sUnicodev := strings.Split(str, "\\u")
	var c string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		c += fmt.Sprintf("%c", temp)
	}
	return c
}
