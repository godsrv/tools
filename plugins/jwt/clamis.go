package jwt

import "github.com/golang-jwt/jwt/v4"

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserId       string
	Username     string
	NickName     string
	AuthorityIds []uint
}

type Auth struct {
	AccessSecret string
	AccessExpire int64
	BufferTime   int64
}
