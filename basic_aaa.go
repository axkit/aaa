package aaa

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axkit/bitset"
	"github.com/axkit/errors"
	"github.com/axkit/vatel"
	"github.com/gbrlsnchs/jwt/v3"
)

// BasicAAA holds data required for implementation AAA interface and axkit/vatel interfaces
// Authorizer, TokenDecoder.
type BasicAAA struct {
	cfg           Config
	us            UserStorer
	rs            RoleStorer
	extraAssigner func(int) (map[string]interface{}, error)
}

// New returns default implementation of AAA based on JWT.
func New(cfg Config, u UserStorer, r RoleStorer) *BasicAAA {
	return &BasicAAA{cfg: cfg, us: u, rs: r}
}

// Init ...
func (a *BasicAAA) Init(ctx context.Context) error {
	return nil
}

// Start ...
func (a *BasicAAA) Start(ctx context.Context) error {
	return nil
}

// SetExtraAssigner receives a funcion what will be called in /sign-in and /refresh-token
// endpoints. Data returned by the function will be assigned to JWT payload attribute "app->extra".
func (a *BasicAAA) SetExtraAssigner(f func(userID int) (map[string]interface{}, error)) {
	a.extraAssigner = f
}

// Endpoints returnts lists of endpoints serving by BasicAAA.
func (a *BasicAAA) Endpoints() []vatel.Endpoint {
	return []vatel.Endpoint{
		{Method: "POST", Path: EPSignIn, NoInputLog: true,
			Controller: func() vatel.Handler { return &SignInController{a: a} }},
		{Method: "POST", Path: EPRefreshToken, NoInputLog: true,
			Controller: func() vatel.Handler { return &RefreshController{a: a} }},
		{Method: "POST", Path: EPIsTokenValid,
			Controller: func() vatel.Handler { return &IsTokenValidController{a: a} }},
	}
}

// SignIn implements sign in logic. In case of succesfull result returns
func (a *BasicAAA) SignIn(login, password string) (*TokenSet, error) {
	if len(login) == 0 {
		return nil, errors.ValidationFailed("empty login")
	}
	if len(password) == 0 {
		return nil, errors.ValidationFailed("empty password").Set("login", login)
	}

	u, err := a.us.UserByCredentials(login, password)
	if err != nil {
		// записать в лог err потому что он подробный (логин не найден, пароль не соответсвует)
		return nil, errors.Catch(err).StatusCode(401).Set("login", login).Msg("invalid cridentials")
	}

	if u.UserLocked() {
		return nil, errors.NewMedium("user is locked").StatusCode(401).Protect().Set("login", login)
	}

	return a.generateTokenSet(u)
}

// GenerateToken generates JWT token without cridentials.
func (a *BasicAAA) GenerateToken(u Userer) (*TokenSet, error) {
	return a.generateTokenSet(u)
}

// Refresh refreshes JWT token.
func (a *BasicAAA) Refresh(encodedToken []byte) (*TokenSet, error) {

	var rt RefreshToken
	_, err := jwt.Verify(encodedToken, jwt.NewHS256([]byte(a.cfg.EncryptionKey)), &rt)
	if err != nil {
		return nil, errors.Catch(err).StatusCode(401).Msg("refresh token verify failed")
	}

	var tv = jwt.ExpirationTimeValidator(time.Now())

	err = tv(&rt.Payload)
	if err != nil {
		return nil, errors.Catch(err).StatusCode(401).Msg("refresh token expired")
	}

	u, err := a.us.UserByID(rt.UserID)
	if err != nil {
		return nil, errors.Catch(err).StatusCode(401)
	}

	return a.generateTokenSet(u)
}

// Decode
func (s *BasicAAA) Decode(encodedToken []byte) (vatel.Tokener, error) {

	if len(encodedToken) == 0 {
		return nil, errors.New("empty access token").StatusCode(401)
	}

	var at Token
	_, err := jwt.Verify(encodedToken, jwt.NewHS256([]byte(s.cfg.EncryptionKey)), &at)
	if err != nil {
		return nil, errors.Catch(err).StatusCode(401).Msg("invalid access token")
	}

	var tv = jwt.ExpirationTimeValidator(time.Now())

	err = tv(&at.Payload)
	if err != nil {
		return nil, errors.Catch(err).StatusCode(401).Msg("access token expired")
	}

	return &at, nil
}

// IsAllowed implements interface axkit/vatel Autorizer. Method receives perms from JTW token
// and endpointPemrs. Return true if all endpointPerms are inside requestPerms.
func (a *BasicAAA) IsAllowed(requestPerms []byte, bitpos ...uint) (bool, error) {
	if n := len(requestPerms); n > 0 && n%2 != 0 {
		return false, errors.New("invalid request permission set").StatusCode(401)
	}
	return bitset.AreSet(requestPerms, bitpos...)
}

func (a *BasicAAA) generateTokenSet(u Userer) (*TokenSet, error) {

	var (
		at  Token
		rt  RefreshToken
		res TokenSet
	)

	roleID := u.UserRole()
	if ok := a.rs.IsRoleExist(roleID); !ok {
		return nil, errors.New("getting user's role failed").StatusCode(500).Protect()
	}

	userID := u.UserID()

	now := time.Now()
	at.Payload = jwt.Payload{
		Issuer:         a.cfg.Issuer,
		Subject:        a.cfg.Subject,
		Audience:       jwt.Audience(a.cfg.Audience),
		ExpirationTime: jwt.NumericDate(now.Add(a.cfg.AccessTokenDuration)),
		//NotBefore:      jwt.NumericDate(now), // lets not use it know
		IssuedAt: jwt.NumericDate(now),
		JWTID:    "test",
	}

	at.App.UserID = userID
	at.App.RoleID = roleID
	at.App.UserLogin = u.UserLogin()

	perms, bs := a.rs.RolePermissions(roleID)
	res.AllowedPermissions = perms
	at.App.PermissionBitSet = json.RawMessage("\"" + bs.String() + "\"")

	if a.extraAssigner != nil {
		var err error
		at.App.ExtraPayload, err = a.extraAssigner(userID)
		if err != nil {
			return nil, errors.Catch(err).StatusCode(500).Msg("building token->extra failed")
		}
	}

	rt.Payload = at.Payload
	rt.ExpirationTime = jwt.NumericDate(now.Add(a.cfg.RefreshTokenDuration))
	if a.cfg.IsRefreshNotBeforeEnabled {
		rt.NotBefore = jwt.NumericDate(now.Add(a.cfg.AccessTokenDuration))
	}

	rt.UserID = userID

	buf, err := jwt.Sign(at, jwt.NewHS256([]byte(a.cfg.EncryptionKey)))
	if err != nil {
		return nil, errors.Catch(err).StatusCode(500).Msg("access token generation failed")
	}
	res.Access = string(buf)

	buf, err = jwt.Sign(rt, jwt.NewHS256([]byte(a.cfg.EncryptionKey)))
	if err != nil {
		return nil, errors.Catch(err).StatusCode(500).Msg("refresh token generation failed")
	}
	res.Refresh = string(buf)

	return &res, nil
}
