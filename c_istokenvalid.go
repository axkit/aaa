package aaa

import "github.com/axkit/vatel"

// IsTokenValidController implements /is-token-valid HTTP endpoint.
type IsTokenValidController struct {
	a      *BasicAAA
	output struct {
		Result string `json:"result"`
	}
}

// Result implements github.com/axkit/vatel Resulter interface.
func (c *IsTokenValidController) Result() interface{} {
	return &c.output
}

// Handle implements github.com/axkit/vatel Handler interface.
func (c *IsTokenValidController) Handle(ctx vatel.Context) error {
	c.output.Result = "ok"
	return nil
}
