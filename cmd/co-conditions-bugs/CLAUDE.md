# co-conditions-bugs

This tool extracts and reports on Jira ticket references found in specific OpenShift origin repository source files.

## Purpose

Monitors OpenShift bugs (OCPBUGS-*) that are referenced as exceptions in CI test code, specifically tracking issues that appear in:
- `pkg/monitortests/clusterversionoperator/legacycvomonitortests/operators.go`
- `test/extended/machines/scale.go`

The tool verifies these bugs exist on the [CO Conditions dashboard](https://redhat.atlassian.net/jira/dashboards/22315) and outputs their current status.

## Authentication

**IMPORTANT:** Unlike other tools in this repository, the Jira issues accessed by this tool are **publicly accessible**. You do **not** need to set up Jira authentication credentials or API tokens to run this tool.

While the tool accepts `--jira-username` and `--jira-password-file` flags (inherited from `internal/flagutil`), they are **optional** for this use case.

## Building

```bash
go build ./cmd/co-conditions-bugs
```

## Running

```bash
# Basic usage - outputs JSON to stdout
./co-conditions-bugs --origin-directory ~/repo/openshift/origin

# Output as Markdown
./co-conditions-bugs --origin-directory ~/repo/openshift/origin --output-format md

# Write to file
./co-conditions-bugs --origin-directory ~/repo/openshift/origin --output-format md --output-file report.md

# Convenience script - generates ./output/co-conditions-bugs.md with build info
./hack/co-conditions-bugs.sh
```

## How It Works

1. **Extract**: Parses hardcoded Go source files from the origin repository for Jira ticket URLs (OCPBUGS-\d+)
2. **Validate**: Queries Jira to ensure all found tickets are linked to one of the CO Conditions tracking issues (OTA-1643, OTA-1626, OTA-362, TRT-1578, OTA-1637)
3. **Report**: Fetches full details for each ticket and outputs as JSON or Markdown

## Output Formats

### JSON
Contains array of ticket objects with fields: key, summary, status, component, resolution, assignee, target_version, release_blocker, notes

### Markdown
Generates a formatted table with:
- Status count summary (Total, New, ASSIGNED, POST, ON_QA, Verified, Closed)
- Full details table with all tickets
- Strikethrough formatting for "Won't Do confirmed" tickets

## Notes Field

The tool includes hardcoded notes for specific tickets (see `notes` map in main.go) to track additional context like:
- Confirmed Won't Do status
- Duplicate relationships
- Removal plans
- Evaluation status
