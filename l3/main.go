package main

import (
	"fmt"
	"os"
	"time"

	redis "github.com/gomodule/redigo/redis"
)

type Book struct {
	BookName  string
	Author    string
	PageCount string
	Press     string
}

func main() {
	c, err := redis.Dial("tcp",
		"172.17.84.205:6379",
		redis.DialKeepAlive(1*time.Second),
		redis.DialPassword("123456"),
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second))

	errCheck(err)

	defer c.Close()

	//Struct
	top1 := Book{
		BookName:  "Crazy golang",
		Author:    "Moon",
		PageCount: "600",
		Press:     "GoodBook",
	}

	if _, err = c.Do("HMSET", redis.Args{}.Add("Top1").AddFlat(&top1)...); err != nil {
		fmt.Println(err)
		return
	}

	//Map
	top2 := map[string]string{
		"BookName":  "Mast C++",
		"Author":    "Diablo",
		"PageCount": "2000",
		"Press":     "BLZ",
	}

	if _, err = c.Do("HMSET", redis.Args{}.Add("Top2").AddFlat(top2)...); err != nil {
		fmt.Println(err)
		return
	}

	top3 := []string{"BookName", "Deep learning",
		"Author", "Barl",
		"PageCount", "2600",
		"Publish House", "BLZ"}
	_, err = c.Do("HMSET", redis.Args{}.Add("Top3").AddFlat(top3)...)
	errCheck(err)

	topx := Book{}

	for _, item := range []string{"Top1", "Top2", "Top3"} {
		value, err := redis.Values(c.Do("HGETALL", item))
		errCheck(err)

		err = redis.ScanStruct(value, &topx)
		errCheck(err)

		fmt.Printf("%s[%+v]\n", item, topx)
	}

	stringsValue, err := redis.Strings(c.Do("HMGET", "Top1", "BookName", "Author"))
	errCheck(err)

	fmt.Printf("hmget:%+v\n", stringsValue)

}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
