package main

import (
	"log"
	"time"
	"fmt"
	"github.com/smcduck/xsys/xmutex"
)

func main() {
	lock := xmutex.NewSingleInstanceLock("your-app-name")
	defer lock.UnLock()

	ok, err := lock.IsUnique()
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		fmt.Println("Is single process")
	} else {
		fmt.Println("Another process running")
	}

	time.Sleep(60 * time.Second)
	log.Println("finished")
}