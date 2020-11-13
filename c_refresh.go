package aaa

import "github.com/axkit/vatel"

// RefreshController implements /refresh-token HTTP endpoint.
type RefreshController struct {
	a     *BasicAAA
	input struct {
		RefreshToken string `json:"refresh_token"`
	}

	output *TokenSet
}

// Input implements github.com/axkit/vatel Inputer interface.
func (a *RefreshController) Input() interface{} {
	return &a.input
}

// Result implements github.com/axkit/vatel Resulter interface.
func (a *RefreshController) Result() interface{} {
	return &a.output
}

// Handle implements github.com/axkit/vatel Handler interface.
func (a *RefreshController) Handle(ctx vatel.Context) error {

	tr, err := a.a.Refresh([]byte(a.input.RefreshToken))
	if err != nil {
		return err
	}
	a.output = tr
	return nil
}
