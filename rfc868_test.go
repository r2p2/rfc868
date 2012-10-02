package rfc868

import (
	"fmt"
	"testing"
)

func TestRFC868(t *testing.T) {
	go ServeTime("localhost:8888")

	ctime, err := RequestTime("localhost:8888")
	if err != nil {
		t.Error(err)
	}
	
	// todo check if time could be ok
}

func BenchmarkRFC868(b *testing.B) {
	go ServeTime("localhost:8888")

	for i := 0; i < b.N; i++ {
		RequestTime("localhost:8888")
	}

}
