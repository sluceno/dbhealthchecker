package dbhealthchecker

import (
	"database/sql"
	"time"
)

type TypeCheck int

const (
	Within TypeCheck = iota
	Higher
	Lower

	defaultRunSuiteEvery         = 24 * time.Hour
	defaultWaitTimeBetweenChecks = 5 * time.Minute
)

// HealthCheck type definition.
type DBHealthChecker struct {
	DB                    *sql.DB
	RunSuiteEvery         time.Duration
	WaitTimeBetweenChecks time.Duration

	healthChecks []HealthCheck
}

type DBHealthCheckerOption func(*DBHealthChecker)

func New(db *sql.DB, opts ...DBHealthCheckerOption) *DBHealthChecker {
	h := &DBHealthChecker{
		DB:                    db,
		RunSuiteEvery:         defaultRunSuiteEvery,
		WaitTimeBetweenChecks: defaultWaitTimeBetweenChecks,
	}

	for _, option := range opts {
		option(h)
	}

	return h
}

func SetWaitTimeBetweenChecks(waitTime time.Duration) DBHealthCheckerOption {
	return func(h *DBHealthChecker) {
		h.WaitTimeBetweenChecks = waitTime
	}
}

func SetRunSuiteEvery(waitTime time.Duration) DBHealthCheckerOption {
	return func(h *DBHealthChecker) {
		h.RunSuiteEvery = waitTime
	}
}

func (h *DBHealthChecker) Add(healthChecks ...HealthCheck) {
	h.healthChecks = append(h.healthChecks, healthChecks...)
}

func (h DBHealthChecker) Run() <-chan HealthCheck {
	out := make(chan HealthCheck)

	ticker := time.NewTicker(h.RunSuiteEvery)

	go func() {
		h.runHealthChecksSuite(out)
		for _ = range ticker.C {
			h.runHealthChecksSuite(out)
		}
	}()

	return out
}

func (h *DBHealthChecker) runHealthChecksSuite(out chan HealthCheck) {
	for _, healthCheck := range h.healthChecks {
		healthCheck.count, healthCheck.err = query(h.DB, healthCheck)

		out <- healthCheck

		time.Sleep(h.WaitTimeBetweenChecks)
	}
}

func query(db *sql.DB, check HealthCheck) (int, error) {
	var count int

	err := db.QueryRow(check.Query).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		return 0, err
	case err != nil:
		return 0, err
	default:
		return count, nil
	}

	return 0, err
}
