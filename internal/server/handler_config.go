package server

import (
	"net/http"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	configs = make(map[string]bool)
	cm      sync.Mutex
)

type configHandler struct{}

func AddConfigUpdate(agent string) {
	configs[agent] = true
}

func (configHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.Lock()
	defer cm.Unlock()

	targetAgent := path.Base(r.URL.String())

	hasNewConfigs, found := configs[targetAgent]
	if !found || !hasNewConfigs {
		_, _ = w.Write([]byte("OK"))

		return
	}

	log.Info().Str("agent", targetAgent).Msg("signal config changed")

	http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)

	delete(configs, targetAgent)
}
