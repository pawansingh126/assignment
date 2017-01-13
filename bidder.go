package main

import (
	"fmt"
	"time"
)

func Bidder(){
	for {
		for len(bidder_chan) !=0 {
			fmt.Println("Kuch to H!!")
		}
		//fmt.Println("")
		time.Sleep(2000 * time.Millisecond)
	}
}