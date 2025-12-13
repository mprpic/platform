# API Reference

The Ambient Code Platform provides a RESTful API for managing projects, agentic sessions, and workflows.

## Interactive Documentation

### Runtime Swagger UI

Access interactive API documentation on a running backend instance:

- **Local Development**: [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)
- **Production**: `https://<backend-url>/swagger/index.html`

### Embedded Documentation

<swagger-ui src="../openapi.yaml"/>

## Quick Start

### Authentication

All API endpoints require Bearer token authentication:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://api.ambient-code.com/api/projects
```

### Base URLs

- **Local**: `http://localhost:8081/api`
- **Production**: `https://<backend-url>/api`

## Common Endpoints

### Projects

```bash
# List projects
GET /api/projects

# Create project
POST /api/projects
{
  "name": "my-project",
  "displayName": "My Project",
  "description": "My AI project"
}

# Get project details
GET /api/projects/{projectName}

# Delete project
DELETE /api/projects/{projectName}
```

### Agentic Sessions

```bash
# Create session
POST /api/projects/{projectName}/agentic-sessions
{
  "displayName": "Feature Implementation",
  "initialPrompt": "Implement feature X",
  "repos": [{"url": "https://github.com/org/repo", "branch": "main"}],
  "interactive": false,
  "timeout": 3600
}

# List sessions
GET /api/projects/{projectName}/agentic-sessions

# Get session details
GET /api/projects/{projectName}/agentic-sessions/{sessionName}

# Start session
POST /api/projects/{projectName}/agentic-sessions/{sessionName}/start

# Stop session
POST /api/projects/{projectName}/agentic-sessions/{sessionName}/stop

# Delete session
DELETE /api/projects/{projectName}/agentic-sessions/{sessionName}
```

### Session Workspace Operations

```bash
# List workspace files
GET /api/projects/{projectName}/agentic-sessions/{sessionName}/workspace?path=/

# Read file content
GET /api/projects/{projectName}/agentic-sessions/{sessionName}/workspace/file/{path}

# Write file content
PUT /api/projects/{projectName}/agentic-sessions/{sessionName}/workspace/file/{path}
```

### Git Operations

```bash
# Get git status
GET /api/projects/{projectName}/agentic-sessions/{sessionName}/git/status?path=repo-name

# Configure git remote
POST /api/projects/{projectName}/agentic-sessions/{sessionName}/git/configure-remote
{
  "path": "repo-name",
  "remoteUrl": "https://github.com/org/repo.git",
  "branch": "main"
}

# Push changes
POST /api/projects/{projectName}/agentic-sessions/{sessionName}/git/push
{
  "path": "repo-name",
  "branch": "feature-branch",
  "message": "Commit message"
}

# Pull changes
POST /api/projects/{projectName}/agentic-sessions/{sessionName}/git/pull
{
  "path": "repo-name",
  "branch": "main"
}
```

### Repository Management

```bash
# Get repository tree
GET /api/projects/{projectName}/repo/tree?repo=org/repo&ref=main&path=/

# Get file blob
GET /api/projects/{projectName}/repo/blob?repo=org/repo&ref=main&path=README.md

# List branches
GET /api/projects/{projectName}/repo/branches?repo=org/repo

# Check seed status
GET /api/projects/{projectName}/repo/seed-status?repo=https://github.com/org/repo

# Seed repository
POST /api/projects/{projectName}/repo/seed
{
  "repositoryUrl": "https://github.com/org/repo",
  "branch": "main",
  "force": false
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Human-readable error message"
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 202 | Accepted (async operation started) |
| 204 | No Content (successful deletion) |
| 400 | Bad Request (invalid parameters) |
| 401 | Unauthorized (invalid/missing token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (resource state prevents operation) |
| 500 | Internal Server Error |
| 502 | Bad Gateway (external service failure) |
| 503 | Service Unavailable (content service unavailable) |

## Authentication & Authorization

### Bearer Token Authentication

All API requests must include a valid Bearer token in the Authorization header:

```
Authorization: Bearer <your-token>
```

### Token Acquisition

**For OpenShift/Kubernetes clusters:**

```bash
# Get service account token
kubectl create token <service-account-name> -n <namespace>

# Or use your user token
oc whoami -t
```

**For GitHub App integration:**

Tokens are minted automatically when you connect your GitHub account via the web UI.

### RBAC Permissions

The API enforces Kubernetes RBAC policies. Users can only:

- List/view projects they have access to
- Create/modify sessions in projects where they have edit/admin permissions
- View session details based on their project permissions

Project-level permissions are managed via:

```bash
# Add user to project
POST /api/projects/{projectName}/permissions
{
  "subjectType": "user",
  "subjectName": "user@example.com",
  "role": "edit"  # or "view", "admin"
}

# List project permissions
GET /api/projects/{projectName}/permissions

# Remove user from project
DELETE /api/projects/{projectName}/permissions/user/{subjectName}
```

## Pagination

List endpoints support pagination with these query parameters:

- `limit`: Number of items per page (default: 20, max: 100)
- `offset`: Starting offset (default: 0)
- `search`: Filter by name/displayName/prompt (optional)

**Response format:**

```json
{
  "items": [...],
  "totalCount": 100,
  "limit": 20,
  "offset": 0,
  "nextOffset": 20,
  "hasMore": true
}
```

## Rate Limiting

Currently no rate limiting is enforced.

## API Versioning

Current API version: `v1alpha1`

⚠️ **Note**: The API is in alpha. Breaking changes may occur without notice.

## WebSocket Support

Real-time session updates are available via WebSocket:

```javascript
const ws = new WebSocket('ws://localhost:8081/ws/projects/{projectName}/sessions/{sessionName}');

ws.onmessage = (event) => {
  const update = JSON.parse(event.data);
  console.log('Session update:', update);
};
```

## Secrets Management

### Runner Secrets (Anthropic API Key)

```bash
# Get runner secrets
GET /api/projects/{projectName}/runner-secrets

# Update runner secrets
PUT /api/projects/{projectName}/runner-secrets
{
  "ANTHROPIC_API_KEY": "your-api-key"
}
```

### Integration Secrets (GitHub, Jira, etc.)

```bash
# Get integration secrets
GET /api/projects/{projectName}/integration-secrets

# Update integration secrets
PUT /api/projects/{projectName}/integration-secrets
{
  "GITHUB_TOKEN": "ghp_...",
  "JIRA_API_TOKEN": "...",
  "JIRA_EMAIL": "user@example.com",
  "JIRA_SERVER": "https://your-instance.atlassian.net"
}
```

## Additional Resources

- **GitHub Repository**: [ambient-code/platform](https://github.com/ambient-code/platform)
- **Issue Tracker**: [GitHub Issues](https://github.com/ambient-code/platform/issues)
- **User Guide**: [Getting Started](../user-guide/getting-started.md)
- **Developer Guide**: [Overview](../index.md)
