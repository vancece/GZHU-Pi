package gzhu_jw

import (
	"github.com/astaxie/beego/logs"
	"reflect"
	"testing"
)

func TestWeekHandle(t *testing.T) {
	var data = []struct {
		Input  string
		Output []int
	}{
		{
			Input:  "4周,8-12周(单),14-16周",
			Output: []int{4, 4, 9, 9, 11, 11, 14, 16},
		},
		{
			Input:  "4周,7周(单),8周(单),10-11周(双)",
			Output: []int{4, 4, 7, 7, 10, 10},
		},
		{
			Input:  "4-1周,7-9周(单)",
			Output: []int{4, 1, 7, 7, 9, 9},
		},
		{
			Input:  "1-8周,9-16周",
			Output: []int{1, 8, 9, 16},
		},
	}

	for _, v := range data {
		output := WeekHandle(v.Input)
		if !reflect.DeepEqual(output, v.Output) {
			t.Errorf("want %v but got %v", v.Output, output)
		} else {
			logs.Debug(output)
		}

	}
}
