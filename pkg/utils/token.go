package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Id     string    `json:"id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
	Hash   []byte    `json:"-"`
}

// GenerateToken() generate a token for a specified ttl (time to life)
func (t *Token) GenerateToken(ttl time.Duration) (*Token, error) {
	token := &Token{Expiry: time.Now().Add(ttl), Id: ulid.Make().String()}

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	token.Token = base64.URLEncoding.EncodeToString(bytes)

	hashed, err := bcrypt.GenerateFromPassword([]byte(token.Token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	token.Hash = hashed

	return token, nil
}
