package action

import "fmt"

func validAllowedStringValues(v string, arr []string) error {
	existed := false
	for _, str := range arr {
		if str == v {
			existed = true
			break
		}
	}

	if !existed {
		return fmt.Errorf("the value %s is not in array %#v", v, arr)
	}
	return nil
}
