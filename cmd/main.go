package main

import (
	"github.com/Tuma78/GolangServ/internal"
)



func main(){
	app := application.New()
	app.RunServer()
}
/// export PORT = 8082 && go run .cmd/main.go