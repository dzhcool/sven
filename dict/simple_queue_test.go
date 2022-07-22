package dict

import (
	"testing"
)

func TestCross(t *testing.T) {
	queue := NewSimpleQueue(10)
	queue.Insert("1")
	v, err := queue.Front()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("ok:", v)
}

func Test(t *testing.T) {
	queue := NewSimpleQueue(10)
	queue.Insert("1")
	queue.Insert("2")
	queue.Insert("3")

	for i := 0; i < 4; i++ {
		v, err := queue.Front()
		if err != nil {
			t.Log(err)
			break
		}
		t.Log(i, " ok:", v)
	}
	return
}
