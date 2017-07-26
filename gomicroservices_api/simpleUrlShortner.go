package main

// package that converts an id and its associated url into a short url,
// and vice versa
import (
	"bytes"
	_ "github.com/go-sql-driver/mysql"
)

// encodes id int input and encodes params base and dict string for mapping
// returning a short url string
func encode(id int, base int, dict string) string {
	var buffer bytes.Buffer
	for id > 0 {
		digit := id % base
		buffer.WriteString(string(a[digit]))
		id /= base
	}

	str := buffer.String()
	return reverse(str)
}

// takes a str string and maps against dictionary of string
// returning the original integer value
func decode(str string, dict string) int {
	id := 0
	for _, j := range str {

		if 'a' <= j && j <= 'z'{
			id = id * base + int(j) - 'a'
		}
		if 'A' <= j && j <= 'Z' {
			id = id * base + int(j) - 'A' + 26
		}
		if '0' <= j && j <= '9' {
			id = id * base + int(j) - '0' + 52
		}
	}

	return id
}

// reverses the input s and returns a string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
