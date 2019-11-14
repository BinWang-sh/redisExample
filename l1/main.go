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

	_, err = c.Do("set", "testkey1", "Hello from redis", "EX", "5")
	errCheck(err)

	r, err := redis.String(c.Do("get", "testkey1"))
	errCheck(err)

	fmt.Println("Get within expire time:", r)

	time.Sleep(8 * time.Second)

	//Check if can get value, after expire time
	r, err = redis.String(c.Do("get", "testkey1"))
	if err != nil {
		fmt.Println(err)
	}

	_, err = c.Do("mset", "name", "Michael", "sex", "M", "age", 23, "postcode", 2343253)
	errCheck(err)

	stringValues, err := redis.Strings(c.Do("mget", "name", "sex"))
	errCheck(err)

	intValues, err := redis.Ints(c.Do("mget", "age", "postcode"))
	errCheck(err)

	for _, v := range stringValues {
		fmt.Println(v)
	}

	for _, i := range intValues {
		fmt.Println(i)
	}

}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
