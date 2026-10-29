package main

import (
	"errorHandler/errorHandler/storage/postgres"
	"log"
)


func main(){
	if err := postgres.Migrate();err !=nil{
		log.Fatal(err)
	}
}