# co-conditions-bugs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a CLI tool that extracts Jira ticket URLs from two Go source files and outputs ticket details as JSON.

**Architecture:** Simple imperative Go program following the monitor-jira-dashboard pattern. Single main.go file with flag parsing, file reading, regex extraction, Jira API calls, and JSON output.

**Tech Stack:** Go standard library (flag, os, regexp, encoding/json, path/filepath), go-jira, logrus, internal/flagutil

---

## File Structure

**New files:**
- `cmd/co-conditions-bugs/main.go` - Complete implementation (main function, option parsing, file reading, Jira interaction, JSON output)

**No modifications to existing files.**

---

### Task 1: Basic CLI structure with flags

**Files:**
- Create: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Create main.go with package declaration and imports**

```go
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
```

- [ ] **Step 2: Build and test flag parsing**

Run:
```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --help
```

Expected: Help text showing `--origin-directory`, `--jira-username`, `--jira-password-file` flags

- [ ] **Step 3: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add basic CLI structure for co-conditions-bugs

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2: Tilde expansion for origin directory

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Add expandPath helper function**

Add after the `validate()` method:

```go
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
```

Add import: `"path/filepath"`

- [ ] **Step 2: Call expandPath in main function**

Replace the `logrus.Info("co-conditions-bugs started")` line in `main()` with:

```go
	expandedDir, err := expandPath(o.originDirectory)
	if err != nil {
		logrus.WithError(err).Fatal("cannot expand origin directory path")
	}

	logrus.Infof("Using origin directory: %s", expandedDir)
```

- [ ] **Step 3: Test tilde expansion**

Run:
```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Log output showing expanded path (e.g., "/Users/username/repo/openshift/origin")

- [ ] **Step 4: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add tilde expansion for origin directory path

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3: Read source files

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Add file reading function**

Add before `main()`:

```go
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
```

- [ ] **Step 2: Call readSourceFiles in main**

Replace `logrus.Infof("Using origin directory: %s", expandedDir)` with:

```go
	logrus.Infof("Using origin directory: %s", expandedDir)

	content, err := readSourceFiles(expandedDir)
	if err != nil {
		logrus.WithError(err).Fatal("cannot read source files")
	}

	logrus.Infof("Read %d bytes from source files", len(content))
```

- [ ] **Step 3: Test file reading (will fail without actual files)**

Run:
```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --origin-directory /tmp/test --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Error "cannot read source files" (unless /tmp/test has those files)

- [ ] **Step 4: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add source file reading functionality

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4: Parse Jira URLs with regex

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Add parseJiraTickets function**

Add import: `"regexp"`

Add before `main()`:

```go
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
```

- [ ] **Step 2: Call parseJiraTickets in main**

Replace `logrus.Infof("Read %d bytes from source files", len(content))` with:

```go
	logrus.Infof("Read %d bytes from source files", len(content))

	tickets := parseJiraTickets(content)
	logrus.Infof("Found %d unique Jira tickets", len(tickets))

	if len(tickets) == 0 {
		logrus.Info("No Jira tickets found")
		return
	}
```

- [ ] **Step 3: Test parsing (create test file)**

Create a test file:
```bash
mkdir -p /tmp/test/pkg/monitortests/clusterversionoperator/legacycvomonitortests
echo '// https://issues.redhat.com/browse/OCPBUGS-12345' > /tmp/test/pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
mkdir -p /tmp/test/test/extended/machines
echo '// https://bugzilla.redhat.com/browse/OCPBUGS-67890' > /tmp/test/test/extended/machines/scale.go
```

Run:
```bash
./co-conditions-bugs --origin-directory /tmp/test --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Log output "Found 2 unique Jira tickets"

- [ ] **Step 4: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add Jira ticket URL parsing with deduplication

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 5: Fetch Jira ticket information

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Add ticketInfo struct and import context**

Add imports:
```go
"context"

