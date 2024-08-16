package main

import (
	"fmt"
	"log"

	"github.com/exam-5/Car-Wash-Api-Gateway/api"
	_ "github.com/exam-5/Car-Wash-Api-Gateway/docs"
)

func main() {
	r := api.NewGin()
	fmt.Println("Server started on port:9080")
	err := r.Run(":9080")
	if err != nil {
		log.Fatal("Error while running server: ", err.Error())
	}
}
