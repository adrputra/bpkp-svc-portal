package main

import (
	"bpkp-svc-portal/app"
	"os"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	app.Start()
}
