package main

import (
	"mylearning/dao"
	"time"
)

func main() {
	chWrite := make(chan bool)
	chRead := make(chan []dao.User)

	go dao.InsertData(chWrite)

	resultOk := <-chWrite
	if resultOk {
		go dao.SelectData(chRead)
		go dao.Write2file(chRead)
	} else {
		panic(resultOk)
	}
	time.Sleep(100 * time.Millisecond)
}
