package main

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/andygrunwald/go-jira"
	"github.com/sirupsen/logrus"
	prowflagutil "sigs.k8s.io/prow/pkg/flagutil"

	"github.com/petr-muller/ota/internal/flagutil"
)

type options struct {
	jira            flagutil.JiraOptions
	originDirectory string
}

type ticketInfo struct {
	Key       string `json:"key"`
	Summary   string `json:"summary"`
	Status    string `json:"status"`
	Component string `json:"component"`
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

func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if len(path) == 1 {
		return homeDir, nil
	}

	return filepath.Join(homeDir, path[1:]), nil
}

func readSourceFiles(baseDir string) ([]byte, error) {
	file1Path := filepath.Join(baseDir, "pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go")
	file2Path := filepath.Join(baseDir, "test/extended/machines/scale.go")

	content1, err := os.ReadFile(file1Path)
	if err != nil {
		return nil, err
	}

	content2, err := os.ReadFile(file2Path)
	if err != nil {
		return nil, err
	}

	// Combine both files with newline separator
	combined := append(content1, '\n')
	combined = append(combined, content2...)

	return combined, nil
}

func parseJiraTickets(content []byte) []string {
	re := regexp.MustCompile(`https://.*/browse/(OCPBUGS-\d+)`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	// Deduplicate using map
	seen := make(map[string]bool)
	var tickets []string

	for _, match := range matches {
		if len(match) > 1 {
			ticketID := match[1]
			if !seen[ticketID] {
				seen[ticketID] = true
				tickets = append(tickets, ticketID)
			}
		}
	}

	return tickets
}

func getFieldValue(opts *prowflagutil.JiraOptions, fieldName string) string {
	val := reflect.ValueOf(opts).Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}

	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Ptr:
		if field.IsNil() {
			return ""
		}
		return field.Elem().String()
	}
	return ""
}

func fetchTicketInfo(jiraClient *jira.Client, ticketID string) (*ticketInfo, error) {
	issue, _, err := jiraClient.Issue.Get(ticketID, nil)
	if err != nil {
		return nil, err
	}

	info := &ticketInfo{
		Key:     issue.Key,
		Summary: issue.Fields.Summary,
		Status:  issue.Fields.Status.Name,
	}

	if len(issue.Fields.Components) > 0 {
		info.Component = issue.Fields.Components[0].Name
	}

	return info, nil
}

func main() {
	o := gatherOptions()
	if err := o.validate(); err != nil {
		logrus.WithError(err).Fatal("invalid options")
	}

	expandedDir, err := expandPath(o.originDirectory)
	if err != nil {
		logrus.WithError(err).Fatal("cannot expand origin directory path")
	}

	logrus.Infof("Using origin directory: %s", expandedDir)

	content, err := readSourceFiles(expandedDir)
	if err != nil {
		logrus.WithError(err).Fatal("cannot read source files")
	}

	logrus.Infof("Read %d bytes from source files", len(content))

	tickets := parseJiraTickets(content)
	logrus.Infof("Found %d unique Jira tickets", len(tickets))

	if len(tickets) == 0 {
		logrus.Info("No Jira tickets found")
		return
	}

	// Create go-jira client using basic auth
	// First validate the options to ensure we have required credentials
	_, err = o.jira.Client()
	if err != nil {
		logrus.WithError(err).Fatal("cannot validate Jira client options")
	}

	// Use reflection to extract credentials from prow's JiraOptions
	var username, password, endpoint string
	jiraOpts := &o.jira.JiraOptions

	// Extract endpoint
	endpointField := getFieldValue(jiraOpts, "endpoint")
	if endpointField != "" {
		endpoint = endpointField
	} else {
		endpoint = "https://redhat.atlassian.net"
	}

	// Extract username
	username = getFieldValue(jiraOpts, "username")

	// Extract password from file if specified
	passwordFileRef := getFieldValue(jiraOpts, "passwordFile")
	if passwordFileRef != "" {
		passwordBytes, err := os.ReadFile(passwordFileRef)
		if err != nil {
			logrus.WithError(err).Fatal("cannot read password file")
		}
		password = string(passwordBytes)
	}

	// Create basic auth transport
	tp := &jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	// Create go-jira client
	jiraClient, err := jira.NewClient(tp.Client(), endpoint)
	if err != nil {
		logrus.WithError(err).Fatal("cannot create Jira client")
	}

	var ticketInfos []*ticketInfo

	for _, ticketID := range tickets {
		logrus.Infof("Fetching ticket: %s", ticketID)
		info, err := fetchTicketInfo(jiraClient, ticketID)
		if err != nil {
			logrus.WithError(err).Fatalf("cannot fetch ticket %s", ticketID)
		}
		ticketInfos = append(ticketInfos, info)
	}

	logrus.Infof("Successfully fetched %d tickets", len(ticketInfos))
}
