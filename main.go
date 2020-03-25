package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
)

func makeBoundString() (string, error) {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", buf[:]), nil
}

func mains(args []string) error {
	bound, err := makeBoundString()
	if err != nil {
		return err
	}
	fmt.Printf("Content-Type: multipart/mixed; boundary=\"%s\"\n\n", bound)

	w := multipart.NewWriter(os.Stdout)
	defer w.Close()
	w.SetBoundary(bound)

	for _, arg := range args {
		header := textproto.MIMEHeader{}

		header.Add("Content-Disposition",
			fmt.Sprintf(`form-data; name="%[1]s"; filename="%[1]s"`, arg))
		header.Add("Content-Type", "application/octet-stream")
		header.Add("Content-Transfer-Encoding", "base64")

		w1, err := w.CreatePart(header)
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
