package cronjob

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	cronjobSubsystem = "cronjob_controller"
	cronNameKey      = "cronjob"
	namespaceKey     = "namespace"
	skipReasonKey    = "reason"
)

const (
	skipReasonConcurrencyPolicy = "concurrencyPolicy"
	skipReasonMissedDeadline    = "missedDeadline"
	skipReasonError             = "error"
)

var schedulingDecisionInvoke = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: cronjobSubsystem,
		Name:      "scheduling_decision_invoke",
		Help:      "Counter that increments when the cronjob controller decides to invoke a CronJob",
	},
	[]string{namespaceKey, cronNameKey})

var schedulingDecisionSkip = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: cronjobSubsystem,
		Name:      "scheduling_decision_skip",
		Help:      "Counter that increments when the cronjob controller decides to skip a CronJob invocation",
	},
	[]string{namespaceKey, cronNameKey, skipReasonKey})

var jobSucceeded = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: cronjobSubsystem,
		Name:      "job_succeeded",
		Help:      "Counter that increments when the cronjob controller detects a child Job has completed with success",
	}, []string{namespaceKey, cronNameKey})

var jobFailed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: cronjobSubsystem,
		Name:      "job_failed",
		Help:      "Counter that increments when the cronjob controller detects a child Job has completed with failure",
	}, []string{namespaceKey, cronNameKey})

var registerOnce sync.Once

func registerMetrics() {
	registerOnce.Do(func() {
		prometheus.MustRegister(schedulingDecisionInvoke)
		prometheus.MustRegister(schedulingDecisionSkip)
		prometheus.MustRegister(jobSucceeded)
		prometheus.MustRegister(jobFailed)
	})
}
