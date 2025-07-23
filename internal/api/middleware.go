package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"
	"teleport-plugin-slack-access-request/internal/config"
	"time"
)

func VerifySlackRequest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slackSig := r.Header.Get("X-Slack-Signature")
			slackTimestamp := r.Header.Get("X-Slack-Request-Timestamp")
			if slackSig == "" || slackTimestamp == "" {
				http.Error(w, "Unauthorized - missing headers", http.StatusBadRequest)
				return
			}

			ts, err := strconv.ParseInt(slackTimestamp, 10, 64)
			if err != nil || abs(time.Now().Unix()-ts) > 60*5 {
				http.Error(w, "Unauthorized - timestamp expired", http.StatusUnauthorized)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unauthorized - read body failed", http.StatusUnauthorized)
				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			base := "v0" + ":" + slackTimestamp + ":" + string(bodyBytes)
			mac := hmac.New(sha256.New, []byte(config.Cfg.Slack.SigningSecret))
			mac.Write([]byte(base))
			expectedSig := "v0=" + hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(slackSig), []byte(expectedSig)) {
				http.Error(w, "Unauthorized - signature mismatch", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
