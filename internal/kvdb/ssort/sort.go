package ssort

import (
	"sort"
	"strconv"
	"strings"
)

const (
	Sep = "_"
)

type timeSuffixSlices []string

func (s timeSuffixSlices) Len() int {
	return len(s)
}

func (s timeSuffixSlices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s timeSuffixSlices) Less(i, j int) bool {
	suffixI := getKeySuffix(s[i])
	suffixJ := getKeySuffix(s[j])

	suffixI64, _ := strconv.ParseInt(suffixI, 10, 64)
	suffixJ64, _ := strconv.ParseInt(suffixJ, 10, 64)

	ret := suffixI64 < suffixJ64
	//fmt.Printf(">>>>>>i:%s, j:%s, ret:%v, suffixI:%d, suffixJ:%d\n", s[i], s[j], ret, suffixI64, suffixJ64)
	return ret
}

func getKeySuffix(key string) string {
	splitKey := strings.Split(key, Sep)
	lenk := len(splitKey)
	return splitKey[lenk-1]
}

func (s timeSuffixSlices) Sort() {
	sort.Stable(s)
}

type SSort func(keys []string)

func SuffixSort(keys []string) {
	suffixKeys := timeSuffixSlices(keys)
	suffixKeys.Sort()
}

func (s timeSuffixSlices) SuffixSort() {
	sort.Stable(s)
}
