package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

func main() {
	result := ToUniqueStringSlice("['name':'pdx','age':10]")
	fmt.Println(result)
}

// ToUniqueStringSlice casts `value` to a slice of non-zero unique strings.
func ToUniqueStringSlice(value any) (result []string) {
	switch val := value.(type) {
	case nil:
		// nothing to cast
	case []string:
		result = val
	case string:
		if val == "" {
			break
		}

		// check if it is a json encoded array of strings
		if strings.Contains(val, "[") {
			if err := json.Unmarshal([]byte(val), &result); err != nil {
				// not a json array, just add the string as single array element
				result = append(result, val)
			}
		} else {
			// just add the string as single array element
			result = append(result, val)
		}
	case json.Marshaler: // eg. JsonArray
		raw, _ := val.MarshalJSON()
		_ = json.Unmarshal(raw, &result)
	default:
		result = cast.ToStringSlice(value)
	}

	return result
}
