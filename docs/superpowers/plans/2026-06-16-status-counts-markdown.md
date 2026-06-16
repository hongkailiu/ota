# Status Counts in Markdown Output Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add status count summary line to markdown output showing total issues and breakdown by status in workflow order.

**Architecture:** Modify the markdown generation case in `cmd/co-conditions-bugs/main.go` to count statuses during iteration and insert a summary line after the header. Use a predefined workflow order with fallback to alphabetical for unknown statuses.

**Tech Stack:** Go standard library (strings, sort)

---

## File Structure

**Files to Modify:**
- `cmd/co-conditions-bugs/main.go:300-331` - Add status counting and summary generation to markdown case

**Files to Test:**
- Manual testing: run tool and verify output format

---

### Task 1: Add Status Counting Logic

**Files:**
- Modify: `cmd/co-conditions-bugs/main.go:300-331`

- [ ] **Step 1: Add workflow order constant and status count map initialization**

Add after line 301 (after `var buf strings.Builder`):

```go
// Define workflow order for status display
workflowOrder := []string{"New", "ASSIGNED", "POST", "ON_QA", "Verified", "Closed"}

// Count statuses
statusCounts := make(map[string]int)
```

- [ ] **Step 2: Add status counting in the ticket iteration loop**

In the existing `for _, ticket := range ticketInfos` loop (line 305), add this as the first line inside the loop:

```go
statusCounts[ticket.Status]++
```

- [ ] **Step 3: Generate status summary string after the loop**

Add after the closing `}` of the for loop (after line 330), before `output = buf.String()`:

```go
// Generate status summary
var summaryParts []string
total := len(ticketInfos)
summaryParts = append(summaryParts, fmt.Sprintf("Total: %d issues", total))

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
```

- [ ] **Step 4: Insert summary into output buffer**

The current code builds the output in `buf` with:
1. Header line (line 302)
2. Table header (line 303)
3. Table separator (line 304)
4. Rows (lines 305-330)

We need to insert the summary after the header line. Replace lines 301-331 with:

```go
var buf strings.Builder

// Define workflow order for status display
workflowOrder := []string{"New", "ASSIGNED", "POST", "ON_QA", "Verified", "Closed"}

// Count statuses
statusCounts := make(map[string]int)

buf.WriteString("## OCPBugs on [jira/dashboards/22315](https://redhat.atlassian.net/jira/dashboards/22315) as exceptions in CI\n")

// Build table rows and count statuses
var tableRows strings.Builder
tableRows.WriteString("| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Assignee | Notes |\n")
tableRows.WriteString("|---|-----|---------|--------|------------|----------------|-----------------|-----------|----------|-------|\n")

for _, ticket := range ticketInfos {
	statusCounts[ticket.Status]++
	
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

	tableRows.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
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
}

// Generate status summary
var summaryParts []string
total := len(ticketInfos)
summaryParts = append(summaryParts, fmt.Sprintf("Total: %d issues", total))

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
buf.WriteString(tableRows.String())

output = buf.String()
```

- [ ] **Step 5: Verify the code compiles**

Run: `go build ./cmd/co-conditions-bugs`

Expected: No compilation errors

- [ ] **Step 6: Commit the changes**

```bash
git add cmd/co-conditions-bugs/main.go
git commit -m "feat: add status counts to markdown output

Add status summary line showing total issue count and breakdown
by status in workflow order (New, ASSIGNED, POST, ON_QA, Verified, Closed).
Summary appears after header and before table.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2: Manual Testing

**Files:**
- Test: Run the tool and verify output

- [ ] **Step 1: Build the tool**

Run: `go build ./cmd/co-conditions-bugs`

Expected: Binary created successfully

- [ ] **Step 2: Run the tool with current data**

Run:
```bash
./co-conditions-bugs --origin-directory ~/repo/openshift/origin --output-format md --jira-username "kerberosid@redhat.com" --jira-password-file ~/.config/ota/.jira-api-token --output-file ./output/co-conditions-bugs.md
```

Expected: File generated at `./output/co-conditions-bugs.md`

- [ ] **Step 3: Verify output format**

Run: `head -10 ./output/co-conditions-bugs.md`

Expected output should look like:
```
## OCPBugs on [jira/dashboards/22315](https://redhat.atlassian.net/jira/dashboards/22315) as exceptions in CI
Total: 30 issues | New: 6 | ASSIGNED: 4 | POST: 5 | ON_QA: 1 | Verified: 2 | Closed: 3

| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Assignee | Notes |
|---|-----|---------|--------|------------|----------------|-----------------|-----------|----------|-------|
```

(Exact counts will vary based on current data)

- [ ] **Step 4: Manually verify status counts match table**

Count each status in the table manually and compare to the summary line.

Run: `grep -oP '(?<=\| )[^|]+(?= \|)' ./output/co-conditions-bugs.md | awk 'NR % 10 == 4' | sort | uniq -c`

Expected: Counts should match the summary line

- [ ] **Step 5: Verify workflow ordering**

Check that statuses appear in order: New, ASSIGNED, POST, ON_QA, Verified, Closed (for those with count > 0).

Any unknown statuses should appear alphabetically after these.

- [ ] **Step 6: Commit the regenerated output**

```bash
git add ./output/co-conditions-bugs.md
git commit -m "Regenerate markdown output with status counts

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Self-Review Checklist

**Spec Coverage:**
- ✓ Output format: `Total: X issues | Status1: N | ...` 
- ✓ Placement: After header, before table, with blank line
- ✓ Only show statuses with count > 0
- ✓ Workflow order: New, ASSIGNED, POST, ON_QA, Verified, Closed
- ✓ Unknown statuses alphabetically at end
- ✓ Count during iteration (no additional loop)
- ✓ No changes to JSON output
- ✓ No changes to table structure

**Placeholder Check:**
- ✓ No TBD/TODO markers
- ✓ All code blocks complete
- ✓ All commands have expected output
- ✓ Exact file paths provided

**Type Consistency:**
- ✓ `statusCounts` is `map[string]int` throughout
- ✓ `workflowOrder` is `[]string` throughout
- ✓ `ticket.Status` is string (existing field)
