package xclock

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateValid(t *testing.T) {
	if DateValid(2018, 2, 31) {
		t.Error("DateValid() error")
		return
	}
	if DateValid(2018, 11, 31) {
		t.Error("DateValid() error")
		return
	}
	if !DateValid(2018, 1, 31) {
		t.Error("DateValid() error")
		return
	}
}

func TestDate_StringYYYYMMDD(t *testing.T) {
	date, err := NewDate(2018, 3, 8, *time.UTC)
	if err != nil {
		t.Error(err)
		return
	}
	dateString := date.StringYYYYMMDD()
	expected := "20180308"
	if dateString != expected {
		t.Errorf("Correct date string %s, but get %s", expected, date.StringYYYYMMDD())
		return
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	dt := Date{}
	b, err := json.Marshal(dt)
	if err != nil {
		t.Error(err)
		return
	}
	if string(b) != "\"\"" {
		t.Errorf("Date json.Marshal(Date{}) error, returns %s", string(b))
		return
	}
}

type test_item struct {
	InDate Date  `json:"InDate"`
}

func TestDate_UnmarshalJSON(t *testing.T) {
	s := `{"InDate":"2018-05-01"}`
	i := test_item{}
	err := json.Unmarshal([]byte(s), &i)
	if err != nil {
		t.Error(err)
		return
	}
	if i.InDate.String() != "2018-05-01" {
		t.Errorf("Date json.Unmarshal() error, returns '%s'", i.InDate.String())
		return
	}
}
