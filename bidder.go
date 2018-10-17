package main

import (
	"net/http"
	"strconv"
	"time"
)

func GetUrl(url string, u *Url_data, c chan bool) {
	timeout := time.Duration(2000 * time.Millisecond)

	client := http.Client{Timeout: timeout}

	res, err := client.Get(url)
	if err != nil {
		u.Status_code = 0
	} else {
		defer res.Body.Close()
		u.Status_code = res.StatusCode
	}
	<-c
}

func Bidder() {
	for {
		for len(bidder_chan) != 0 {

			incr := <-bidder_chan
			req_str := "request:" + strconv.Itoa(incr)

			req_data, _ := RedisGet(req_str)

			is_complete := make(chan bool, len(req_data.Urls))

			for i, u := range req_data.Urls {
				is_complete <- true
				go GetUrl(u.Url, &req_data.Urls[i], is_complete)
			}

			for len(is_complete) != 0 {
				time.Sleep(500 * time.Millisecond)
			}
			req_data.Status_msg = "Reducing"
			er := RedisSet(req_data, req_str)
			if er != nil {
				panic(er)
			}
			reducer_chan <- incr
		}
		time.Sleep(2000 * time.Millisecond)
	}
}
