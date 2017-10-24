package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
)

func task() {
	fmt.Println("I am runnning worker task.")
}

func main() {
	fmt.Println("monitor run")

	gocron.Every(1).Minute().Do(task)
	gocron.Start()
}


