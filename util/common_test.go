package util

import (
	"log"
	"testing"

	"github.com/luqmanarifin/kentang/model"
)

func TestEntriesToSortedMap(t *testing.T) {
	cok := make(map[string]int)
	log.Printf("%d\n", cok["asu"])
	cok["asu"] = 1
	log.Printf("%d\n", cok["asu"])
	entries := []model.Entry{
		model.Entry{Keyword: "niki"},
		model.Entry{Keyword: "niki"},
		model.Entry{Keyword: "luq"},
		model.Entry{Keyword: "luq"},
		model.Entry{Keyword: "niki"},
		model.Entry{Keyword: "niki"},
		model.Entry{Keyword: "bird"},
		model.Entry{Keyword: "asu"},
	}
	m := EntriesToSortedMap(entries)
	for key, value := range m {
		log.Printf("%s: %d\n", key, value)
	}
}
