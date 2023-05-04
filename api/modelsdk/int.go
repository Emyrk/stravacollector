package modelsdk

import (
	"encoding/json"
	"strconv"
)

type StringInt int64

func (s StringInt) String() string {
	return strconv.FormatInt(int64(s), 10)
}

func (s StringInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *StringInt) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err == nil && str != "" {
		v, err := strconv.ParseInt(str, 10, 64)
		*s = StringInt(v)
		return err
	}
	var v int64
	err = json.Unmarshal(data, &v)
	*s = StringInt(v)
	return err
}
