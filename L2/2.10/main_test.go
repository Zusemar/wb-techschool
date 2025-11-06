package main

import (
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	opts := parseFlags([]string{"-nru"})
	if !opts.n || !opts.r || !opts.u {
		t.Errorf("parseFlags failed: %+v", opts)
	}

	opts = parseFlags([]string{"-k", "2"})
	if !opts.k || opts.N != 2 {
		t.Errorf("parseFlags -k failed: %+v", opts)
	}
}

func TestSortStrBasic(t *testing.T) {
	lines := []string{"banana", "apple", "cherry"}
	got := sortStr(lines, sortOptions{})
	want := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestSortStrNumericReverse(t *testing.T) {
	lines := []string{"10", "2", "5"}
	got := sortStr(lines, sortOptions{n: true, r: true})
	want := []string{"10", "5", "2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestSortStrColumn(t *testing.T) {
	lines := []string{"apple\t3", "banana\t1", "cherry\t2"}
	got := sortStr(lines, sortOptions{k: true, N: 2})
	want := []string{"banana\t1", "cherry\t2", "apple\t3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestUnique(t *testing.T) {
	in := []string{"a", "a", "b", "b", "c"}
	want := []string{"a", "b", "c"}
	got := unique(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
