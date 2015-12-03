package main

import (
	"log"
	"os"

	"github.com/arteev/zbarnet/app"
)

//TODO: write barcode to pipe out
//TODO: Optional Raw format without base64

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			os.Exit(1)
		}
	}()
	app.New().Run()
}
