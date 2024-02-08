package entity

import (
	"encoding/json"
	"time"
)

type Test struct {
	Questions string
	Answers   string
	Signature string
	SignedAt  time.Time
}

func (u *Test) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Test) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
