package proxy

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// Proxy implements http.Handler accepting requests and relaying them
// to the defined upstream address server.
type Proxy struct {
	Upstream string
	Client   http.Client
}

func (t *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := t.Upstream + r.URL.String()

	slog.Debug("proxying request", "method", r.Method, "url", url)

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		slog.Error("failed constructing request",
			"method", r.Method, "url", url, "err", err.Error())
		respondError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed constructing request: %s", err.Error()))
		return
	}

	req.Header = r.Header

	resp, err := t.Client.Do(req)
	if err != nil {
		slog.Error("upstrea mrequest failed",
			"method", r.Method, "url", url, "err", err.Error())
		respondError(w, http.StatusInternalServerError,
			fmt.Sprintf("upstrea mrequest failed: %s", err.Error()))
		return
	}

	for key, values := range resp.Header {
		for _, val := range values {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		slog.Error("failed copying response body",
			"method", r.Method, "url", url, "err", err.Error())
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
