package main

import (
	"testing"
	"strconv"
)

func TestEncode(t *testing.T) {
	for _, c := range []struct {
		in, want string
	}{
		{"1", "b"},
		{"26", "A"},
		{"51", "Z"},
		{"61", "9"},
	}{
		n, _:= strconv.Atoi(c.in)
		got := encode(n, base, a)
		if got != c.want {
			t.Errorf("Encode(%q) == %q, want %q", c.in, got, c.want)
		}
	}


}

func TestReverse(t *testing.T) {
	for _, c := range []struct {
		in, want string
	} {
		{ "44$4", "4$44"},
		{"caruso", "osurac"},
		{"Hello, 世界", "界世 ,olleH"},
	} {
		got := reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}