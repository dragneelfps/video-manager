package utils

import "encoding/json"

func StructToJsonString(s any) string {
	bs, _ := json.Marshal(s)
	return string(bs)
}
