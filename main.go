package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const defaultWidth = 240

func getWidth() int64 {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	fmt.Printf("out: %#v\n", string(out))
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		log.Println(err)
		return defaultWidth
	}
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	width, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return width
}

func padRight(str string, fill string, length int) string {
	dt := length - len(str)
	if dt <= 0 {
		return str
	}
	return str + strings.Repeat(fill, dt)
}

func main() {
	width := getWidth()

	fmt.Println("width: ", width)
	buf := bufio.NewScanner(os.Stdin)
	msgKey := padRight("message", " ", 40)

	for buf.Scan() {
		var err error
		m := make(map[string]any)
		if err = buf.Err(); err != nil {
			fmt.Println("ERROR:", err)
		}
		line := buf.Bytes()
		if err = json.Unmarshal(line, &m); err != nil {
			fmt.Println(string(line))
		} else {
			var msg any
			for k, v := range m {
				if k == "message" {
					msg = v
				} else {
					fmt.Printf("%s%v\n", padRight(k, " ", 40), v)
				}
			}

			if msg != nil {
				fmt.Printf("%s%v\n", msgKey, msg)
			}

		}
		fmt.Println(strings.Repeat("_", int(width)))

	}
}
