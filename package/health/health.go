package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

type Checker interface {
	Check(ctx context.Context) error
	Name() string
}

type Health struct {
	mu      sync.RWMutex
	status  Status
	checks  map[string]Checker
	started time.Time
	version string
	service string
}

func New(service, version string) *Health {
	return &Health{
		status:  StatusHealthy,
		checks:  make(map[string]Checker),
		started: time.Now(),
		version: version,
		service: service,
	}
}

func (h *Health) Register(name string, checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = checker
}

func (h *Health) Check(ctx context.Context) map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]interface{})
	result["service"] = h.service
	result["version"] = h.version
	result["status"] = string(h.status)
	result["uptime"] = time.Since(h.started).String()
	result["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	checks := make(map[string]string)
	allHealthy := true

	for name, checker := range h.checks {
		checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := checker.Check(checkCtx)
		cancel()

		if err != nil {
			checks[name] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks[name] = "healthy"
		}
	}

	result["checks"] = checks

	if !allHealthy {
		result["status"] = string(StatusUnhealthy)
	}

	return result
}

func (h *Health) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		result := h.Check(ctx)

		status := http.StatusOK
		if result["status"] == string(StatusUnhealthy) {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(result)
	}
}

func (h *Health) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		result := h.Check(ctx)
		status := result["status"].(string)

		if status == string(StatusHealthy) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
		}
	}
}

func (h *Health) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
	}
}

type DBChecker struct {
	db *sql.DB
}

func NewDBChecker(db *sql.DB) *DBChecker {
	return &DBChecker{db: db}
}

func (c *DBChecker) Name() string {
	return "database"
}

func (c *DBChecker) Check(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
