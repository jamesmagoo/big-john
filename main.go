package main

func main() {
	server := NewAPIServer(":5001")
	server.Run()
}
