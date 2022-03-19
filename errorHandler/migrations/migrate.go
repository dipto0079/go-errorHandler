package main

import (
	"blog/category/storage/postgres"
	"log"
)


func main(){
	if err := postgres.Migrate();err !=nil{
		log.Fatal(err)
	}
}