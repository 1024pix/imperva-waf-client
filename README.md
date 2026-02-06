# Imperva Cloud WAF Go Client

A Golang client for the Imperva Cloud WAF API.

## Features

This client implements the following Imperva Cloud WAF APIs:

* **Custom Rules API** (v2)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/rules-api-definition.htm?operationId=operations-Rules-postsitessiteIdrules)
    * Implemented: Create, Read, Update, Delete, List Rules.
* **Session Management API** (v3)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/session-release-api.htm?operationId=operations-Session_Release_API-releaseSession)
    * Implemented: Release Blocked Session.
* **Traffic Statistics & Logs** (v1)
    * [Documentation](https://docs-cybersec.thalesgroup.com/bundle/api-docs/page/traffic-stats-api-definition.htm)
    * Implemented: Get Visits (Logs), Get Statistics.

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

See `cmd/example/main.go` for a complete example.

```bash
go run cmd/example/main.go
```

### Safety Warning

The example CLI acts on the `site_id` specified in your configuration. It will prompt for confirmation before proceeding to ensure you are not running tests against a production site unintentionally.
