package limitmap

import (
	"fmt"
	"tebot/test"
	"testing"
	"time"
)

func TestNewLimitMap(t *testing.T) {
	type args struct {
		maxSize int
	}
	tests := []struct {
		name string
		args args
		want *LimitMap
	}{
		{args: args{100}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Init()
			lm := NewLimitMap(200, 100)

			for i := 1; i <= 105; i++ {
				lm.Add(i, fmt.Sprintf("Value %d", i%10)) // this will generate 10 unique strings
			}

			lm.Display()
			//v := lm.Get(101)
			//fmt.Println(v)
			time.Sleep(20 * time.Second)
		})
	}
}
