package util

import (
	"sort"

	"github.com/luqmanarifin/kentang/model"
)

type Pair struct {
	Key   int
	Value string
}

type ByKey []Pair

func (s ByKey) Len() int {
	return len(s)
}

func (s ByKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByKey) Less(i, j int) bool {
	return s[i].Key > s[j].Key
}

func EntriesToSortedMap(entries []model.Entry) []Pair {
	m := make(map[string]int)
	for _, entry := range entries {
		m[entry.Keyword]++
	}

	var pairs []Pair
	for k, v := range m {
		pairs = append(pairs, Pair{
			Key:   v,
			Value: k,
		})
	}
	sort.Sort(ByKey(pairs))
	return pairs
}
