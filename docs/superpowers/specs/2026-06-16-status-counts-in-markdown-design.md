# Status Counts in Markdown Output - Design Specification

**Date:** 2026-06-16  
**Tool:** `cmd/co-conditions-bugs`  
**Scope:** Add status counting to markdown output

## Overview

Add a status summary line to the markdown output that shows the total count of issues and a breakdown by status. The summary will appear at the top of the document, immediately after the header and before the issue table.

## Requirements

### Output Format
- Display format: `Total: X issues | Status1: N | Status2: M | ...`
- Placement: After the header, before the table, with blank line separation
- Only show statuses with count > 0
- Example: `Total: 30 issues | New: 6 | ASSIGNED: 4 | POST: 5 | ON_QA: 1 | Verified: 2 | Closed: 3`

### Status Ordering
Statuses will be displayed in workflow order:
1. New
2. ASSIGNED
3. POST
4. ON_QA
5. Verified
6. Closed

Any statuses not in this predefined list will appear alphabetically at the end.

## Implementation Approach

### Data Collection
- Use `map[string]int` to count status occurrences
- Count during the existing iteration over `ticketInfos` in the markdown generation case
- No additional iteration needed

### Workflow Order Definition
- Define a constant slice with the canonical workflow order
- Place near the markdown generation logic in `main()`
- Use this to order the output statuses

### Output Generation
1. Build status counts map during ticket iteration
2. Generate summary line using workflow order
3. Insert summary after header line
4. Add blank line before table

## Code Changes

### File: `cmd/co-conditions-bugs/main.go`

**Location:** In the `case "md":` section of the `main()` function (around line 300)

**Changes:**
1. Define workflow order constant at the start of the md case
2. Initialize status count map before the ticket iteration loop
3. Increment status counts during ticket iteration
4. After the loop, generate the status summary string
5. Insert summary into output buffer after the header

**Algorithm for status summary generation:**
1. Start with total count
2. Iterate through workflow order slice
3. For each status in workflow order, if count > 0, append to summary
4. Collect any remaining statuses not in workflow order
5. Sort remaining statuses alphabetically and append if count > 0
6. Join all parts with " | " separator

## Expected Output Example

```markdown
## OCPBugs on [jira/dashboards/22315](https://redhat.atlassian.net/jira/dashboards/22315) as exceptions in CI
Total: 30 issues | New: 6 | ASSIGNED: 4 | POST: 5 | ON_QA: 1 | Verified: 2 | Closed: 3

| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Assignee | Notes |
|---|-----|---------|--------|------------|----------------|-----------------|-----------|----------|-------|
...
```

## Non-Goals

- This design does not modify the JSON output format
- No changes to the table structure or content
- No filtering or grouping by status
- No percentage calculations or statistical analysis

## Testing

Manual testing approach:
1. Run the tool against current data
2. Verify status counts match table entries
3. Verify workflow ordering is correct
4. Verify statuses with 0 count are not shown
5. Verify unknown statuses (if any) appear alphabetically at the end
