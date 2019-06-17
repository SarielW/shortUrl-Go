package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"

	"github.com/go-redis/redis"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var count int
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

type Urlmap struct {
	gorm.Model
	LongUrl string
	ShortUrl string
}

func main() {
	Init()
	router := gin.Default()
	router.POST("/originUrl", transHandler)
	router.GET("/getUrl/:shortUrl", getHandler)
	router.Run(":8081")

}

func Init() {
	db, err := gorm.Open("mysql", "root:admin@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var urlmap Urlmap
	db.HasTable(&Urlmap{})

	db.CreateTable(&Urlmap{})
	db.Last(&urlmap).GetErrors()
	id := urlmap.ID + 1
	count = int(id)
}

func transHandler(c *gin.Context) {
	origin := c.PostForm("origin")
	getUrl(origin)
}

func getHandler(c *gin.Context) {
	url := c.Param("shortUrl")
	c.Redirect(http.StatusMovedPermanently, findUrl(url))
}

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:	"localhost:6379",
		Password: "",
		DB: 0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
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

// redis存 short-url：long-url
func getUrl(s string) string{
	db, err := gorm.Open("mysql", "root:admin@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var urlmap Urlmap
	if db.First(&urlmap, "long_url = ?", s).RecordNotFound() {
		client := createClient()
		err := client.Set( from10to62(count), s, 0).Err()
		if err != nil {
			panic(err)
		}
		urlmap := Urlmap{LongUrl: s, ShortUrl: from10to62(count)}
		count++
		db.Create(&urlmap)

	}
	db.First(&urlmap, "long_url = ?", s)
	fmt.Println(urlmap.ShortUrl)
	return  urlmap.ShortUrl
}

func findUrl(s string) string {
	client := createClient()
	val, err := client.Get(s).Result()
	if err == nil {
		fmt.Println("using redis find:" + val)
		return val
	}
	db, err := gorm.Open("mysql", "root:admin@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var urlmap Urlmap
	if db.First(&urlmap, "short_url = ?", s).RecordNotFound() {
		fmt.Println("无此网址")
	}
	db.First(&urlmap, "short_url = ?", s)
	return  urlmap.LongUrl
}
