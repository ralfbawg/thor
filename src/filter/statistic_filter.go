package filter

import (
	"common/logging"
	"net/http"
)

type StatisticFilter struct {
	*BaseFilter
}

func (c *StatisticFilter) before() {
	logging.Debug("do statistic filter before")
}

func (c *StatisticFilter) do(w http.ResponseWriter, r *http.Request) {
	logging.Debug("do statistic filter")
}
