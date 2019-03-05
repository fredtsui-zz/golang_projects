package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GeneralResponse struct {
	ResponseType string
	Message      string
	Group        int
}

var size int = 30
var maxTimeout int = 1000

var c chan int = make(chan int, size)
var b chan int64 = make(chan int64, size)

func main() {
	args := os.Args
	if len(args) > 0 {
		arg1, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Usage: " + os.Args[0] + " [numRequests [ maxDelay ]]")
		} else {
			size = arg1
		}
		if len(args) > 1 {
			arg2, err2 := strconv.Atoi(os.Args[2])
			if err2 != nil {
				fmt.Println("Usage: " + os.Args[0] + " [numRequests [ maxDelay (in Milliseconds) ]]")
			} else {
				maxTimeout = arg2
			}
		}
	}
	for i := 0; i < size; i++ {
		go getReq(i)
	}
	for i := 0; i < size; i++ {
		seed := (time.Now().UnixNano() / int64(time.Millisecond))
		seed = seed % 100
		b <- seed
	}
	// wait for go threads to finish before this returns
	for i := 0; i < size; i++ {
		<-c
	}
}

func getReq(x int) {
	seed := <-b
	rand.NewSource(seed)
	time.Sleep(time.Duration(rand.Intn(maxTimeout)) * time.Millisecond)
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
