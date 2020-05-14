package gqueue

import (
	"github.com/beeker1121/goque"
	"log"
)

var Queue *goque.Queue

func Setup() {
	var initErr error
	Queue, initErr = goque.OpenQueue("data_dir")
	if initErr != nil {
		log.Fatal(initErr)
	}
}

func Close() {
	Queue.Close()
}
