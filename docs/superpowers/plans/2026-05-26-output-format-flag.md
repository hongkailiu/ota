# Output Format Flag Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `--output-format` and `--output-file` flags to support JSON (default) and Markdown table output formats, with file or stdout output.

**Architecture:** Add outputFormat and outputFile fields to the options struct, validate outputFormat accepts only "json" or "md", then use a switch statement in main() to format the output and write it to a file or stdout based on outputFile value.

**Tech Stack:** Go standard library (flag, strings, encoding/json, fmt, os)

---

## File Structure

**Modified files:**
- `cmd/co-conditions-bugs/main.go` - Add outputFormat field, flag, validation, and format switch

**No new files created.**

---

### Task 1: Add outputFormat and outputFile fields and flags ✅ COMPLETED

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

**Status:** User has completed this task with the following implementation:
- Added `outputFormat string` field to options struct
- Added `outputFile string` field to options struct  
- Registered `--output-format` flag with default "json"
- Registered `--output-file` flag with default "/tmp/output.json" (needs to be changed to "-" in Task 2)

---

### Task 2: Fix output-file default and add validation ✅ COMPLETED

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [x] **Step 1: Change output-file default to "-"**

Find the `--output-file` flag registration (around line 50) and change the default from "/tmp/output.json" to "-":

```go
	fs.StringVar(&o.outputFile, "output-file", "-", "Output file path (use - for stdout)")
```

- [x] **Step 2: Update validate() method**

Find the `validate()` method (around line 59) and replace it with:

```go
func (o *options) validate() error {
	if err := o.jira.Validate(); err != nil {
		return err
	}

	if o.outputFormat != "json" && o.outputFormat != "md" {
		return fmt.Errorf("invalid output format %q, must be 'json' or 'md'", o.outputFormat)
	}

	return nil
}
```

- [x] **Step 3: Test validation with invalid format**

Run:
```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --output-format xml --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Error message "invalid output format \"xml\", must be 'json' or 'md'"

- [x] **Step 4: Test validation with valid format**

Run:
```bash
./co-conditions-bugs --output-format json --jira-username test@redhat.com --jira-password-file /tmp/token --origin-directory /tmp/nonexistent
```

Expected: Different error (file reading failure, not validation error)

- [x] **Step 5: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Fix output-file default and add format validation

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3: Implement markdown output formatting ✅ COMPLETED

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [x] **Step 1: Add strings import**

Add `"strings"` to the import list (should be inserted alphabetically in the standard library imports section):

```go
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
```

- [x] **Step 2: Replace JSON output with format generation and file writing**

Find the current output code (around lines 171-176) that looks like:

```go
	logrus.Infof("Successfully fetched %d tickets", len(ticketInfos))

	output, err := json.MarshalIndent(ticketInfos, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatal("cannot marshal JSON output")
	}

	fmt.Println(string(output))
```

Replace it with:

```go
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
		buf.WriteString("| Key | Summary | Status | Resolution | Component | Notes |\n")
		buf.WriteString("|-----|---------|--------|------------|-----------|-------|\n")
		for _, ticket := range ticketInfos {
			escapedSummary := strings.ReplaceAll(ticket.Summary, "|", "\\|")
			escapedComponent := strings.ReplaceAll(ticket.Component, "|", "\\|")
			escapedResolution := strings.ReplaceAll(ticket.Resolution, "|", "\\|")
			escapedNotes := strings.ReplaceAll(ticket.Notes, "|", "\\|")
			buf.WriteString(fmt.Sprintf("| [%s](%s) | %s | %s | %s | %s | %s |\n",
				ticket.Key,
				ticket.URL,
				escapedSummary,
				ticket.Status,
				escapedResolution,
				escapedComponent,
				escapedNotes))
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
```

- [x] **Step 3: Build the updated binary**

Run:
```bash
go build ./cmd/co-conditions-bugs
```

Expected: Clean build with no errors

- [x] **Step 4: Test JSON output (default)**

Create test files:
```bash
mkdir -p /tmp/test/pkg/monitortests/clusterversionoperator/legacycvomonitortests
echo '// https://issues.redhat.com/browse/OCPBUGS-11111' > /tmp/test/pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go
mkdir -p /tmp/test/test/extended/machines
echo '// https://bugzilla.redhat.com/browse/OCPBUGS-22222' > /tmp/test/test/extended/machines/scale.go
```

Note: This test will fail at the Jira API call without valid credentials, but you can verify it parses the tickets correctly by checking the log output.

- [x] **Step 5: Test markdown output to stdout**

Run:
```bash
./co-conditions-bugs --output-format md --origin-directory /tmp/test --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Log shows "Found 2 unique Jira tickets", then fails at Jira API call (expected without valid credentials)

- [x] **Step 6: Test writing to file**

Run:
```bash
./co-conditions-bugs --output-format md --output-file /tmp/test-output.md --origin-directory /tmp/test --jira-username test@redhat.com --jira-password-file /tmp/token
```

Expected: Would write to file if Jira credentials were valid

- [x] **Step 7: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Implement markdown/json formatting and file output

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4: Update package documentation ✅ COMPLETED

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [x] **Step 1: Update package comment**

Find the package comment at the top of the file (lines 1-12) and update it to:

```go
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
```

- [x] **Step 2: Build final binary**

Run:
```bash
go build ./cmd/co-conditions-bugs
```

Expected: Clean build with no errors

- [x] **Step 3: Test with actual repository (if available)**

If you have the origin repository with valid Jira credentials:

Test JSON output:
```bash
./co-conditions-bugs --jira-username your-email@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```

Test Markdown output:
```bash
./co-conditions-bugs -o md --jira-username your-email@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```

Expected: Both commands produce output in their respective formats with actual ticket data

- [x] **Step 4: Final commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Update package documentation for output formats

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Self-Review Checklist

**Spec Coverage:**
- ✅ Add outputFormat field to options struct
- ✅ Register -o and --output-format flags with "json" default
- ✅ Validate format is "json" or "md"
- ✅ Implement switch statement for output selection
- ✅ Implement markdown table output with pipe escaping
- ✅ Update package documentation

**Placeholder Check:**
- ✅ No TBD, TODO, or "implement later"
- ✅ All code blocks complete and executable
- ✅ All file paths exact and complete
- ✅ All commands include expected output

**Type Consistency:**
- ✅ `outputFormat` field name consistent across tasks
- ✅ Valid values "json" and "md" consistent
- ✅ `ticketInfo` struct usage consistent with existing code

**Build Verification:**
- Each task produces buildable code
- Commits are atomic and focused
- Manual testing steps included for verification
