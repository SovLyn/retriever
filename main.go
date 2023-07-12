package main

import (
	"log"
	"os"
	"strings"
	"time"

	"example.com/m/v2/gdelt"
	"example.com/m/v2/global"
)

func main() {
	// Initialize
	global.Log = log.New(os.Stdout, "INFO:", log.LUTC|log.Ldate|log.Ltime)
	if !strings.HasSuffix(global.TempStorage, "/") {
		global.TempStorage += "/"
	}

	global.Log.Println("program started")

	result, _ := gdelt.GetRawData(time.Date(2023, 7, 10, 0, 0, 0, 0, time.UTC))

	for i := 0; i < 50; i++ {
		println(result[i].String())
	}
}
