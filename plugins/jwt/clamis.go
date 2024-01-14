package jwt

import "github.com/golang-jwt/jwt/v4"

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserId       string `json:"user_id"`
	Username     string `json:"username"`
	NickName     string `json:"nick_name"`
	AuthorityIds []uint `json:"authority_ids"`
}

type AuthConf struct {
	AccessSecret string `json:"access_secret"`
	AccessExpire int64  `json:"access_expire"`
	BufferTime   int64  `json:"buffer_time"`
}
