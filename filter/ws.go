/*
	Reads lsof network files output from stdin and returns number of established websocket connections to <listening_port>

	Requirements
		lsof -nP -i output to stdin
		<listening_port> is required
		go 1.13

	Usage:
		... | go run filter/ws.go <listening_port>
*/

package main

import (
	"fmt"
	"bufio"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: missing port arg")
		return
	}
	port := os.Args[1]
	var n int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pattern := fmt.Sprintf(`:%s->.*ESTABLISHED`, port)
		match, err := regexp.Match(pattern, scanner.Bytes())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error on regex:", err)
			return
		}
		if match {
			n++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error on stdin:", err)
	}
	fmt.Println("* filter/ws.go")
	fmt.Println(n, "websockets")
}
