package modelsdk

import (
	"github.com/Emyrk/strava/api/modelsdk/sdktype"
)

// StringInt encodes IDs as strings to avoid precision loss in JavaScript.
// Keep this a named type for the TypeScript generator, and explicitly forward
// its methods because a defined type does not inherit sdktype.StringInt's methods.
type StringInt sdktype.StringInt

func (s StringInt) String() string {
	return sdktype.StringInt(s).String()
}

func (s StringInt) MarshalJSON() ([]byte, error) {
	return sdktype.StringInt(s).MarshalJSON()
}

func (s *StringInt) UnmarshalJSON(data []byte) error {
	return (*sdktype.StringInt)(s).UnmarshalJSON(data)
}
