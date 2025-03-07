package main

import (
	"fmt"

	"github.com/Liphium/magic/wizard/util"
)

func main() {
	fmt.Println("Welcome to Wizard!")

	var tokenString string

	fmt.Print("Type JWT please: ")
	fmt.Scan(&tokenString)

	fmt.Println("JWT:", tokenString)

	util.PostRequestBackend("/api/wizard/init", map[string]interface{}{"token": tokenString, "host": "localhost"})
}
