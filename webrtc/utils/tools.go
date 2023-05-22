package utils

import (
	"encoding/base64"
	"encoding/json"
)

func Encode(obj any) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(marshal)
}

func Decode(in string, obj any) {
	decodeString, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(decodeString, obj)
	if err != nil {
		panic(err)
	}

}
