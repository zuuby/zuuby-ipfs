package utils

import (
	"encoding/json"
)

func StructToJSON(s *struct{}) {
	data, _ = json.Marshal(s)
	return string(data)
}

func JSONToStruct(data []byte, s *struct{}) {
	json.UnMarshal(data, s)
}
