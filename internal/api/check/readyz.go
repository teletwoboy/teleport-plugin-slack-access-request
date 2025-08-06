package check

import (
	"net/http"
	"sync/atomic"
)

func Readyz(w http.ResponseWriter, _ *http.Request, isReady *atomic.Value) {
	if isReady == nil || !isReady.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
