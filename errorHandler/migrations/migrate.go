package main

import (
	"github.com/dipto0079/go-errorHandler/errorHandler/storage/postgres"
	"log"
)


func main(){
	if err := postgres.Migrate();err !=nil{
		log.Fatal(err)
	}
}