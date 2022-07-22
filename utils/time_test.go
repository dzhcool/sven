package utils

import "testing"

func TestDateToTime(t *testing.T) {
	a := DateToTime("2022-07-11")
	b := DateToTime("2022-07-11 11:19:00")
	if a != b {
		t.Fatalf("date to time failed. a:%d, b:%d", a, b)
		return
	}
	t.Log("date to time success")
}
