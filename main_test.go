package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFirstUse(t *testing.T) {
	target := "./test/target/target.txt"
	dir := "./test"
	dirEndOfSlash := dir + "/"

	result := isFirstUse(target, dir)
	if !result {
		t.Errorf("expected %v, but %v", true, result)
	}

	result = isFirstUse(target, dirEndOfSlash)
	if !result {
		t.Errorf("expected %v, but %v", true, result)
	}

	dist := "./test/target.txt"
	err := os.Link(target, dist)
	if err != nil {
		t.Fatal("could not create hard link")
	}
	defer func() {
		err = os.Remove(dist)
		if err != nil {
			t.Fatal("could not delete hard link")
		}
	}()

	result = isFirstUse(target, dir)
	if result {
		t.Errorf("expected %v, but %v", false, result)
	}

	result = isFirstUse(target, dirEndOfSlash)
	if result {
		t.Errorf("expected %v, but %v", false, result)
	}
}

func TestHasOriginalFile(t *testing.T) {
	source := "./test/target/target.txt"
	dist := "./test/target.txt"
	dir := "./test"

	result := hasOrginalFile(source, dir)
	if result {
		t.Errorf("expected %v, but %v", false, result)
	}

	err := os.Link(source, dist)
	if err != nil {
		t.Fatal("could not create hard link")
	}
	defer func() {
		err = os.Remove(dist)
		if err != nil {
			t.Fatal("could not delete hard link")
		}
	}()

	result = hasOrginalFile(source, dir)
	if !result {
		t.Errorf("expected %v, but %v", true, result)
	}

}

func TestIsSameFile(t *testing.T) {
	source := "./test/target/target.txt"
	dist := "./test/target.txt"

	err := os.Link(source, dist)
	if err != nil {
		t.Fatal("could not create hard link")
	}

	tests := []struct {
		source   string
		dist     string
		expected bool
	}{
		{source, source, true},
		{source, dist, true},
		{source, "./main.go", false},
	}

	for _, test := range tests {
		result := isSameFile(test.source, test.dist)
		if test.expected != result {
			t.Errorf("expected %v, but %v", test.expected, result)
		}
	}

	err = os.Rename(dist, dist+".txt")
	if err != nil {
		t.Fatal("could not rename")
	}
	defer func() {
		err = os.Remove(dist + ".txt")
		if err != nil {
			t.Fatal("could not delete hard link")
		}
	}()

	result := isSameFile(source, dist+".txt")
	if !result {
		t.Errorf("expected true, but %v", result)
	}
}

func TestIsSameExt(t *testing.T) {
	source := "./test/target/target.txt"
	dist := "./test/target.txt"

	err := os.Link(source, dist)
	if err != nil {
		t.Fatal("could not create hard link")
	}
	defer func() {
		err = os.Remove(dist)
		if err != nil {
			t.Fatal("could not delete hard link")
		}
	}()

	tests := []struct {
		source   string
		dist     string
		expected bool
	}{
		{source, dist, true},
		{source, dist + "/", false},
		{source, "./main.go", false},
	}

	for _, test := range tests {
		ext := filepath.Ext(test.source)
		result := isSameExt(ext, test.dist)
		if test.expected != result {
			t.Errorf("expected %v, but %v", test.expected, result)
		}
	}

}

func TestNormarizeDir(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/test", "/test/"},
		{"/test/", "/test/"},
		{"/test/test", "/test/test/"},
		{"/test/test/", "/test/test/"},
	}

	for _, test := range tests {
		result := normarizeDir(test.input)
		if filepath.FromSlash(test.expected) != result {
			t.Errorf("expected %v, but %v", test.expected, result)
		}
	}
}
