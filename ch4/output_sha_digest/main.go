package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)


func main() {
	stdin_scanner := bufio.NewScanner(os.Stdin)
	algorithmName := flag.String("a", "sha256", "The name of secure hash algorithm")
	flag.Parse()
	for {
		stdin_scanner.Scan()
		input_argument := stdin_scanner.Text()
		if len(input_argument) == 0 {
			break
		}
		switch *algorithmName {
		case "sha256":
			fmt.Printf("%x\n", sha256.Sum256([]byte(input_argument)))
		case "sha384":
			fmt.Printf("%x\n", sha512.Sum384([]byte(input_argument)))
		case "sha512":
			fmt.Printf("%x\n", sha512.Sum512([]byte(input_argument)))
		default:
			fmt.Println("Argument 'a' should be 'sha256', 'sha384' or 'sha512'")
			return
		}
	}
}
