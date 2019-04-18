package filter

import (
	"common/logging"
	"net/http"
)

type AuthFilter struct {
	*BaseApiFilter
}

func (c *AuthFilter) before() {
	logging.Debug("do auth filter before")
}

func (c *AuthFilter) do(w http.ResponseWriter, r *http.Request) {
	logging.Debug("do auth filter")
}

type RegFilter struct {
	*BaseWsFilter
}

func (c *RegFilter) do(msg []byte) bool {
	return true
}
