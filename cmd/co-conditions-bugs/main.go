// co-conditions-bugs extracts Jira ticket references from specific Go source files
// in the OpenShift origin repository and outputs their details as JSON or Markdown.
//
// It reads two hardcoded files:
//   - pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
//   - test/extended/machines/scale.go
//
// Usage:
//
//	co-conditions-bugs --origin-directory ~/repo/openshift/origin --output-format md
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	andyjira "github.com/andygrunwald/go-jira"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
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
	Number           int    `json:"count"`
	Key              string `json:"key"`
	Summary          string `json:"summary"`
	Status           string `json:"status"`
	Component        string `json:"component"`
	IsPixaaComponent bool   `json:"is_pixaa_component"`
	Resolution       string `json:"resolution"`
	Assignee         string `json:"assignee"`
	URL              string `json:"url"`
	Notes            string `json:"notes"`
	TargetVersion    string `json:"target_version"`
	ReleaseBlocker   string `json:"release_blocker"`
	ParentID         string `json:"parent_id,omitempty"`
	ParentKey        string `json:"parent_key,omitempty"`
}

func gatherOptions() options {
	var o options
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	o.jira.AddFlags(fs)
	fs.StringVar(&o.originDirectory, "origin-directory", "~/repo/openshift/origin", "Path to the origin repository directory")
	fs.StringVar(&o.outputFormat, "output-format", "json", "output format. Either json or md")
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

func extractCustomField(unknowns map[string]interface{}, fieldID string) string {
	value, exists := unknowns[fieldID]
	if !exists || value == nil {
		return ""
	}

	// Handle direct string
	if str, ok := value.(string); ok {
		return str
	}

	// Handle object with Name field (common pattern)
	if obj, ok := value.(map[string]interface{}); ok {
		if name, ok := obj["name"].(string); ok {
			return name
		}
		if val, ok := obj["value"].(string); ok {
			return val
		}
	}

	// Handle array with objects (for Target Version field)
	if arr, ok := value.([]interface{}); ok {
		if len(arr) > 0 {
			if obj, ok := arr[0].(map[string]interface{}); ok {
				if name, ok := obj["name"].(string); ok {
					return name
				}
			}
		}
	}

	return ""
}

var (
	notes = map[string]string{
		"OCPBUGS-22382": "Won't Do confirmed",
		"OCPBUGS-23744": "Won't Do confirmed: OLMv0 in maintenance mode",
		"OCPBUGS-65583": "Won't Do confirmed: OLMv0 in maintenance mode",
		"OCPBUGS-20056": "Possibly dup of  OCPBUGS-66027; Descoped",
		"OCPBUGS-42837": "To be removed in 5.1",
		"OCPBUGS-65984": "Two-Nodes clusters not fixed; Descoped",
		"OCPBUGS-64852": "Under evaluation",
		"OCPBUGS-86308": "Descoped",
		"OCPBUGS-38661": "Descoped",
		"OCPBUGS-23746": "Descoped",
		"OCPBUGS-62633": "Descoped",
		"OCPBUGS-66213": "Descoped",
		"OCPBUGS-38662": "Descoped",
		"OCPBUGS-65896": "Descoped",
		"OCPBUGS-63116": "Descoped",
		"OCPBUGS-38678": "Descoped", // not an exception in o/origin
		"OCPBUGS-38663": "Descoped",
		"OCPBUGS-66027": "Descoped", // not an exception in o/origin
		"OCPBUGS-62629": "Descoped",
		"OCPBUGS-66225": "Under evaluation",
		"OCPBUGS-65647": "Won't Do confirmed",
	}

	pixaaComponents = sets.New[string](
		"Cloud Compute",
		"Cluster Autoscaler", // ?
		"Management Console",
		"Cloud Compute", // ?
		"OLM",
	)
)

const (
	customFieldTargetVersion  = "customfield_10855"
	customFieldReleaseBlocker = "customfield_10847"
)

func isPixaaComponent(component string) bool {
	parts := strings.Split(component, "/")
	if len(parts) > 0 {
		firstPart := strings.TrimSpace(parts[0])
		return pixaaComponents.Has(firstPart)
	}
	return false
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
		Notes:   notes[issue.Key],
	}

	if len(issue.Fields.Components) > 0 {
		info.Component = issue.Fields.Components[0].Name
		info.IsPixaaComponent = isPixaaComponent(info.Component)
	}

	if issue.Fields.Resolution != nil {
		info.Resolution = issue.Fields.Resolution.Name
	}

	if issue.Fields.Assignee != nil {
		info.Assignee = issue.Fields.Assignee.DisplayName
	}

	// Extract parent information
	if issue.Fields.Parent != nil {
		info.ParentID = issue.Fields.Parent.ID
		info.ParentKey = issue.Fields.Parent.Key
	}

	// Extract custom fields
	info.TargetVersion = extractCustomField(issue.Fields.Unknowns, customFieldTargetVersion)
	info.ReleaseBlocker = extractCustomField(issue.Fields.Unknowns, customFieldReleaseBlocker)

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

	sort.Strings(tickets)

	jiraClient, err := o.jira.Client()
	if err != nil {
		logrus.WithError(err).Fatal("cannot create Jira client")
	}

	ctx := context.Background()
	jql := "issue in (linkedIssues(OTA-1643), linkedIssues(OTA-1626), linkedIssues(OTA-362), linkedIssues(TRT-1578), linkedIssues(OTA-1637)) AND project = \"OpenShift Bugs\""
	// Currently about 70 in total
	issues, _, err := jiraClient.SearchV2JqlWithContext(ctx, jql, &andyjira.SearchOptionsV2{MaxResults: 200, Fields: []string{"id", "key", "parent"}})
	if err != nil {
		logrus.WithError(err).Fatal("cannot search Jira issues")
	}
	logrus.Infof("Found %d Jira tickets on the dashboard", len(issues))

	delta := sets.New[string](tickets...)
	for _, issue := range issues {
		var parentKey string
		if issue.Fields != nil && issue.Fields.Parent != nil {
			parentKey = issue.Fields.Parent.Key
		}
		logrus.WithField("id", issue.ID).WithField("key", issue.Key).WithField("parentKey", parentKey).Debug("Found issue")

		delta.Delete(issue.Key)
	}
	if delta.Len() > 0 {
		logrus.WithField("delta", sets.List[string](delta)).Fatal("jira issues in o/origin but not on the dashboard")
	}

	var ticketInfos []*ticketInfo

	for i, ticketID := range tickets {
		logrus.Infof("Fetching ticket: %s", ticketID)
		info, err := fetchTicketInfo(jiraClient, ticketID)
		if err != nil {
			logrus.WithError(err).Fatalf("cannot fetch ticket %s", ticketID)
		}
		info.Number = i
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

		// Define workflow order for status display
		workflowOrder := []string{"New", "ASSIGNED", "POST", "ON_QA", "Verified", "Closed"}

		// Filter out descoped/won't-do tickets and track them separately
		var activeTickets []*ticketInfo
		var descopedTickets []string
		var descopedTicketKeys []string // Track keys for deduplication
		var wontDoTickets []string

		for _, ticket := range ticketInfos {
			linkedKey := fmt.Sprintf("[%s](%s)", ticket.Key, ticket.URL)
			if strings.Contains(ticket.Notes, "Won't Do confirmed") {
				wontDoTickets = append(wontDoTickets, linkedKey)
			} else if strings.Contains(ticket.Notes, "Descoped") {
				descopedTickets = append(descopedTickets, linkedKey)
				descopedTicketKeys = append(descopedTicketKeys, ticket.Key)
			} else {
				activeTickets = append(activeTickets, ticket)
			}
		}

		// Find other descoped bugs from notes map that aren't in descopedTickets
		var descopedOthers []string
		for key, note := range notes {
			if strings.Contains(note, "Descoped") {
				// Check if this key is not already in descopedTicketKeys
				found := false
				for _, dk := range descopedTicketKeys {
					if key == dk {
						found = true
						break
					}
				}
				if !found {
					// Try to find the URL for this ticket in ticketInfos
					ticketFound := false
					for _, ticket := range ticketInfos {
						if ticket.Key == key {
							linkedKey := fmt.Sprintf("[%s](%s)", ticket.Key, ticket.URL)
							descopedOthers = append(descopedOthers, linkedKey)
							ticketFound = true
							break
						}
					}
					// If ticket not in ticketInfos (not found in source files), construct URL manually
					if !ticketFound {
						linkedKey := fmt.Sprintf("[%s](https://redhat.atlassian.net/browse/%s)", key, key)
						descopedOthers = append(descopedOthers, linkedKey)
					}
				}
			}
		}

		// Renumber active tickets
		for i := range activeTickets {
			activeTickets[i].Number = i
		}

		// Count statuses and Pixaa components
		statusCounts := make(map[string]int)
		pixaaCount := 0
		pixaaNotClosedCount := 0

		buf.WriteString("## OCPBugs on [jira/dashboards/22315](https://redhat.atlassian.net/jira/dashboards/22315) as exceptions in CI\n")

		// Build table rows and count statuses
		var tableRows strings.Builder
		tableRows.WriteString("| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Pixaa | Assignee | Parent | Notes |\n")
		tableRows.WriteString("|---|-----|---------|--------|------------|----------------|-----------------|-----------|-------|----------|--------|-------|\n")

		for _, ticket := range activeTickets {
			statusCounts[ticket.Status]++
			if ticket.IsPixaaComponent {
				pixaaCount++
				if ticket.Status != "Closed" {
					pixaaNotClosedCount++
				}
			}

			escapedSummary := strings.ReplaceAll(ticket.Summary, "|", "\\|")
			escapedComponent := strings.ReplaceAll(ticket.Component, "|", "\\|")
			escapedResolution := strings.ReplaceAll(ticket.Resolution, "|", "\\|")
			escapedAssignee := strings.ReplaceAll(ticket.Assignee, "|", "\\|")
			escapedNotes := strings.ReplaceAll(ticket.Notes, "|", "\\|")
			escapedTargetVersion := strings.ReplaceAll(ticket.TargetVersion, "|", "\\|")
			escapedReleaseBlocker := strings.ReplaceAll(ticket.ReleaseBlocker, "|", "\\|")

			keyField := fmt.Sprintf("[%s](%s)", ticket.Key, ticket.URL)
			if strings.Contains(ticket.Notes, "Won't Do confirmed") {
				keyField = fmt.Sprintf("~~%s~~", keyField)
			}

			parentKeyField := ""
			if ticket.ParentKey != "" {
				parentKeyField = ticket.ParentKey
			}

			pixaaField := ""
			if ticket.IsPixaaComponent {
				pixaaField = "✓"
			}

			tableRows.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
				ticket.Number,
				keyField,
				escapedSummary,
				ticket.Status,
				escapedResolution,
				escapedTargetVersion,
				escapedReleaseBlocker,
				escapedComponent,
				pixaaField,
				escapedAssignee,
				parentKeyField,
				escapedNotes))
		}

		// Generate status summary
		var summaryParts []string
		total := len(activeTickets)
		summaryParts = append(summaryParts, fmt.Sprintf("Total: %d issues", total))
		summaryParts = append(summaryParts, fmt.Sprintf("Pixaa: %d (not closed: %d)", pixaaCount, pixaaNotClosedCount))

		// Add counts for statuses in workflow order
		for _, status := range workflowOrder {
			if count, exists := statusCounts[status]; exists && count > 0 {
				summaryParts = append(summaryParts, fmt.Sprintf("%s: %d", status, count))
			}
		}

		// Add counts for any remaining statuses not in workflow order (alphabetically)
		var otherStatuses []string
		for status := range statusCounts {
			found := false
			for _, wfStatus := range workflowOrder {
				if status == wfStatus {
					found = true
					break
				}
			}
			if !found {
				otherStatuses = append(otherStatuses, status)
			}
		}
		sort.Strings(otherStatuses)
		for _, status := range otherStatuses {
			if count := statusCounts[status]; count > 0 {
				summaryParts = append(summaryParts, fmt.Sprintf("%s: %d", status, count))
			}
		}

		statusSummary := strings.Join(summaryParts, " | ")

		// Write summary and table to buffer
		buf.WriteString(statusSummary)
		buf.WriteString("\n\n")

		// add Pixaa aggregated info
		buf.WriteString(fmt.Sprintf("PIXAA: open/total: %d/%d\n\n", pixaaNotClosedCount, pixaaCount))

		// Add Pixaa components list
		pixaaComponentsList := sets.List(pixaaComponents)
		sort.Strings(pixaaComponentsList)
		buf.WriteString(fmt.Sprintf("PIXAA components: %s\n\n", strings.Join(pixaaComponentsList, ", ")))

		if len(descopedTickets) > 0 {
			buf.WriteString(fmt.Sprintf("Descoped (not in table): %d issues - %s\n\n", len(descopedTickets), strings.Join(descopedTickets, ", ")))
		}

		if len(descopedOthers) > 0 {
			buf.WriteString(fmt.Sprintf("Descoped (others): %d issues - %s\n\n", len(descopedOthers), strings.Join(descopedOthers, ", ")))
		} else {
			buf.WriteString("Descoped (others): 0 issues\n\n")
		}

		if len(wontDoTickets) > 0 {
			buf.WriteString(fmt.Sprintf("Won't Do: %d issues - %s\n\n", len(wontDoTickets), strings.Join(wontDoTickets, ", ")))
		}

		buf.WriteString(tableRows.String())

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
