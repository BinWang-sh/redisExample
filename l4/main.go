package main

import (
	"fmt"
	"os"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(2 * ROUNTINECOUNT)

	for i := 0; i < ROUNTINECOUNT; i++ {
		go func(routineNum int) {
			for cnt := 0; cnt < 1000; cnt++ {
				c := redisPool.Get()
				//defer c.Close()

				key := fmt.Sprintf("key_%d_%d", routineNum, cnt)
				value := fmt.Sprintf("value_%d_%d", routineNum, cnt)

				_, err := c.Do("set", key, value)
				if err != nil {
					fmt.Printf("set %s:%v\n", key, err)
				}
				fmt.Printf("s %s\n", value)

				if cnt%50 == 0 {
					aCount := redisPool.Stats().ActiveCount
					wCount := redisPool.Stats().WaitCount
					fmt.Printf("activeCount:%d, waitCount:%d\n", aCount, wCount)
				}

				c.Close()
				//time.Sleep(50 * time.Millisecond)
			}

			wg.Done()
		}(i)

		go func(routineNum int) {
			for cnt := 0; cnt < 1000; cnt++ {
				c := redisPool.Get()
				//defer c.Close()
				key := fmt.Sprintf("key_%d_%d", routineNum, cnt)
				value, err := redis.String(c.Do("get", key))
				if err != nil {
					fmt.Printf("get %s:%v\n", key, err)
				}

				fmt.Printf("g %s\n", value)
				c.Close()
			}
			wg.Done()
		}(i)

	}
	wg.Wait()

}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
