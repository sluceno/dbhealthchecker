package dbhealthchecker

import "fmt"

type HealthCheck struct {
	Name          string
	Query         string
	ConditionType TypeCheck
	Max           int
	Min           int
	Thereshold    int

	count int
	err   error
}

func (h HealthCheck) Count() int {
	return h.count
}

func (h HealthCheck) Error() error {
	return h.err
}

func (h HealthCheck) Healthy() bool {
	switch h.ConditionType {
	case Within:
		return h.healthCheckWithin()
	case Higher:
		return h.healthCheckHigher()
	case Lower:
		return h.healthCheckLower()
	}

	return false
}

func (h HealthCheck) String() string {
	switch h.ConditionType {
	case Within:
		return fmt.Sprintf("Found %d items. Expected to be between %d and %d.", h.count, h.Min, h.Max)
	case Higher:
		return fmt.Sprintf("Found %d items. Expected to be higher than %d.", h.count, h.Thereshold)
	case Lower:
		return fmt.Sprintf("Found %d items. Expected to be lower than %d.", h.count, h.Thereshold)
	}

	return ""
}

func (h HealthCheck) healthCheckWithin() bool {
	if h.count < h.Min || h.count > h.Max {
		return false
	}

	return true
}

func (h HealthCheck) healthCheckHigher() bool {
	if h.count <= h.Thereshold {
		return false
	}

	return true
}

func (h HealthCheck) healthCheckLower() bool {
	if h.count > h.Thereshold {
		return false
	}

	return true
}
