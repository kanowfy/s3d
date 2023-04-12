package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadArgs(t *testing.T) {
	args := []string{"s3d", "s/hey/yo/", "in.txt"}
	got, err := readArguments(args)
	want := &config{
		pattern: "s/hey/yo/",
		infile:  "in.txt",
	}
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestReadArgsErr(t *testing.T) {
	args := []string{"s3d", "s/hey/yo/"}
	_, err := readArguments(args)
	if err == nil {
		t.Fatalf("expected to fail from: %v", ErrNotEnoughArguments)
	}
}

func TestParsePattern(t *testing.T) {
	patternTests := []struct {
		name    string
		pattern string
		want    *replaceOption
	}{
		{
			name:    "first instance match",
			pattern: "s/hey/zo/",
			want:    &replaceOption{from: "hey", to: "zo", isGlobal: false},
		},

		{
			name:    "global match",
			pattern: "s/lmao/xd/g",
			want:    &replaceOption{from: "lmao", to: "xd", isGlobal: true},
		},
	}

	for _, tt := range patternTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePattern(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	replaceTests := []struct {
		name string
		data io.Reader
		opts *replaceOption
		want string
	}{
		{
			name: "first instance match",
			data: strings.NewReader("one two one three\ntwo one two three\n"),
			opts: &replaceOption{from: "one", to: "five", isGlobal: false},
			want: "five two one three\ntwo five two three\n",
		},
		{
			name: "global match",
			data: strings.NewReader("one two one three\ntwo one two three\n"),
			opts: &replaceOption{from: "one", to: "five", isGlobal: true},
			want: "five two five three\ntwo five two three\n",
		},
	}

	for _, tt := range replaceTests {
		t.Run(tt.name, func(t *testing.T) {
			got := replace(tt.data, tt.opts)
			if tt.want != got {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
