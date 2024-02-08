package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int
	Username string
	Tests    []*Test
}

func NewUser(id int, username string) *User {
	return &User{
		ID:       id,
		Username: username,
		Tests:    []*Test{},
	}
}

func (u *User) AddSignedTest(questions, answers string) string {
	uuid := uuid.New()
	test := &Test{
		Questions: questions,
		Answers:   answers,
		Signature: uuid.String(),
		SignedAt:  time.Now(),
	}
	u.Tests = append(u.Tests, test)
	return uuid.String()
}

func (u *User) VerifySignature(signature string) (bool, string, time.Time) {
	for _, t := range u.Tests {
		if t.Signature == signature {
			return true, t.Answers, t.SignedAt
		}
	}
	return false, "", time.Time{}
}

func (u *User) String() string {
	return fmt.Sprintf("User %d: %s Tests: %d)", u.ID, u.Username, len(u.Tests))
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
