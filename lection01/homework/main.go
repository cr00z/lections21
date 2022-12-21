package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const (
	black = 0
	// red
	// green
	// yellow
	// blue
	// magenta
	// cyan
	white = 7
)

func sandglass(params ...string) error {
	size := 20
	char := 'X'
	color := white
	var err error
	for _, param := range params {
		values := strings.Split(param, "=")
		switch values[0] {
		case "size":
			size, err = strconv.Atoi(values[1])
			if err != nil {
				return errors.Wrap(err, "unknown size")
			}
		case "char":
			if utf8.RuneCountInString(values[1]) > 1 {
				return errors.New("char length more then 1")
			}
			char, _ = utf8.DecodeRuneInString(values[1])
		case "color":
			color, err = strconv.Atoi(values[1])
			if err != nil {
				return errors.Wrap(err, "unknown color")
			}
			if color < black || color > white {
				return errors.New("unknown color")
			}
		default:
			return errors.New("unknown parameter")
		}
	}
	fmt.Print("\033[3", strconv.Itoa(color), "m")
	var out string
	for i := 0; i < size; i++ {
		if i == 0 || i == size-1 {
			out = strings.Repeat(string(char), size)
		} else {
			str := make([]rune, size)
			for j := 0; j < size; j++ {
				str[j] = ' '
			}
			str[i] = char
			str[size-i-1] = char
			out = string(str)
		}
		fmt.Println(out)
	}
	fmt.Print("\033[m")
	return nil
}

func main() {
	err := sandglass("size=11", "char=@", "color=2")
	if err != nil {
		log.Fatal(err)
	}
}
