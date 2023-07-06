package main

import "log"

func main() {

	storage, err := NewPostgreStore("jamesc", "gobank", "''")
	if err != nil {
		panic(err)
	}

	if err := storage.Init(); err != nil {

		log.Fatal(err)
	}
	server := NewAPIServer(":8080", storage)

	server.Run()
}
