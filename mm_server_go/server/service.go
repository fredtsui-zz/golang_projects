package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var appPath string = "/mm"

var wQueue chan GroupInfo

var baton1, baton2 chan int
var lock chan int
var count int = 0
var members string
var groupnum int = 0
var groupsize int = 5
var port int = 8000

type GeneralResponse struct {
	ResponseType string `json:"responseType,omitempty"`
	Message      string `json:"message,omitempty"`
	Group        int    `json:"group,omitempty"`
}

type GroupInfo struct {
	GroupNum int
	Members  string
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 {
		arg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Println("Usage: ", os.Args[0], " [ port [ group_size ]]")
		} else {
			port = arg
		}
	}
	if len(args) == 2 {
		arg, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Println("Usage: ", os.Args[0], " [ port [ group_size ]]")
		} else {
			groupsize = arg
		}
	}

	router := mux.NewRouter()
	go mmAdmin()
	router.HandleFunc(appPath+"/test/{message}", testHandler).Methods("GET")
	router.HandleFunc(appPath+"/{id}", mmHandler).Methods("GET")
	portStr := strconv.Itoa(port)
	log.Println("service hosting on: ", portStr, " port")
	log.Fatal(http.ListenAndServe(":"+portStr, router))

}

func mmAdmin() { // g is group size
	wQueue = make(chan GroupInfo, groupsize)
	lock = make(chan int, 1) // lock
	lock <- 1
	baton1 = make(chan int, groupsize)
	baton2 = make(chan int, 1)
	for {
		<-baton2
		{
			<-lock
			log.Println(groupnum)
			groupnum = groupnum + 1
			count = 0
			res := GroupInfo{groupnum, members}
			members = ""
			for i := 0; i < groupsize; i++ {
				wQueue <- res
			}
			for i := 0; i < groupsize; i++ {
				<-baton1
			}
			lock <- 1
		}
	}
}

func mmHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.Println("received: " + id)
	baton1 <- 1
	<-lock
	members = members + " " + id
	count = count + 1
	if count == groupsize {
		baton2 <- 1
	}
	lock <- 1
	resp := new(GeneralResponse)
	ginfo := <-wQueue
	resp.Group = ginfo.GroupNum
	resp.Message = ginfo.Members

	json.NewEncoder(w).Encode(resp)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println("message received, " + params["message"])

	resp := new(GeneralResponse)
	resp.ResponseType = "succ"
	json.NewEncoder(w).Encode(resp)
}
