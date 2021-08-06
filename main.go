package main

import (
	"github.com/joho/godotenv"

	"github.com/bgildson/unico-challenge/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
