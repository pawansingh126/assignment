package main

import (
	"strconv"
	"time"
)

func Reducer() {
	for {
		for len(reducer_chan) != 0 {
			incr := <-reducer_chan
			req_str := "request:" + strconv.Itoa(incr)

			req_data, _ := RedisGet(req_str)
			req_data.Result = make(map[int][]string)

			for _, u := range req_data.Urls {
				req_data.Result[u.Status_code] = append(req_data.Result[u.Status_code], u.Url)
			}

			req_data.Status_msg = "Complete"
			er := RedisSet(req_data, req_str)
			if er != nil {
				panic(er)
			}
		}
		time.Sleep(2000 * time.Millisecond)
	}
}
