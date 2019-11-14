package main

import (
	"fmt"
	"os"
	"time"

	redis "github.com/gomodule/redigo/redis"
)

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

	_, err = c.Do("hset", "books", "name", "golang", "author", "Moon", "pages", "4000")
	errCheck(err)

	v, err := redis.String(c.Do("hget", "books", "name"))
	errCheck(err)
	fmt.Println("book.name:", v)

	v, err = redis.String(c.Do("hget", "books", "author"))
	errCheck(err)
	fmt.Println("book.author:", v)

}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
