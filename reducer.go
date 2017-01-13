package main

import (
	"fmt"
	"time"
)

func Reducer(){
	for {
		for len(reducer_chan) !=0 {
			fmt.Println("Kuch to H!!")

		}
		//fmt.Println("Kuch Nhi H!!")
		time.Sleep(2000 * time.Millisecond)
	}
}