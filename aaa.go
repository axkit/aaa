package aaa

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/axkit/bitset"
	"github.com/axkit/vatel"
	"github.com/gbrlsnchs/jwt/v3"
)

var (
	// EPSignIn holds endpoint path to sign in.
	EPSignIn = "/auth/sign-in"

	// EPIsTokenValid holds endpoint path to is token valid.
	EPIsTokenValid = "/auth/is-token-valid"

	// EPRefreshToken holds endpoint path to refresh token.
	EPRefreshToken = "/auth/refresh-token"
)

// Config describes JWT configuration.
type Config struct {
	AccessTokenDuration       time.Duration
	RefreshTokenDuration      time.Duration
	IsRefreshNotBeforeEnabled bool
	Issuer                    string
	Subject                   string
	Audience                  []string
	EncryptionKey             string
}

// DefaultConfig holds default JWT configuration.
var DefaultConfig = Config{
	AccessTokenDuration:       time.Minute * 30,
	RefreshTokenDuration:      time.Hour * 24 * 30,
	IsRefreshNotBeforeEnabled: false,
	Issuer:                    "",
	Subject:                   "",
	Audience:                  []string{""},
	EncryptionKey:             "default",
}

// Userer is an interface what wraps access methods to User's attributes.
type Userer interface {
	UserID() int
	UserLogin() string
	UserRole() int
	UserLocked() bool
}

// UserStorer is an interface what wraps metods UserByCridentials and UserByID.
//
// UserByCredentials returns a user (object implementing interface Userer) if
// user with login and password is found.
//
// UserByID returns a user (object implementing interface Userer) identified by userID.
type UserStorer interface {
	UserByCredentials(login, password string) (Userer, error)
	UserByID(userID int) (Userer, error)
}

// RoleStorer is an interface what wraps methods IsRoleExist and RolePermissions.
//
// IsRoleExist returns true if role is roleID is exists.
//
// RolePermissions returns array of permissions and BitSet permission representation.
type RoleStorer interface {
	IsRoleExist(roleID int) bool
	RolePermissions(roleID int) ([]string, bitset.BitSet)
}

type AAA interface {

	// SignIn предоставляет метод для аутентификации пользователя.
	SignIn(login, password string) (*TokenSet, error)

	// ForceSignIn генерирует JWT токены для пользователя.
	// может использоваться для принудительной аутентификации пользователя, при
	// переходе по ссылки из письма активации адреса email.
	ForceSignIn(Userer) (*TokenSet, error)

	// routePerms []string принимает токен в виде base64 строки, проверяет на валидность.
	// и конвертирует в объект типа Tokener.
	//Authorize(encodedToken []byte, perms ...string) (*Token, error)
	//DecodeToken(encodedToken []byte) (*Token, error)

	// RefreshToken принимает токен в виде base64 строки, проверяет на валидность,
	// обновляет и возвращает новый токен.
	RefreshToken(encodedToken []byte) (*TokenSet, error)

	SetExtraAssigner(func(userID int) map[string]interface{})
}

// ApplicationPayload defines attributes what will be injected into JWT access token.
type ApplicationPayload struct {
	UserID           int                    `json:"user"`
	UserLogin        string                 `json:"login"`
	RoleID           int                    `json:"role"`
	PermissionBitSet json.RawMessage        `json:"perms,omitempty"`
	IsDebug          bool                   `json:"debug,omitempty"`
	ExtraPayload     map[string]interface{} `json:"extra,omitempty"`
}

// Token implements interface axkit/vatel Tokener.
type Token struct {
	jwt.Payload
	App ApplicationPayload `json:"app"`
}

// SystemPayload returns JWT system attributes related to standard.
func (t *Token) SystemPayload() map[string]interface{} {
	res := map[string]interface{}{}
	res["exp"] = t.Payload.ExpirationTime
	return res
}

func (t *Token) ApplicationPayload() vatel.TokenPayloader {
	return &t.App
}

//
type RefreshToken struct {
	jwt.Payload
	UserID int `json:"user"`
}

// TokenSet describes response on successfull sign in and refresh token requests.
type TokenSet struct {
	Access             string   `json:"access_token"`
	Refresh            string   `json:"refresh_token"`
	AllowedPermissions []string `json:"allowed_permissions"`
}

func (t *Token) JSON() []byte {
	return nil
}

func (t *ApplicationPayload) Login() string {
	return t.UserLogin
}

func (t *ApplicationPayload) User() int {
	return t.UserID
}

func (t *ApplicationPayload) Role() int {
	return t.RoleID
}

func (t *ApplicationPayload) Perms() []byte {
	if len(t.PermissionBitSet) < 2 || bytes.Equal(t.PermissionBitSet, []byte(`""`)) {
		return nil
	}

	// cut quotas on the sides.
	t.PermissionBitSet = t.PermissionBitSet[1 : len(t.PermissionBitSet)-1]
	return []byte(t.PermissionBitSet)
}

func (t *ApplicationPayload) Extra() interface{} {
	return t.ExtraPayload
}
