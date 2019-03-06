package filter

import (
	"common/logging"
	"net/http"
)

type AuthFilter struct {
	*BaseFilter
}

func (c *AuthFilter)before()  {
	logging.Debug("do auth filter before")
}

func (c *AuthFilter)do(w http.ResponseWriter,r *http.Request)   {
	logging.Debug("do auth filter")
}
