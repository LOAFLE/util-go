package json

import (
	"encoding/json"
	"fmt"
	"reflect"

	our "git.loafle.net/overflow/util-go/reflect"
)

// SetValueWithJSONStringArray set the value of json string array
// raw([]byte) is ["1", {"a": 1}, [1, 2]]
// targets([]interface{}) is array of pointer ex) *int, *string, *[], *map, *struct
func SetValueWithJSONStringArrayBytes(raw []byte, targets []interface{}) error {
	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return err
	}

	return SetValueWithJSONStringArray(values, targets)
}

func SetValueWithJSONStringArray(values []string, targets []interface{}) error {
	if len(targets) != len(values) {
		return fmt.Errorf("Count of raw[%d] and targets[%d] is not same", len(values), len(targets))
	}

	for indexI := 0; indexI < len(values); indexI++ {
		target := targets[indexI]
		value := values[indexI]

		if reflect.Ptr != reflect.TypeOf(target).Kind() {
			return fmt.Errorf("Type of target[%d] must be ptr but is %s, value=%s", indexI, reflect.TypeOf(target).Kind(), value)
		}

		switch reflect.TypeOf(target).Elem().Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
			if err := json.Unmarshal([]byte(value), &target); nil != err {
				return err
			}
		case reflect.Ptr:
			return fmt.Errorf("Type of target[%d] cannot be double ptr, value=%s", indexI, value)
		default:
			cv, err := our.ConvertToType(value, reflect.TypeOf(target).Elem())
			if nil != err {
				return fmt.Errorf("Type conversion of value[%s] has been failed to %s[%d]", value, reflect.TypeOf(target).Elem().Kind(), indexI)
			}

			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(cv))
		}
	}

	return nil
}
