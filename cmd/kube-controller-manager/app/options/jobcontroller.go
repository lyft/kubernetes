/*
Copyright 2018 The Kubernetes Authors.

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

package options

import (
	"github.com/spf13/pflag"

	kubectrlmgrconfig "k8s.io/kubernetes/pkg/controller/apis/config"
)

// JobControllerOptions holds the JobController options.
type JobControllerOptions struct {
	*kubectrlmgrconfig.JobControllerConfiguration
}

// AddFlags adds flags related to JobController for controller manager to the specified FlagSet.
func (o *JobControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up JobController config with options.
func (o *JobControllerOptions) ApplyTo(cfg *kubectrlmgrconfig.JobControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentJobSyncs = o.ConcurrentJobSyncs

	return nil
}

// Validate checks validation of JobControllerOptions.
func (o *JobControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}

// CronJobControllerOptions holds the CronJobController options.
type CronJobControllerOptions struct {
	*kubectrlmgrconfig.CronJobControllerConfiguration
}

// AddFlags adds flags related to CronJobController for controller manager to the specified FlagSet.
func (o *CronJobControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.DurationVar(&o.CronJobControllerSyncPeriod.Duration, "cronjob-controller-sync-period", o.CronJobControllerSyncPeriod.Duration, "Period for syncing cronjobs.")
}

// ApplyTo fills up CronJobController config with options.
func (o *CronJobControllerOptions) ApplyTo(cfg *kubectrlmgrconfig.CronJobControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.CronJobControllerSyncPeriod = o.CronJobControllerSyncPeriod

	return nil
}

// Validate checks validation of CronJobControllerOptions.
func (o *CronJobControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
