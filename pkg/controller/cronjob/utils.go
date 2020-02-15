/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cronjob

import (
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"k8s.io/klog"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

// Utilities for dealing with Jobs and CronJobs and time.

func inActiveList(sj batchv1beta1.CronJob, uid types.UID) bool {
	for _, j := range sj.Status.Active {
		if j.UID == uid {
			return true
		}
	}
	return false
}

func deleteFromActiveList(sj *batchv1beta1.CronJob, uid types.UID) {
	if sj == nil {
		return
	}
	newActive := []v1.ObjectReference{}
	for _, j := range sj.Status.Active {
		if j.UID != uid {
			newActive = append(newActive, j)
		}
	}
	sj.Status.Active = newActive
}

// getParentUIDFromJob extracts UID of job's parent and whether it was found
func getParentUIDFromJob(j batchv1.Job) (types.UID, bool) {
	controllerRef := metav1.GetControllerOf(&j)

	if controllerRef == nil {
		return types.UID(""), false
	}

	if controllerRef.Kind != "CronJob" {
		klog.V(4).Infof("Job with non-CronJob parent, name %s namespace %s", j.Name, j.Namespace)
		return types.UID(""), false
	}

	return controllerRef.UID, true
}

// groupJobsByParent groups jobs into a map keyed by the job parent UID (e.g. scheduledJob).
// It has no receiver, to facilitate testing.
func groupJobsByParent(js []batchv1.Job) map[types.UID][]batchv1.Job {
	jobsBySj := make(map[types.UID][]batchv1.Job)
	for _, job := range js {
		parentUID, found := getParentUIDFromJob(job)
		if !found {
			klog.V(4).Infof("Unable to get parent uid from job %s in namespace %s", job.Name, job.Namespace)
			continue
		}
		jobsBySj[parentUID] = append(jobsBySj[parentUID], job)
	}
	return jobsBySj
}

func parseSchedule(scheduleString string) (cron.SpecSchedule, error) {
	tempSchedule, err := cron.ParseStandard(scheduleString)
	if err != nil {
		return cron.SpecSchedule{}, fmt.Errorf("unparseable schedule: %s : %s", scheduleString, err)
	}

	var schedule *cron.SpecSchedule
	if t, ok := tempSchedule.(*cron.SpecSchedule); ok && t != nil {
		schedule = t
	} else {
		return cron.SpecSchedule{}, errors.New("couldn't assert Schedule as SpecSchedule")
	}

	return *schedule, nil
}

// getRecentUnmetScheduleTimes gets a slice of times (from latest to oldest) that have passed when a Job should have started but did not.
// TODO(vllry): (latest, missedTimes, err) would be a better return type, but its easier to patch this way.
// A limited number of missed start times are returned.
// If there were missed times prior to the last known start time, then those are not returned.
func getRecentUnmetScheduleTimes(sj batchv1beta1.CronJob, now time.Time) ([]time.Time, error) {
	starts := []time.Time{}

	schedule, err := parseSchedule(sj.Spec.Schedule)
	if err != nil {
		return starts, err
	}

	var earliestTime time.Time
	if sj.Status.LastScheduleTime != nil {
		earliestTime = sj.Status.LastScheduleTime.Time
	} else {
		// If none found, then this is either a recently created scheduledJob,
		// or the active/completed info was somehow lost (contract for status
		// in kubernetes says it may need to be recreated), or that we have
		// started a job, but have not noticed it yet (distributed systems can
		// have arbitrary delays).  In any case, use the creation time of the
		// CronJob as last known start time.
		earliestTime = sj.ObjectMeta.CreationTimestamp.Time
	}
	if sj.Spec.StartingDeadlineSeconds != nil {
		// Controller is not going to schedule anything below this point
		schedulingDeadline := now.Add(-time.Second * time.Duration(*sj.Spec.StartingDeadlineSeconds))

		if schedulingDeadline.After(earliestTime) {
			earliestTime = schedulingDeadline
		}
	}
	if earliestTime.After(now) {
		return []time.Time{}, nil
	}

	// Step backward from now, until earliestTime.
	// Record all missed schedules between now and then.
	for t := PreviousScheduleTime(schedule, now); t.After(earliestTime); t = PreviousScheduleTime(schedule, t) {
		starts = append(starts, t)
		t = t.Add(-time.Second) // Make sure we don't return this same time again
		// We only care about knowing 0, 1, or multiple missed start times.
		if len(starts) > 1 {
			return starts, nil
		}
	}

	return starts, nil
}

// PreviousScheduleTime returns the most recent time this schedule is valid, before or at the given time.
// If no time can be found to satisfy the schedule, return the zero time.
func PreviousScheduleTime(schedule cron.SpecSchedule, t time.Time) time.Time {
	// General approach:
	// For Month, Day, Hour, Minute, Second:
	// Check if the time value matches.  If yes, continue to the next field.
	// If the field doesn't match the schedule, then decrement the field until it matches.
	// While decrementing the field, a wrap-around brings it back to the beginning
	// of the field list (since it is necessary to re-verify previous field
	// values)

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(-time.Duration(t.Nanosecond()) * time.Nanosecond)

	// If no time is found within five years, return zero.
	yearLimit := t.Year() - 5

WRAP:

	if t.Year() < yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	for 1<<uint(t.Month())&schedule.Month == 0 {
		// Step back into previous month.
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Add(-time.Second)

		// Wrap around the year.
		if t.Month() == time.December {
			goto WRAP
		}
	}

	// Now get a day in that month.
	for !dayMatches(schedule, t) {
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(-time.Second)
		date := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
		if t.Day() == date.Day() {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&schedule.Hour == 0 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location()).Add(-time.Second)

		if t.Hour() == 23 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&schedule.Minute == 0 {
		t = t.Truncate(time.Minute).Add(-time.Second)

		if t.Minute() == 59 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&schedule.Second == 0 {
		t = t.Add(-time.Second)

		if t.Second() == 59 {
			goto WRAP
		}
	}

	return t // If we've reached the end, we have a valid match (or none exist in range).
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func dayMatches(specSchedule cron.SpecSchedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&specSchedule.Dom > 0
		dowMatch bool = 1<<uint(t.Weekday())&specSchedule.Dow > 0
		starBit uint64 = 1 << 63
	)
	if specSchedule.Dom&starBit > 0 || specSchedule.Dow&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}

// getJobFromTemplate makes a Job from a CronJob
func getJobFromTemplate(sj *batchv1beta1.CronJob, scheduledTime time.Time) (*batchv1.Job, error) {
	labels := copyLabels(&sj.Spec.JobTemplate)
	annotations := copyAnnotations(&sj.Spec.JobTemplate)
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name := fmt.Sprintf("%s-%d", sj.Name, getTimeHash(scheduledTime))

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Annotations:     annotations,
			Name:            name,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(sj, controllerKind)},
		},
	}
	if err := legacyscheme.Scheme.Convert(&sj.Spec.JobTemplate.Spec, &job.Spec, nil); err != nil {
		return nil, fmt.Errorf("unable to convert job template: %v", err)
	}
	return job, nil
}

// getTimeHash returns Unix Epoch Time
func getTimeHash(scheduledTime time.Time) int64 {
	return scheduledTime.Unix()
}

func getFinishedStatus(j *batchv1.Job) (bool, batchv1.JobConditionType) {
	for _, c := range j.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
			return true, c.Type
		}
	}
	return false, ""
}

func IsJobFinished(j *batchv1.Job) bool {
	isFinished, _ := getFinishedStatus(j)
	return isFinished
}

// byJobStartTime sorts a list of jobs by start timestamp, using their names as a tie breaker.
type byJobStartTime []batchv1.Job

func (o byJobStartTime) Len() int      { return len(o) }
func (o byJobStartTime) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o byJobStartTime) Less(i, j int) bool {
	if o[i].Status.StartTime == nil && o[j].Status.StartTime != nil {
		return false
	}
	if o[i].Status.StartTime != nil && o[j].Status.StartTime == nil {
		return true
	}
	if o[i].Status.StartTime.Equal(o[j].Status.StartTime) {
		return o[i].Name < o[j].Name
	}
	return o[i].Status.StartTime.Before(o[j].Status.StartTime)
}
