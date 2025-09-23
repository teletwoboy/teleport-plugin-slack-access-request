package metric

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// statusRecorder wraps http.ResponseWriter to capture status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// InstrumentHTTP is Middleware for chi
func InstrumentHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HTTPInFlightRequests.Inc()       // 현재 처리 중인 요청 수 증가
		defer HTTPInFlightRequests.Dec() // 함수 종료 시 감소

		start := time.Now()                                              // 요청 시작 시간
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // 상태 코드 기록용 래퍼

		next.ServeHTTP(rec, r) // 실제 핸들러 실행 (여기서 라우팅 최종 결정)

		// 실행 이후에 최종 패턴을 읽는다
		routeCtx := chi.RouteContext(r.Context()) // chi 라우트 컨텍스트
		path := routeCtx.RoutePattern()           // 최종 매칭된 라우트 패턴 (예: /api/v1/users/{id})
		if path == "" {
			path = "unknown" // 매칭 실패시 카디널리티 안전한 기본값
		}

		duration := time.Since(start).Seconds() // 처리 시간(초)
		statusStr := strconv.Itoa(rec.status)   // 상태코드 문자열

		HTTPRequestsTotal.WithLabelValues(r.Method, path, statusStr).Inc()               // 요청 수 증가
		HTTPRequestDuration.WithLabelValues(r.Method, path, statusStr).Observe(duration) // 지연시간 관측
	})
}
