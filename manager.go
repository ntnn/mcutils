package mcutils

import (
	mctrl "sigs.k8s.io/multicluster-runtime"
)

// SilentManagerOpts modifies and returns a set of manager.Options that
// disables leader election and any network listeners.
// Useful e.g. for testing.
func SilentManagerOpts(opts mctrl.Options) mctrl.Options {
	opts.LeaderElection = false
	opts.Metrics.BindAddress = "0"
	opts.HealthProbeBindAddress = "0"
	opts.PprofBindAddress = "0"
	return opts
}
