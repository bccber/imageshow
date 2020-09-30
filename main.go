package main

import (
	"imageshow/routers"
)

func main() {
	router := routers.InitRouter()
	router.Run(":8080")
}
