package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"
)

var urlMap = make(map[string]string)
var count int = 1
var dict map[int]byte = map[int]byte{
	0:  '0',
	1:  '1',
	2:  '2',
	3:  '3',
	4:  '4',
	5:  '5',
	6:  '6',
	7:  '7',
	8:  '8',
	9:  '9',
	10: 'a',
	11: 'b',
	12: 'c',
	13: 'd',
	14: 'e',
	15: 'f',
	16: 'g',
	17: 'h',
	18: 'i',
	19: 'j',
	20: 'k',
	21: 'l',
	22: 'm',
	23: 'n',
	24: 'o',
	25: 'p',
	26: 'q',
	27: 'r',
	28: 's',
	29: 't',
	30: 'u',
	31: 'v',
	32: 'w',
	33: 'x',
	34: 'y',
	35: 'z',
	36: 'A',
	37: 'B',
	38: 'C',
	39: 'D',
	40: 'E',
	41: 'F',
	42: 'G',
	43: 'H',
	44: 'I',
	45: 'J',
	46: 'K',
	47: 'L',
	48: 'M',
	49: 'N',
	50: 'O',
	51: 'P',
	52: 'Q',
	53: 'R',
	54: 'S',
	55: 'T',
	56: 'U',
	57: 'V',
	58: 'W',
	59: 'X',
	60: 'Y',
	61: 'Z'}

func main() {
	router := gin.Default()
	router.GET("/shortUrl/:shortUrl", transHandler)
	router.GET("/getUrl/:getUrl", getHandler)
	router.Run(":8081")

}

func transHandler(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	transform(shortUrl)
	fmt.Println(urlMap)
}

func getHandler(c *gin.Context) {
	getUrl := c.Param("getUrl")
	for key, value := range urlMap {
		fmt.Println(value)
		if value == getUrl {
			c.Redirect(http.StatusMovedPermanently, "https://"+key)
		}
	}
	c.String(http.StatusOK, "无该短网址")
}

func transform(s string) string {
	if _, ok := urlMap[s]; ok {
		return urlMap[s]
	} else {
		urlMap[s] = from10to62(count)
		count++
		return urlMap[s]
	}
}

func from10to62(num int) string {
	var str62 []byte
	for {
		var result byte
		var tmp []byte

		number := num % 62
		result = dict[number]

		tmp = append(tmp, result)

		str62 = append(tmp, str62...)
		num = num / 62

		if num == 0 {
			break
		}
	}
	return string(str62)
}
