package main

import (
	"fmt"
	"os"
	"time"

	redis "github.com/gomodule/redigo/redis"
)

const (
	MAXIDLE       = 50
	MAXACTIVE     = 5000
	IDLETIMEOUT   = 30 * time.Second
	ROUNTINECOUNT = 50
)

func deferClose(con *redis.Conn) {
	fmt.Println("close")
	(*con).Close()
}

func main() {

	redisPool := &redis.Pool{
		MaxIdle:     MAXIDLE,
		MaxActive:   MAXACTIVE,
		IdleTimeout: IDLETIMEOUT,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				"172.17.84.205:6379",
				redis.DialKeepAlive(20*time.Second),
				redis.DialPassword("123456"),
				redis.DialConnectTimeout(15*time.Second),
				redis.DialReadTimeout(15*time.Second),
				redis.DialWriteTimeout(15*time.Second))

			if err != nil {
				fmt.Println(err)
			}
			return c, err
		},
	}

	con := redisPool.Get()

	con.Do("SET", "SPKey1", "specialValue")

	//If SPKey1 does not exist set SPKey1 to updatedValue
	con.Do("SETNX", "SPKey1", "updatedValue")

	v, _ := redis.String(con.Do("GET", "SPKey1"))

	fmt.Println(v)

	//If SPKey1 exist set to newValue
	_, err := con.Do("SETEX", "SPKey1", 5, "newValue")
	errCheck(err)

	v, _ = redis.String(con.Do("GET", "SPKey1"))

	//LPUSH push values to the head of array
	_, err = con.Do("LPUSH", "languages", "python", "java")
	errCheck(err)

	_, err = con.Do("RPUSH", "languages", "golang")
	errCheck(err)

	//LPOP pop up the value at the head of array
	_, err = con.Do("LPOP", "languages")

	//LRANGE get values from array
	values, _ := redis.Values(con.Do("LRANGE", "languages", "0", "100"))
	for _, v := range values {
		fmt.Println(string(v.([]byte)))
	}
}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
