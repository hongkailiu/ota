# Design: Add Target Version and Release Blocker to ticketInfo

**Date:** 2026-06-11  
**Status:** Approved  
**Component:** cmd/co-conditions-bugs

## Overview

Extend the `co-conditions-bugs` tool to retrieve and display two additional Jira custom fields for each issue: "Target Version" and "Release Blocker". These fields provide critical context about which release version a bug targets and its release blocker approval status.

## Requirements

- Add "Target Version" field (single version string, e.g., "4.18.0")
- Add "Release Blocker" field (status string, e.g., "Approved", "Rejected")
- Display both fields in JSON output format
- Display both fields in Markdown table output format
- Handle missing/null field values gracefully (empty strings)
- No breaking changes to existing output structure

## Architecture

### Approach

Use the go-jira library's `issue.Fields.Unknowns` map to access Red Hat Jira custom fields. Custom fields in Jira have IDs like `customfield_XXXXX` and are not part of the standard go-jira `Issue` struct.

### Components Modified

**1. ticketInfo struct** (`cmd/co-conditions-bugs/main.go:39-49`)
- Add `TargetVersion string` field with JSON tag `target_version`
- Add `ReleaseBlocker string` field with JSON tag `release_blocker`

**2. fetchTicketInfo function** (`cmd/co-conditions-bugs/main.go:147-174`)
- Extract custom field values from `issue.Fields.Unknowns` map
- Use helper function for safe type assertion and extraction
- Populate new ticketInfo fields

**3. Markdown output generation** (`cmd/co-conditions-bugs/main.go:253-280`)
- Add "Target Version" and "Release Blocker" columns to table header
- Insert columns after "Resolution" and before "Component"
- Apply pipe character escaping to field values

## Data Flow

### Field Discovery Phase

Before implementation, discover the exact custom field IDs:
1. Fetch a sample Jira issue using the existing Jira client
2. Inspect `issue.Fields.Unknowns` map contents
3. Identify field IDs matching "Target Version" and "Release Blocker"
4. Common Red Hat Jira custom field IDs to check:
   - Target Version: `customfield_12319940` (typical pattern)
   - Release Blocker: `customfield_12316840` (typical pattern)

### Runtime Extraction

1. `fetchTicketInfo` receives `*jira.Issue` from API call
2. For each custom field:
   - Check if field ID exists in `issue.Fields.Unknowns`
   - Perform safe type assertion (handle string, object, or array types)
   - Extract string value (for objects, use `.Name` or `.Value` sub-field)
   - Default to empty string if missing, null, or type mismatch
3. Populate `ticketInfo` struct fields
4. Return populated struct

### Output Generation

**JSON format:**
- New fields automatically included via struct JSON tags
- No code changes needed beyond struct definition

**Markdown format:**
- New table structure:
  ```
  | # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Assignee | Notes |
  ```
- Escape pipe characters in field values using `strings.ReplaceAll`
- Maintain consistent column alignment

## Implementation Details

### Helper Function

Create `extractCustomField` helper function:

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

### Field ID Constants

Define constants for custom field IDs after discovery:

```go
const (
    customFieldTargetVersion  = "customfield_XXXXX"  // To be discovered
    customFieldReleaseBlocker = "customfield_YYYYY"  // To be discovered
)
```

### Modified ticketInfo Struct

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

## Error Handling

### Missing or Invalid Field IDs
- If custom field IDs are incorrect or have changed, fields will be empty
- No fatal errors - gracefully degrade to empty values
- Consider adding debug logging when fields are not found

### Type Assertion Failures
- Custom fields may return unexpected types (string, object, array)
- Use two-value type assertions (`value, ok := ...`)
- For complex types, try extracting common sub-fields (`Name`, `Value`)
- Default to empty string on type mismatch

### Null or Missing Values
- Some issues may not have these fields populated
- Treat as empty strings in output
- No special placeholder text ("N/A", "Unknown") needed

## Testing Approach

1. **Field Discovery Test:** Run tool against known issue to discover field IDs
2. **Populated Fields Test:** Verify extraction works for issues with values set
3. **Empty Fields Test:** Verify graceful handling when fields are null/missing
4. **JSON Output Test:** Verify new fields appear in JSON output
5. **Markdown Output Test:** Verify table renders correctly with new columns

## Non-Goals

- Not adding filtering or sorting by these new fields
- Not validating field values (e.g., version format validation)
- Not handling multiple target versions (single string only)
- Not backfilling historical data or caching

## Future Considerations

- If field IDs change across Jira environments, may need configuration
- Could add validation to warn about unexpected field value formats
- May want to make field IDs configurable via flags or config file
