# Target Version and Release Blocker Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Target Version and Release Blocker custom fields to the co-conditions-bugs tool's ticketInfo output in both JSON and Markdown formats.

**Architecture:** Extract Red Hat Jira custom fields from the go-jira library's `Unknowns` map, add them to the ticketInfo struct, and display them in both output formats. The implementation includes a field discovery phase to find the correct custom field IDs, followed by extraction logic and output formatting changes.

**Tech Stack:** Go, go-jira (andygrunwald/go-jira v1.17.0), Prow Jira client

---

### Task 1: Discover Custom Field IDs

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

This task adds temporary debug code to discover the actual custom field IDs for "Target Version" and "Release Blocker" in Red Hat's Jira instance.

- [ ] **Step 1: Add debug output to fetchTicketInfo**

Add this code after line 151 (after `issue, err := jiraClient.GetIssue(ticketID)`):

```go
	// Temporary debug: dump all custom fields
	logrus.Info("=== Custom Fields Debug ===")
	for key, value := range issue.Fields.Unknowns {
		logrus.Infof("Field ID: %s, Value: %+v", key, value)
	}
	logrus.Info("=== End Custom Fields Debug ===")
```

- [ ] **Step 2: Build and run the tool**

```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format json | head -50
```

Expected: Console output showing all custom field IDs and their values for the first ticket.

- [ ] **Step 3: Identify the field IDs**

Look through the debug output for fields matching:
- "Target Version" or similar (likely has a value like "4.18.0")
- "Release Blocker" or similar (likely has a value like "Approved" or "Rejected")

Note the field IDs (format: `customfield_XXXXX`) for use in the next task.

- [ ] **Step 4: Do not commit yet**

We'll remove the debug code and add the proper implementation in the next tasks.

---

### Task 2: Add Helper Function and Constants

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Remove the debug code from Task 1**

Delete the debug lines added after line 151 in `fetchTicketInfo`.

- [ ] **Step 2: Add custom field ID constants**

Add these constants after the `notes` map declaration (after line 145), replacing `XXXXX` and `YYYYY` with the actual field IDs discovered in Task 1:

```go
const (
	customFieldTargetVersion  = "customfield_XXXXX"
	customFieldReleaseBlocker = "customfield_YYYYY"
)
```

- [ ] **Step 3: Add extractCustomField helper function**

Add this function after the `parseJiraTickets` function (after line 136):

```go
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

	return ""
}
```

- [ ] **Step 4: Verify code compiles**

```bash
go build ./cmd/co-conditions-bugs
```

Expected: Successful compilation with no errors.

- [ ] **Step 5: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add helper function for extracting Jira custom fields

Add extractCustomField helper to safely extract string values from
Jira's Unknowns map with type assertion handling.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3: Extend ticketInfo Struct

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go:39-49`

- [ ] **Step 1: Add new fields to ticketInfo struct**

Modify the `ticketInfo` struct to add the two new fields at the end:

```go
type ticketInfo struct {
	Number         int    `json:"count"`
	Key            string `json:"key"`
	Summary        string `json:"summary"`
	Status         string `json:"status"`
	Component      string `json:"component"`
	Resolution     string `json:"resolution"`
	Assignee       string `json:"assignee"`
	URL            string `json:"url"`
	Notes          string `json:"notes"`
	TargetVersion  string `json:"target_version"`
	ReleaseBlocker string `json:"release_blocker"`
}
```

- [ ] **Step 2: Verify code compiles**

```bash
go build ./cmd/co-conditions-bugs
```

Expected: Successful compilation with no errors.

- [ ] **Step 3: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add TargetVersion and ReleaseBlocker fields to ticketInfo

Extend ticketInfo struct to include target_version and release_blocker
fields for JSON output.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4: Extract Custom Fields in fetchTicketInfo

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go:147-174`

- [ ] **Step 1: Add field extraction to fetchTicketInfo**

Add these lines before the `return info, nil` statement (after line 171, after the assignee extraction):

```go
	// Extract custom fields
	info.TargetVersion = extractCustomField(issue.Fields.Unknowns, customFieldTargetVersion)
	info.ReleaseBlocker = extractCustomField(issue.Fields.Unknowns, customFieldReleaseBlocker)
```

- [ ] **Step 2: Build and test JSON output**

```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format json > /tmp/test-output.json
```

Expected: Successful execution with no errors.

- [ ] **Step 3: Verify JSON contains new fields**

