package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const defaultWidth = 240

type (
	config struct {
		Width   int
		Filters map[string]string
	}
)

func (c *config) FromArgs() *config {
	flag.IntVar(&c.Width, "w", defaultWidth, "-w {width}")
	flag.Func("f", "-f key=value", func(x string) error {
		if c.Filters == nil {
			c.Filters = make(map[string]string)
		}
		parts := strings.Split(strings.TrimSpace(x), "=")
		c.Filters[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		return nil
	})

	flag.Parse()
	return c
}

func (c *config) String() string {
	bts, _ := json.MarshalIndent(c, "", "\t")
	return string(bts)
}

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
	cfg := new(config).FromArgs()
	fmt.Println(cfg.String())

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
			var (
				msg  any
				skip bool
			)

			for k, v := range m {
				target, ok := cfg.Filters[k]
				if ok && target != v {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

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
		fmt.Println(strings.Repeat("_", cfg.Width))
	}
}
