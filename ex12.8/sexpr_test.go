// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package sexpr

import (
	"reflect"
	"testing"
	"fmt"
	"bytes"
)

// Test verifies that encoding and decoding a complex data value
// produces an equal result.
//
// The test does not make direct assertions about the encoded output
// because the output depends on map iteration order, which is
// nondeterministic.  The output of the t.Log statements can be
// inspected by running the test with the -v flag:
//
// 	$ go test -v gopl.io/ch12/sexpr
//
func Test(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	// Encode it
	data, err := Marshal(strangelove)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshal() = %s\n", data)
	fmt.Println(string(data))

	// Decode it
	var movie Movie
	if err := Unmarshal(data, &movie); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	t.Logf("Unmarshal() = %+v\n", movie)

	// Check equality.
	if !reflect.DeepEqual(movie, strangelove) {
		t.Fatal("not equal")
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		v bool
		want string
	}{
		{true, "t"},
		{false, "nil"},
	}
	for _, test := range tests {
		data, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%s): %s", err)
		}
		if string(data) != test.want {
			t.Errorf("Marshal(%s) got %s, wanted %s", test.v, data, test.want)
		}
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		v float32
		want string
	}{
		{3.2e9, "3.2e+09"},
		{1.0, "1"},
		{0, "0"},
	}
	for _, test := range tests {
		data, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%s): %s", err)
		}
		if string(data) != test.want {
			t.Errorf("Marshal(%s) got %s, wanted %s", test.v, data, test.want)
		}
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		v float64
		want string
	}{
		{3.2e9, "3.2e+09"},
		{1.0, "1"},
		{0, "0"},
	}
	for _, test := range tests {
		data, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%s): %s", err)
		}
		if string(data) != test.want {
			t.Errorf("Marshal(%s) got %s, wanted %s", test.v, data, test.want)
		}
	}
}

func TestComplex64(t *testing.T) {
	tests := []struct {
		v complex64
		want string
	}{
		{0 + 0i, "#C(0 0)"},
		{3-2i, "#C(3 -2)"},
		{-1e9 + -2.2e9i, "#C(-1e+09 -2.2e+09)"},
	}
	for _, test := range tests {
		data, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%s): %s", err)
		}
		if string(data) != test.want {
			t.Errorf("Marshal(%s) got %s, wanted %s", test.v, data, test.want)
		}
	}
}

func TestComplex128(t *testing.T) {
	tests := []struct {
		v complex128
		want string
	}{
		{0 + 0i, "#C(0 0)"},
		{3-2i, "#C(3 -2)"},
		{-1e9 + -2.2e9i, "#C(-1e+09 -2.2e+09)"},
	}
	for _, test := range tests {
		data, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%s): %s", err)
		}
		if string(data) != test.want {
			t.Errorf("Marshal(%s) got %s, wanted %s", test.v, data, test.want)
		}
	}
}

func TestStreamingDecode(t *testing.T) {
	type Book struct {
		Title, Author string
	}
	book := Book{"Point Counterpoint", "Aldous Huxley"}
	data, err := Marshal(book)
	if err != nil {
		t.Errorf("setting up test: %s", err)
		return
	}
	data = bytes.Repeat(data, 2)
	t.Logf("%s", data)

	dec := NewDecoder(bytes.NewReader(data))
	books := make([]Book, 0)
	for dec.More() {
		var b Book
		err := dec.Decode(&b)
		if err != nil {
			t.Errorf("Error after decoding %s: %s", books, err)
			return
		}
		books = append(books, b)
	}
	want := []Book{book, book}
	if !reflect.DeepEqual(want, books) {
		t.Errorf("Got %s, want %s", books, want)
	}
}
