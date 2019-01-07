package json

import (
	"encoding/json"
	"strconv"
)

func NumberToInt(n json.Number) (int, error) {
	n64, err := strconv.ParseInt(n.String(), 10, 32)
	if nil != err {
		return 0, err
	}
	return int(n64), nil
}
