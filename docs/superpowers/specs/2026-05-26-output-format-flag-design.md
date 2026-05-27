# Output Format Flag Design

**Date:** 2026-05-26

## Overview

Add `--output-format` and `--output-file` flags to the `co-conditions-bugs` command to support both JSON (default) and Markdown table output formats, with the ability to write to a file or stdout.

## Purpose

Enable users to output Jira ticket information in different formats and destinations for different use cases:
- JSON for programmatic consumption and integration with other tools
- Markdown for human-readable reports and documentation
- File output for saving results permanently
- Stdout output for piping to other commands

## Architecture

### Flag Addition & Validation

**New fields in `options` struct:**
```go
type options struct {
    jira            flagutil.JiraOptions
    originDirectory string
    outputFormat    string  // NEW: "json" or "md"
    outputFile      string  // NEW: file path or "-" for stdout
}
```

**Flag registration in `gatherOptions()`:**
```go
fs.StringVar(&o.outputFormat, "output-format", "json", "Output format (json or md)")
fs.StringVar(&o.outputFile, "output-file", "-", "Output file path (use - for stdout)")
```

**Validation in `validate()` method:**
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

### Output Formatting and File Writing

**Replace the current JSON-only output block in `main()` with format switch and file writing logic:**

The implementation will:
1. Format the output based on `outputFormat` (json or md)
2. Write to a file if `outputFile` is a path, or to stdout if `outputFile` is "-"

**Format generation:**
```go
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
    buf.WriteString("|-----|---------|-----------|----------|\n")
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
```

**File writing logic:**
```go
if o.outputFile == "-" {
    fmt.Print(output)
} else {
    if err := os.WriteFile(o.outputFile, []byte(output), 0644); err != nil {
        logrus.WithError(err).Fatalf("cannot write to file %s", o.outputFile)
    }
    logrus.Infof("Output written to %s", o.outputFile)
}
```

### Markdown Table Format

The markdown output will be a standard GitHub-flavored markdown table with clickable links in the Key column:

```
| Key | Summary | Status | Component |
|-----|---------|--------|-----------|
| [OCPBUGS-12345](https://redhat.atlassian.net/browse/OCPBUGS-12345) | Bug description | Open | Component Name |
| [OCPBUGS-67890](https://redhat.atlassian.net/browse/OCPBUGS-67890) | Another bug | Closed | Other Component |
```

**Special character handling:**
- Pipe characters (`|`) in ticket data will be escaped or replaced to avoid breaking the table format
- Use `strings.ReplaceAll(field, "|", "\\|")` for each field value
- Key column uses markdown link format: `[KEY](URL)` for clickable links to Jira tickets

### Error Handling

- Invalid output format values are caught during validation with a clear error message
- JSON marshaling errors continue to use the existing error handling pattern
- Markdown output has no marshaling step, so no additional error handling needed beyond escaping

## Testing Strategy

Manual testing with both formats:

1. **JSON output (default):**
```bash
./co-conditions-bugs --jira-username user@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```
Expected: JSON array output (existing behavior)

2. **JSON output (explicit):**
```bash
./co-conditions-bugs -o json --jira-username user@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```
Expected: Same JSON array output

3. **Markdown output:**
```bash
./co-conditions-bugs -o md --jira-username user@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```
Expected: Markdown table with ticket information

4. **Invalid format:**
```bash
./co-conditions-bugs -o xml --jira-username user@redhat.com --jira-password-file ~/.config/ota/jira-api-token
```
Expected: Error message "invalid output format "xml", must be 'json' or 'md'"

5. **Tickets with special characters:**
- Test with tickets that have pipe characters in summary/component
- Verify markdown table is not broken

## Implementation Details

### Import Requirements

Need to add `strings` import for pipe character escaping and string building:
```go
import (
    // existing imports...
    "strings"
)
```

### File Writing

When `outputFile` is "-" (the default), write to stdout using `fmt.Print()`.

When `outputFile` is a path:
- Use `os.WriteFile()` to write the formatted output
- Use file permissions `0644` (rw-r--r--)
- Log success message with the file path
- Exit with fatal error if write fails

### Helper Function (Optional)

For cleaner code, could extract markdown formatting to a helper function:
```go
func formatMarkdown(tickets []*ticketInfo) {
    fmt.Println("| Key | Summary | Status | Component |")
    fmt.Println("|-----|---------|--------|-----------|")
    for _, ticket := range tickets {
        escapedSummary := strings.ReplaceAll(ticket.Summary, "|", "\\|")
        escapedComponent := strings.ReplaceAll(ticket.Component, "|", "\\|")
        fmt.Printf("| %s | %s | %s | %s |\n", 
            ticket.Key, 
            escapedSummary, 
            ticket.Status, 
            escapedComponent)
    }
}
```

However, keeping it inline in the switch statement is acceptable given the simplicity.

## File Structure

Only one file modified:
- `cmd/co-conditions-bugs/main.go`

## Future Enhancements (Not in Scope)

- Additional output formats (CSV, YAML, HTML)
- Custom markdown templates
- Colorized terminal output
- Output to file instead of stdout

## Documentation Updates

Update the package comment to reflect the new flag:

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
```
