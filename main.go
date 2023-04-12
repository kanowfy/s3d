package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrIncorrectFormat    = errors.New("incorrect format")
)

type config struct {
	pattern string
	infile  string
}

type replaceOption struct {
	from     string
	to       string
	isGlobal bool
}

func readArguments(args []string) (*config, error) {
	if len(args) < 3 {
		return nil, ErrNotEnoughArguments
	}

	return &config{
		pattern: args[1],
		infile:  args[2],
	}, nil
}

func parsePattern(p string) (*replaceOption, error) {
	tokens := strings.Split(p, "/")
	if len(tokens) < 4 {
		return nil, ErrIncorrectFormat
	}
	opts := &replaceOption{
		from: tokens[1],
		to:   tokens[2],
	}
	if strings.ToLower(tokens[3]) == "g" {
		opts.isGlobal = true
	} else {
		opts.isGlobal = false
	}

	return opts, nil
}

func replace(data io.Reader, opts *replaceOption) string {
	var sb strings.Builder
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, opts.from) {
			if opts.isGlobal {
				str = strings.ReplaceAll(str, opts.from, opts.to)
			} else {
				str = strings.Replace(str, opts.from, opts.to, 1)
			}
		}

		sb.WriteString(fmt.Sprintf("%s\n", str))
	}

	return sb.String()
}

func main() {
	config, err := readArguments(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	replaceOpts, err := parsePattern(config.pattern)
	if err != nil {
		if errors.Is(err, ErrIncorrectFormat) {
			fmt.Println("incorrect format, should be s/from/to/[g]")
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	f, err := os.Open(config.infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	output := replace(f, replaceOpts)
	fmt.Print(output)
}
