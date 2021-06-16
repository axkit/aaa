package aaa

import (
	"github.com/axkit/errors"
	"github.com/axkit/vatel"
)

// SignInController implements sign in HTTP endpoint.
type SignInController struct {
	a     *BasicAAA
	input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	output *TokenSet
}

// Input returns reference to incoming struct.
func (c *SignInController) Input() interface{} {
	return &c.input
}

// Result returns reference to sucessfull output.
func (c *SignInController) Result() interface{} {
	return c.output
}

// Handle implements github.com/axkit/vatel Handler interface.
func (c *SignInController) Handle(ctx vatel.Context) error {
	tr, err := c.a.SignIn(c.input.Login, c.input.Password) // , string(ctx.UserAgent()), ctx.RemoteIP().String())
	if err != nil {
		return errors.Catch(err).StatusCode(401).Severity(errors.Medium).Msg("sign in failed")
	}
	c.output = tr
	return nil
}
