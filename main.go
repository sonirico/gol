package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
	Opcode string

	Filter struct {
		Key   string
		Value string
		Op    Opcode
	}

	config struct {
		Width   int
		Filters map[string]Filter
	}
)

const (
	OpEqual    = "="
	OpNotEqual = "!="
)

func (op Opcode) String() string {
	return string(op)
}

func (f Filter) Match(other any) bool {
	switch f.Op {
	case OpEqual:
		return other == f.Value
	case OpNotEqual:
		return other != f.Value
	}

	return false
}

func ParseOpCode(x string) (Opcode, error) {
	switch {
	case strings.Contains(x, "!="):
		return OpNotEqual, nil
	case strings.Contains(x, "<>"):
		return OpNotEqual, nil
	case strings.Contains(x, "="):
		return OpEqual, nil
	}
	return "", errors.New("invalid operator")
}

func ParseFilter(x string) (f Filter, err error) {
	x = strings.TrimSpace(x)
	op, err := ParseOpCode(x)
	if err != nil {
		return f, err
	}
	parts := strings.Split(strings.TrimSpace(x), op.String())
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return Filter{
		Key:   key,
		Value: value,
		Op:    op,
	}, nil
}

func (c *config) FromArgs() *config {
	flag.IntVar(&c.Width, "w", defaultWidth, "-w {width}")
	flag.Func("f", "-f key=value", func(x string) error {
		if c.Filters == nil {
			c.Filters = make(map[string]Filter)
		}

		filter, err := ParseFilter(x)
		if err != nil {
			return err
		}

		c.Filters[filter.Key] = filter

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
				filter, ok := cfg.Filters[k]
				if ok && !filter.Match(v) {
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
