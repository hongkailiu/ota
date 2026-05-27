// co-conditions-bugs extracts Jira ticket references from specific Go source files
// in the OpenShift origin repository and outputs their details as JSON or Markdown.
//
// It reads two hardcoded files:
//   - pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
//   - test/extended/machines/scale.go
//
// Usage:
//
//	co-conditions-bugs --origin-directory ~/repo/openshift/origin \
//	  --jira-username user@redhat.com \
//	  --jira-password-file ~/.config/ota/jira-api-token \
//	  --output-format md
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/jira"

	"github.com/petr-muller/ota/internal/flagutil"
)

type options struct {
	jira            flagutil.JiraOptions
	originDirectory string
	outputFormat    string
	outputFile      string
}

type ticketInfo struct {
	Key       string `json:"key"`
	Summary   string `json:"summary"`
	Status    string `json:"status"`
	Component string `json:"component"`
	URL       string `json:"url"`
}

func gatherOptions() options {
	var o options
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	o.jira.AddFlags(fs)
	fs.StringVar(&o.originDirectory, "origin-directory", "~/repo/openshift/origin", "Path to the origin repository directory")
	fs.StringVar(&o.outputFormat, "output-format", "json", "output format")
	fs.StringVar(&o.outputFile, "output-file", "-", "Output file path (use - for stdout)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.WithError(err).Fatalf("cannot parse args: '%s'", os.Args[1:])
	}

	return o
}

func (o *options) validate() error {
	if err := o.jira.Validate(); err != nil {
		return err
	}

	if o.outputFormat != "json" && o.outputFormat != "md" {
		return fmt.Errorf("invalid output format %q, must be 'json' or 'md'", o.outputFormat)
	}

	return nil
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

func fetchTicketInfo(jiraClient jira.Client, ticketID string) (*ticketInfo, error) {
	issue, err := jiraClient.GetIssue(ticketID)
	if err != nil {
		return nil, err
	}

	info := &ticketInfo{
		Key:     issue.Key,
		Summary: issue.Fields.Summary,
		Status:  issue.Fields.Status.Name,
		URL:     strings.Split(issue.Self, "/rest/api")[0] + "/browse/" + issue.Key,
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

	jiraClient, err := o.jira.Client()
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

	var output string

	switch o.outputFormat {
	case "json":
		jsonOutput, err := json.MarshalIndent(ticketInfos, "", "  ")
		if err != nil {
			logrus.WithError(err).Fatal("cannot marshal JSON output")
		}
		output = string(jsonOutput)

	case "md":
		var buf strings.Builder
		buf.WriteString("| Key | Summary | Status | Component |\n")
		buf.WriteString("|-----|---------|--------|-----------|\n")
		for _, ticket := range ticketInfos {
			escapedSummary := strings.ReplaceAll(ticket.Summary, "|", "\\|")
			escapedComponent := strings.ReplaceAll(ticket.Component, "|", "\\|")
			buf.WriteString(fmt.Sprintf("| [%s](%s) | %s | %s | %s |\n",
				ticket.Key,
				ticket.URL,
				escapedSummary,
				ticket.Status,
				escapedComponent))
		}
		output = buf.String()
	}

	if o.outputFile == "-" {
		fmt.Print(output)
	} else {
		if err := os.WriteFile(o.outputFile, []byte(output), 0644); err != nil {
			logrus.WithError(err).Fatalf("cannot write to file %s", o.outputFile)
		}
		logrus.Infof("Output written to %s", o.outputFile)
	}
}
