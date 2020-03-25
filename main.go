package main

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
)

func mains(args []string) error {
	w := multipart.NewWriter(os.Stdout)
	defer w.Close()
	for _, arg := range args {
		w1, err := w.CreateFormFile(arg, arg)
		if err != nil {
			return fmt.Errorf("%s: %w", arg, err)
		}
		fd, err := os.Open(arg)
		if err != nil {
			return fmt.Errorf("%s: %w", arg, err)
		}
		dst := make([]byte, 76)
		src := make([]byte, base64.StdEncoding.DecodedLen(len(dst)))
		for {
			n, err := fd.Read(src)
			base64.StdEncoding.Encode(dst, src)
			w1.Write(dst[:base64.StdEncoding.EncodedLen(n)])
			w1.Write([]byte{'\r', '\n'})
			if err != nil {
				break
			}
		}
		fd.Close()
		if err != nil {
			return fmt.Errorf("%s: %w", arg, err)
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
