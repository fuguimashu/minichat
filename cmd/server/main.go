package main

import "minichat/internal/server"

func main() {
	server := server.New("0.0.0.0", 2577)
	server.Start()
}
