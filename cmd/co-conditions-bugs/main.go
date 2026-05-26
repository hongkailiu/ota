package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/petr-muller/ota/internal/flagutil"
)

type options struct {
	jira            flagutil.JiraOptions
	originDirectory string
}

func gatherOptions() options {
	var o options
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	o.jira.AddFlags(fs)
	fs.StringVar(&o.originDirectory, "origin-directory", "~/repo/openshift/origin", "Path to the origin repository directory")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.WithError(err).Fatalf("cannot parse args: '%s'", os.Args[1:])
	}

	return o
}

func (o *options) validate() error {
	return o.jira.Validate()
}

func main() {
	o := gatherOptions()
	if err := o.validate(); err != nil {
		logrus.WithError(err).Fatal("invalid options")
	}

	logrus.Info("co-conditions-bugs started")
}
