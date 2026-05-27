# co-conditions-bugs Command Design

**Date:** 2026-05-26

## Overview

The `co-conditions-bugs` command is a simple CLI tool that reads two Go source files from a local directory, parses Jira ticket URLs from those files, fetches ticket information via the Jira API, and outputs the results as JSON.

## Purpose

Enable automated extraction and reporting of Jira bug references embedded in Go source code files, specifically targeting cluster operator condition-related bugs referenced in the codebase.

## Architecture

### Command Flow

1. **Parse CLI Arguments**
   - Accept a directory path via `--origin-directory` flag (default: `~/repo/openshift/origin`)
   - Expand `~` to user's home directory if present
   - Accept Jira authentication flags via `flagutil.JiraOptions` (username, password file)
   - Validate all required flags are provided

2. **Read Source Files**
   - Construct paths to two hardcoded files relative to the directory:
     - `pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go`
     - `test/extended/machines/scale.go`
   - Read entire file contents into memory

3. **Parse Jira References**
   - Use regex pattern `https://.*/browse/(OCPBUGS-\d+)` to extract ticket IDs
   - Deduplicate ticket IDs using a map (to handle tickets referenced in both files or multiple times)

4. **Fetch Jira Data**
   - Create Jira client using existing `flagutil.JiraOptions.Client()`
   - For each unique ticket ID, call Jira API to fetch issue details
   - Extract relevant fields: Key, Summary, Status, Component

5. **Output Results**
   - Marshal the ticket information to JSON
   - Print to stdout with indentation for readability

### Data Structures

```go
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
```

### Error Handling

Follow the `monitor-jira-dashboard` pattern:
- Use `logrus.Fatal()` for unrecoverable errors (missing flags, file not found, Jira API errors)
- Exit with non-zero status on any error
- No graceful degradation - fail fast and clearly

### Dependencies

- Standard library: `flag`, `os`, `regexp`, `encoding/json`, `path/filepath`
- External: `github.com/sirupsen/logrus`, `github.com/andygrunwald/go-jira`
- Internal: `github.com/petr-muller/ota/internal/flagutil`

## Implementation Details

### File Reading

Use `os.ReadFile()` to read both files completely into memory. Files are expected to be source code files, typically small (<100KB), so full memory read is acceptable.

Handle tilde expansion for the origin directory path using `os.UserHomeDir()` to replace `~` with the user's home directory path.

### Regex Parsing

Compile regex once: `regexp.MustCompile(`https://.*/browse/(OCPBUGS-\d+)`)`
Use `FindAllStringSubmatch()` to extract all matches and capture groups.

### Deduplication

Use a `map[string]bool` to track seen ticket IDs during parsing. Only fetch each unique ticket once.

### Jira Ticket Fetching

Iterate through unique ticket IDs and fetch each one sequentially. If any Jira API call fails (network error, ticket not found, authentication failure), exit immediately with logrus.Fatal(). Do not attempt to continue fetching remaining tickets.

### Jira Field Extraction

For Component field, handle the case where multiple components exist by:
- Taking the first component if available
- Using empty string if no components

For Status, use `issue.Fields.Status.Name`

### JSON Output

Output an array of `ticketInfo` structs using `json.MarshalIndent()` with 2-space indentation for human readability.

## Testing Strategy

Manual testing approach:
1. Create test directory with sample .go files containing Jira URLs
2. Run command and verify JSON output
3. Test with duplicate tickets across files
4. Test with no tickets found
5. Test with invalid directory path

Future enhancement: unit tests could be added for the parsing logic if extracted to a separate function.

## Command Line Interface

```bash
./co-conditions-bugs --origin-directory /path/to/source/directory
```

### Flags

- `--origin-directory`: Path to the directory containing the source files (default: `~/repo/openshift/origin`)

## Example Output

```json
[
  {
    "key": "OCPBUGS-12345",
    "summary": "Cluster operator degraded condition stuck",
    "status": "In Progress",
    "component": "Cluster Version Operator"
  },
  {
    "key": "OCPBUGS-67890",
    "summary": "Available condition flapping on upgrade",
    "status": "Closed",
    "component": "etcd"
  }
]
```

## File Structure

```
cmd/co-conditions-bugs/
└── main.go
```

Single file implementation following the simplicity of `monitor-jira-dashboard`.

## Future Enhancements (Not in Scope)

- Make file paths configurable via flags instead of hardcoded
- Support for multiple file patterns/globs
- CSV or table output format options
- Filtering by status or component
- Caching of Jira results to avoid repeated API calls
