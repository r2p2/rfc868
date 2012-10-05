package rfc868

import (
	"fmt"
	"testing"
)

func TestRFC868(t *testing.T) {
	go ServeTime("localhost:8888")

	/*ctime, err := RequestTime("localhost:8888")
	if err != nil {
		t.Error(err)
	}
	*/
	// todo check if time is correct
}

func BenchmarkRFC868(b *testing.B) {
	go ServeTime("localhost:8888")

	c, err := NewClient("localhost:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < b.N; i++ {
		c.RequestTime()
	}

}

// You have to run a server separatly.
func BenchmarkClientI(b *testing.B) {
	c, err := NewClient("127.0.0.1:1024")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < b.N; i++ {
		c.RequestTime()
	}
}
