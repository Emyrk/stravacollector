package strava

import (
	"fmt"
	"strings"
)

type Error struct {
	Message string          `json:"message"`
	Errors  []DetailedError `json:"errors"`
}

func (e Error) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf(":%s-%s-%s", err.Code, err.Resource, err.Field))
	}
	return b.String()
}

func (e Error) ContainsCode(code string) bool {
	for _, err := range e.Errors {
		if err.Code == code {
			return true
		}
	}
	return false
}

type DetailedError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}
