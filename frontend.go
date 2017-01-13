package main

import (
	"fmt"
	"net/http"
    "io"
    "strings"
    "github.com/garyburd/redigo/redis"
    "encoding/json"
    //"redis"
)

var bidder_chan = make(chan int)
var reducer_chan = make(chan int)
var conn, err = redis.Dial("tcp", ":6379")

type url_data struct {
    url string
    status_code int
}

type Request struct {
    status_msg string
    urls []url_data
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "http://locahost:9090%s Not found", r.URL.Path)
}

func Frontend(w http.ResponseWriter, r *http.Request) {
    if len(r.URL.Path[len("/requests/"):])>0 {
        req_id := (r.URL.Path[len("/requests/"):])
        old_req := "request:"+req_id
        data, _ := conn.Do("GET", old_req)
        fmt.Println(data)
    } else {
        switch r.Method {
            case "POST": 
                    req_data := Request{}
                    reader, _ := r.MultipartReader()
                    b := make([]byte, 100)
                    file_content := ""
                    for {
                        part, err := reader.NextPart()
                        if err == io.EOF{
                                break
                            }
                        for {
                            num, _ := part.Read(b)
                            file_content = file_content + string(b[:num])
                            if num==0{
                                break
                            }
                        }
                        req_id, er := conn.Do("INCR", "requests")
                        if er != nil{}
                        req_data.status_msg = "Bidding"
                        for _, url := range strings.SplitAfter(file_content, "\n"){
                            req_data.urls = append(req_data.urls, url_data{url,0})
                        }
                        data, _ := json.Marshal(req_data)
                        fmt.Println(req_data)
                        fmt.Println(data)
                        new_req := "request:" + string(req_id.(int64))
                        fmt.Println(new_req)
                        conn.Do("SET", new_req, data)
                        w.Header().Set("req_id", string(req_id.(int64)))
                    }
            case "GET":
                    fmt.Fprintf(w, "http://locahost:9090%s Not found on the server", r.URL.Path)
        }
    }
}

func main() {
    go Bidder()
    go Reducer()
    conn.Do("SET", "requests", 1)
    defer conn.Close()
    http.HandleFunc("/", NotFound)
    http.HandleFunc("/requests/", Frontend)
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        fmt.Println("ListenAndServe: ", err)
    }
}