"github.com/andygrunwald/go-jira"
```

Add after `options` struct:

```go
type ticketInfo struct {
	Key       string `json:"key"`
	Summary   string `json:"summary"`
	Status    string `json:"status"`
	Component string `json:"component"`
}
```

- [ ] **Step 2: Add fetchTicketInfo function**

Add before `main()`:

```go
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
```

- [ ] **Step 3: Add Jira client creation and ticket fetching in main**

Replace the block after `if len(tickets) == 0` with:

```go
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
```

- [ ] **Step 4: Test with invalid credentials (expected to fail)**

Run:
```bash
echo "invalid-token" > /tmp/token
./co-conditions-bugs --origin-directory /tmp/test --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Error "cannot create Jira client" or "cannot fetch ticket" (authentication error)

- [ ] **Step 5: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add Jira ticket fetching via API

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 6: Output JSON

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Add JSON output**

Add import: `"encoding/json"`

Replace `logrus.Infof("Successfully fetched %d tickets", len(ticketInfos))` with:

```go
	logrus.Infof("Successfully fetched %d tickets", len(ticketInfos))

	output, err := json.MarshalIndent(ticketInfos, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatal("cannot marshal JSON output")
	}

	fmt.Println(string(output))
```

Add import: `"fmt"`

- [ ] **Step 2: Test JSON output with mock data**

Create test files with real-looking ticket references:
```bash
echo '// See https://issues.redhat.com/browse/OCPBUGS-11111
// Also https://bugzilla.redhat.com/browse/OCPBUGS-22222
// Duplicate: https://issues.redhat.com/browse/OCPBUGS-11111' > /tmp/test/pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
```

Run:
```bash
./co-conditions-bugs --origin-directory /tmp/test --jira-username valid@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```

Expected: JSON array output with ticket details (if credentials are valid and tickets exist)

- [ ] **Step 3: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add JSON output for ticket information

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 7: Final testing and documentation

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go` (add package comment)

- [ ] **Step 1: Add package documentation comment**

At the top of `main.go`, before `package main`, add:

```go
// co-conditions-bugs extracts Jira ticket references from specific Go source files
// in the OpenShift origin repository and outputs their details as JSON.
//
// It reads two hardcoded files:
//   - pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
//   - test/extended/machines/scale.go
//
// Usage:
//
//	co-conditions-bugs --origin-directory ~/repo/openshift/origin \
//	  --jira-username user@redhat.com \
//	  --jira-password-file ~/.config/ota/jira-api-token
package main
```

- [ ] **Step 2: Build final binary**

Run:
```bash
go build ./cmd/co-conditions-bugs
```

Expected: Clean build with no errors

- [ ] **Step 3: Test with actual origin repository (if available)**

If you have the origin repository locally:
```bash
./co-conditions-bugs --jira-username your-email@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```

Expected: JSON output with actual ticket information from the origin repository files

- [ ] **Step 4: Test with custom directory**

```bash
./co-conditions-bugs --origin-directory /path/to/custom/dir --jira-username user@redhat.com --jira-password-file /path/to/token
```

Expected: Reads files from custom directory

- [ ] **Step 5: Final commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add package documentation for co-conditions-bugs

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Self-Review Checklist

**Spec Coverage:**
- ✅ CLI structure with flags (--origin-directory, Jira auth)
- ✅ Default value for --origin-directory (~/repo/openshift/origin)
- ✅ Tilde expansion for directory paths
- ✅ Read two specific Go files from directory
- ✅ Regex parsing of Jira URLs with pattern `https://.*/browse/(OCPBUGS-\d+)`
- ✅ Deduplication of ticket IDs
- ✅ Jira API client creation using flagutil.JiraOptions
- ✅ Fetch ticket details (Key, Summary, Status, Component)
- ✅ JSON output with indentation
- ✅ Error handling with logrus.Fatal pattern

**Placeholder Check:**
- ✅ No TBD, TODO, or "implement later"
- ✅ All code blocks complete and executable
- ✅ All file paths exact and complete
- ✅ All commands include expected output

**Type Consistency:**
- ✅ `options` struct consistent across tasks
- ✅ `ticketInfo` struct used consistently
- ✅ Function names consistent (expandPath, readSourceFiles, parseJiraTickets, fetchTicketInfo)
- ✅ Variable names consistent (expandedDir, content, tickets, ticketInfos)

**Build Verification:**
- Each task produces buildable code
- Commits are atomic and focused
- TDD principles followed where practical (build → test → commit)
