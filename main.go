package main

import "github.com/sjjwantfish/web-server/web"

func main() {
	server := web.HTTPServer{
		Addr: ":8080",
	}
	server.Start()
}