```bash
cat /tmp/test-output.json | jq '.[0] | {key: .key, target_version: .target_version, release_blocker: .release_blocker}'
```

Expected: Output showing the key, target_version, and release_blocker fields. Values should be populated if the Jira issue has these fields set, or empty strings if not.

- [ ] **Step 4: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Extract Target Version and Release Blocker from Jira

Use extractCustomField helper to populate TargetVersion and
ReleaseBlocker fields from Jira custom fields.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 5: Update Markdown Output Format

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go:253-280`

- [ ] **Step 1: Update markdown table header**

Modify line 256 to include the new columns after "Resolution":

```go
			buf.WriteString("| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Assignee | Notes |\n")
```

- [ ] **Step 2: Update markdown table separator**

Modify line 257 to match the new column count:

```go
			buf.WriteString("|---|-----|---------|--------|------------|----------------|-----------------|-----------|----------|-------|\n")
```

- [ ] **Step 3: Add field escaping for new columns**

Add these lines after line 263 (after `escapedAssignee` assignment):

```go
				escapedTargetVersion := strings.ReplaceAll(ticket.TargetVersion, "|", "\\|")
				escapedReleaseBlocker := strings.ReplaceAll(ticket.ReleaseBlocker, "|", "\\|")
```

- [ ] **Step 4: Update markdown table row format**

Modify the `buf.WriteString(fmt.Sprintf(...))` call (lines 270-278) to include the new fields:

```go
				buf.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
					ticket.Number,
					keyField,
					escapedSummary,
					ticket.Status,
					escapedResolution,
					escapedTargetVersion,
					escapedReleaseBlocker,
					escapedComponent,
					escapedAssignee,
					escapedNotes))
```

- [ ] **Step 5: Build and test markdown output**

```bash
go build ./cmd/co-conditions-bugs
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format md > /tmp/test-output.md
```

Expected: Successful execution with no errors.

- [ ] **Step 6: Verify markdown table structure**

```bash
head -5 /tmp/test-output.md
```

Expected: Output showing a table with the new "Target Version" and "Release Blocker" columns after "Resolution" and before "Component". The header row and separator should align properly.

- [ ] **Step 7: Commit**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "Add Target Version and Release Blocker to markdown output

Include new columns in markdown table format with proper escaping
and alignment.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 6: End-to-End Verification

**Files:**
- Test: `cmd/co-conditions-bugs/main.go`

- [ ] **Step 1: Test with actual data - JSON format**

```bash
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format json --output-file /tmp/final-test.json
cat /tmp/final-test.json | jq '.[0:3] | .[] | {key: .key, status: .status, target_version: .target_version, release_blocker: .release_blocker}'
```

Expected: First 3 tickets showing key, status, target_version, and release_blocker fields. Values should be populated or empty strings.

- [ ] **Step 2: Test with actual data - Markdown format**

```bash
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format md --output-file /tmp/final-test.md
head -10 /tmp/final-test.md
```

Expected: Markdown table with proper formatting, including the new columns. All pipe characters in field values should be escaped.

- [ ] **Step 3: Verify existing output file workflow**

```bash
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format md --output-file ./output/co-conditions-bugs.md
cat ./output/co-conditions-bugs.md | head -20
```

Expected: File written to ./output/co-conditions-bugs.md with the new columns included. Table should render correctly with all fields present.

- [ ] **Step 4: Test with missing field values**

Review the output and verify that:
- Tickets without Target Version show empty cells (not errors)
- Tickets without Release Blocker show empty cells (not errors)
- No crashes or warnings about missing fields

Expected: Clean output with empty strings for missing values, no errors.

- [ ] **Step 5: Final commit and update output file if needed**

```bash
# If output/co-conditions-bugs.md exists and should be regenerated:
./co-conditions-bugs --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --origin-directory ~/repo/openshift/origin --output-format md --output-file ./output/co-conditions-bugs.md
git add ./output/co-conditions-bugs.md
git commit -m "Regenerate output with Target Version and Release Blocker

Update output file to include new custom fields.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Notes

- **Field ID Discovery:** Task 1 requires manual inspection of debug output to find the correct custom field IDs. These IDs must be used in Task 2's constants.
- **Empty Values:** The implementation gracefully handles missing or null custom field values by returning empty strings.
- **Type Safety:** The `extractCustomField` helper handles multiple potential return types from Jira's API (direct strings, objects with `name` or `value` fields).
- **Breaking Changes:** None - the changes are purely additive. Existing JSON and Markdown output remains compatible.
