package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type GeneralResponse struct {
	ResponseType string
	Message      string
	Group        int
}

var size int = 30

var c chan int = make(chan int, size)

func main() {
	for i := 0; i < size; i++ {
		go getReq(i)
	}
	for i := 0; i < size; i++ {
		<-c
	}
}

func getReq(x int) {
	s := strconv.Itoa(x)
	fmt.Println("go getReq: " + s)
	resp, err := http.Get("http://localhost:8000/mm/" + s)
	if err != nil {
		fmt.Println("error!")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	res := GeneralResponse{}
	json.Unmarshal(body, &res)
	fmt.Printf("%d:\t%d\n", x, res.Group)
	c <- x
}
