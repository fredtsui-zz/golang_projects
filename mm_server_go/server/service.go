package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var appPath string = "/mm"

var wQueue chan int

var baton1 chan int
var lock chan int
var count int = 0
var groupnum int = 0

type GeneralResponse struct {
	ResponseType string   `json:"responseType,omitempty"`
	Message      []string `json:"message,omitempty"`
	Group        int      `json:"group,omitempty"`
}

func main() {
	router := mux.NewRouter()
	go mmAdmin(3)
	router.HandleFunc(appPath+"/test/{message}", testHandler).Methods("GET")
	router.HandleFunc(appPath+"/{id}", mmHandler).Methods("GET")
	log.Println("service hosting on: 8000 port")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func mmAdmin(g int) { // g is group size
	wQueue = make(chan int, g)
	lock = make(chan int, 1) // lock
	lock <- 1
	baton1 = make(chan int, g)
	for {
		if count == g {
			<-lock
			log.Println(groupnum)
			groupnum = groupnum + 1
			count = 0
			for i := 0; i < g; i++ {
				wQueue <- groupnum
			}
			for i := 0; i < g; i++ {
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
	count = count + 1
	lock <- 1
	resp := new(GeneralResponse)
	resp.Group = <-wQueue
	//log.Println(result)
	//resp.Message = result

	json.NewEncoder(w).Encode(resp)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println("message received, " + params["message"])

	resp := new(GeneralResponse)
	resp.ResponseType = "succ"
	//resp.Message = params["message"]
	json.NewEncoder(w).Encode(resp)
}
