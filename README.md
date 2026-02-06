# Imperva Cloud WAF Go Client

A Golang client for the Imperva Cloud WAF API.

## Features

This client implements the following Imperva Cloud WAF APIs:

* **Custom Rules API** (v2)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/rules-api-definition.htm?operationId=operations-Rules-postsitessiteIdrules)
* **Session Management API** (v3)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/session-release-api.htm?operationId=operations-Session_Release_API-releaseSession)
* **Traffic Statistics & Logs** (v1)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/traffic-stats-api-definition.htm)

## Implemented Endpoints

### Site Management (v1)
*   `ListSites`: Lists all sites for the account (`POST /api/prov/v1/sites/list`)
*   `GetSiteStatus`: Retrieves the status of a specific site (`POST /api/prov/v1/sites/status`)

### Custom Rules (v2 & v3)
*   `ListRules`: Lists all rules for a site (`GET /api/prov/v3/rules`)
*   `CreateRule`: Creates a new custom rule (`POST /api/prov/v2/sites/{siteId}/rules`)
*   `GetRule`: Retrieves a specific rule (`GET /api/prov/v2/sites/{siteId}/rules/{ruleId}`)
*   `UpdateRule`: Updates an existing rule (`POST /api/prov/v2/sites/{siteId}/rules/{ruleId}`)
*   `DeleteRule`: Deletes a rule (`DELETE /api/prov/v2/sites/{siteId}/rules/{ruleId}`)

### Session Management (v3)
*   `ReleaseSession`: Releases a blocked session (`POST /v3/sites/{siteId}/sessions/{sessionId}/release`)

### Traffic Statistics & Logs (v1)
*   `GetVisits`: Retrieves traffic logs/visits (`POST /api/visits/v1`)
*   `GetStats`: Retrieves aggregated traffic statistics (`POST /api/stats/v1`)

## Installation

```bash
go get github.com/1024pix/imperva-waf-client
```

## Configuration

Copy `config.json.example` to `config.json` and fill in your credentials:

```json
{
    "api_id": "YOUR_API_ID",
    "api_key": "YOUR_API_KEY",
    "account_id": "YOUR_ACCOUNT_ID",
    "site_id": 12345,
    "host": "https://my.imperva.com"
}
```

## Usage

The example CLI (`cmd/example`) demonstrates how to use the client.

### Running the Example

Since the example is split across multiple files in the `main` package, you must run it by targeting the directory or including all files:

```bash
# Run with interactive mode (default)
go run ./cmd/example

# Run a specific action on a specific site (non-interactive)
go run ./cmd/example -site=12345 -action=rules
```

**Available Flags:**
*   `-site`: Specify the Site ID to test (skips interactive selection).
*   `-action`: Action to perform: `status`, `rules`, `stats`, or `all` (default: `all`).
*   `-config`: Path to configuration file (default: `config.json`).
```

### Safety Warning

The example CLI acts on the `site_id` specified in your configuration. It will prompt for confirmation before proceeding to ensure you are not running tests against a production site unintentionally.
