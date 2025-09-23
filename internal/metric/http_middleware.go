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
		routeCtx := chi.RouteContext(r.Context()) // chi 라우트 컨텍스트 가져오기
		path := routeCtx.RoutePattern()           // 매칭된 라우트 패턴 추출 (예: /users/{id})

		HTTPInFlightRequests.Inc()       // 현재 처리 중인 요청 수 증가
		defer HTTPInFlightRequests.Dec() // 함수 종료 시 처리 중인 요청 수 감소

		start := time.Now()                                              // 요청 시작 시간 기록
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // 상태 코드 기록용 ResponseWriter 래퍼 생성 (기본 200)

		next.ServeHTTP(rec, r) // 실제 핸들러 실행 (응답 기록은 rec에 저장됨)

		duration := time.Since(start).Seconds() // 요청 처리 시간 계산 (초 단위)
		statusStr := strconv.Itoa(rec.status)   // 상태 코드를 문자열로 변환

		HTTPRequestsTotal.WithLabelValues(r.Method, path, statusStr).Inc()               // 요청 수 카운터 증가
		HTTPRequestDuration.WithLabelValues(r.Method, path, statusStr).Observe(duration) // 요청 처리 시간 기록 (히스토그램에)
	})
}
