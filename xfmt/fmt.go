package xfmt

import (
	"encoding/json"
	"fmt"
)

func MarshalAndPrintln(x interface{}) {
	buf, err := json.Marshal(x)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
}