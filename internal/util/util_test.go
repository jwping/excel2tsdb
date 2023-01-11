package util

import (
	"fmt"
	"testing"
)

func TestGetExcel(t *testing.T) {
	dataMetrics, err := GetMetrics("/home/jwping/gopath/excel2tsdb/testdir", false, nil, []int{1}, 1, 4)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("dataMetrics: %+v\n", dataMetrics)
}
