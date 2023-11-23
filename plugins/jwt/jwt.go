package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

type JWT struct {
	SecretKey []byte
}

func NewJwt(secretKey string) *JWT {
	return &JWT{
		[]byte(secretKey),
	}
}

// CreateClaim return the CustomClaims
func (j *JWT) CreateClaims(bufferTime int64, expiresAt int64, baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: bufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),                                  // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresAt) * time.Second)), // 过期时间 7天  配置文件
		},
	}
	return claims
}

// NewJwtToken returns the jwt token from the given data.
func (j *JWT) CreateToken(secretKey string, claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Claims = claims
	return token.SignedString(j.SecretKey)
}

// parse token return the customClaims struct
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SecretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// 解析token获取用户userId
func (j *JWT) GetUserID(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.BaseClaims.UserId, nil
}

// 解析token获取用户角色
func (j *JWT) GetUserAuthorityId(tokenString string) ([]uint, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims.BaseClaims.AuthorityIds, nil
}
