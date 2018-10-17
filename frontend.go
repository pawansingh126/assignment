package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/http"
	"strconv"
	"strings"
	//"redis"
)

var bidder_chan = make(chan int, 100)
var reducer_chan = make(chan int, 100)
var conn, err = redis.Dial("tcp", ":6379")

type Url_data struct {
	Url         string
	Status_code int
}

type Request struct {
	Status_msg string
	Urls       []Url_data
	Result     map[int][]string
}

func ReadFile(r *http.Request) (*Request, error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 100)
	file_content := ""

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		for {
			num, _ := part.Read(buffer)
			file_content += string(buffer[:num])
			if num == 0 {
				break
			}
		}
	}

	URLs := strings.Split(file_content, "\n")
	req_data := Request{Status_msg: "Bidding"}

	for _, u := range URLs {
		req_data.Urls = append(req_data.Urls, Url_data{Url: u})
	}
	return &req_data, nil
}

func RedisSet(r *Request, req_str string) error {
	marshalled, er := json.Marshal(r)
	if er != nil {
		fmt.Println(er)
		panic(er)
	}
	_, err := conn.Do("SET", req_str, marshalled)
	return err
}

func RedisGet(req_str string) (*Request, error) {
	req_data, er := redis.Bytes(conn.Do("GET", req_str))
	if er != nil {
		panic(er)
	}
	req := &Request{}
	err := json.Unmarshal(req_data, req)
	if err != nil {
		panic(err)
	}
	return req, nil
}

func RedisIncr() (int64, error) {
	id, er := redis.Int64(conn.Do("INCR", "requests"))
	if er != nil {
		panic(err)
	}
	return id, nil
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http://locahost:9090%s Not found", r.URL.Path)
}

func Frontend(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path[len("/requests/"):]) > 0 {

		req_id := "request:" + r.URL.Path[len("/requests/"):]
		req_data, err := RedisGet(req_id)
		if err != nil {
			panic(err)
		}
		if req_data.Status_msg == "Complete" {
			w.Header().Set("Content-Type", "application/json")
			res, er := json.MarshalIndent(req_data.Result, "", "    ")
			if er != nil {
				panic(er)
			}
			w.Write(res)
		} else {
			w.Header().Set("Retry-After", fmt.Sprintf("%v", 10))
			fmt.Fprintln(w, "Current Status:", req_data.Status_msg)
		}

	} else {
		switch r.Method {
		case "POST":
			req_data, err := ReadFile(r)
			if err != nil {
				fmt.Fprintf(w, "Improper Data File Sent")
				return
			}
			incr, _ := RedisIncr()
			incr_str := strconv.FormatInt(incr, 10)
			w.Header().Set("req_id", fmt.Sprintf("%v", incr))
			new_req := "request:" + incr_str
			err = RedisSet(req_data, new_req)
			if err != nil {
				panic(err)
			}
			bidder_chan <- int(incr)
		case "GET":
			fmt.Fprintf(w, "http://locahost:9090%s Not found on the server", r.URL.Path)
		}
	}
}

func main() {
	go Bidder()
	go Reducer()

	conn.Do("SET", "requests", 0)
	defer conn.Close()
	http.HandleFunc("/", NotFound)
	http.HandleFunc("/requests/", Frontend)

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
