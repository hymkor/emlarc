package main

import (
	"io"
	"os"
	"encoding/base64"
)

func main() {
	r := base64.NewDecoder(base64.StdEncoding,os.Stdin)
	io.Copy(os.Stdout, r)
}
