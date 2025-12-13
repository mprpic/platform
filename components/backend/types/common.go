// Package types defines common type definitions for AgenticSession, ProjectSettings, and RFE workflows.
package types

// Common types used across the application

// GitRepository represents a Git repository configuration
// @Description Git repository URL and branch information
type GitRepository struct {
	URL      string       `json:"url" example:"https://github.com/org/repo.git"`
	Branch   *string      `json:"branch,omitempty" example:"main"`
	Provider ProviderType `json:"provider,omitempty"` // Optional: auto-detected if not specified
}

// UserContext contains user identity information
// @Description User authentication and authorization context
type UserContext struct {
	UserID      string   `json:"userId" binding:"required" example:"user123"`
	DisplayName string   `json:"displayName" binding:"required" example:"John Doe"`
	Groups      []string `json:"groups" binding:"required" example:"developers,admins"`
}

// BotAccountRef references a bot service account
// @Description Reference to a Kubernetes service account for bot operations
type BotAccountRef struct {
	Name string `json:"name" binding:"required" example:"bot-account"`
}

// ResourceOverrides specifies custom resource limits
// @Description Kubernetes resource overrides for session pods
type ResourceOverrides struct {
	CPU           string `json:"cpu,omitempty" example:"2000m"`
	Memory        string `json:"memory,omitempty" example:"4Gi"`
	StorageClass  string `json:"storageClass,omitempty" example:"fast-ssd"`
	PriorityClass string `json:"priorityClass,omitempty" example:"high-priority"`
}

// LLMSettings configures the language model
// @Description Language model configuration including model, temperature, and token limits
type LLMSettings struct {
	Model       string  `json:"model" example:"claude-sonnet-4"`
	Temperature float64 `json:"temperature" example:"0.7"`
	MaxTokens   int     `json:"maxTokens" example:"4096"`
}

type GitConfig struct {
	Repositories []GitRepository `json:"repositories,omitempty"`
}

type Paths struct {
	Workspace string `json:"workspace,omitempty"`
	Messages  string `json:"messages,omitempty"`
	Inbox     string `json:"inbox,omitempty"`
}

// Common repository browsing types (used by both GitHub and GitLab)

// Branch represents a Git branch (common format for UI)
type Branch struct {
	Name      string     `json:"name"`
	Protected bool       `json:"protected"`
	Default   bool       `json:"default,omitempty"`
	Commit    CommitInfo `json:"commit,omitempty"`
}

// CommitInfo represents basic commit information
type CommitInfo struct {
	SHA       string `json:"sha"`
	Message   string `json:"message,omitempty"`
	Author    string `json:"author,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// TreeEntry represents a file or directory in a repository
type TreeEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "blob" (file) or "tree" (directory)
	Mode string `json:"mode,omitempty"`
	SHA  string `json:"sha,omitempty"`
	Size int    `json:"size,omitempty"`
}

// FileContent represents file contents from a repository
type FileContent struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"` // "base64" or "utf-8"
	Size     int    `json:"size"`
	SHA      string `json:"sha,omitempty"`
}

// BoolPtr returns a pointer to the given bool value.
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns a pointer to the given string value.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int value.
func IntPtr(i int) *int {
	return &i
}

// PaginationParams represents common pagination request parameters
// @Description Pagination parameters for list endpoints
type PaginationParams struct {
	Limit    int    `form:"limit" example:"20"`     // Number of items per page (default: 20, max: 100)
	Offset   int    `form:"offset" example:"0"`     // Offset for offset-based pagination
	Continue string `form:"continue" example:""`    // Continuation token for k8s-style pagination
	Search   string `form:"search" example:"query"` // Search/filter term
}

// PaginatedResponse is a generic paginated response structure
// @Description Generic paginated response with items and pagination metadata
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalCount int         `json:"totalCount" example:"100"`
	Limit      int         `json:"limit" example:"20"`
	Offset     int         `json:"offset" example:"0"`
	HasMore    bool        `json:"hasMore" example:"true"`
	Continue   string      `json:"continue,omitempty" example:"continue-token"`   // For k8s-style pagination
	NextOffset *int        `json:"nextOffset,omitempty" example:"20"` // For offset-based pagination
}

// DefaultPaginationLimit is the default number of items per page
const DefaultPaginationLimit = 20

// MaxPaginationLimit is the maximum allowed items per page
const MaxPaginationLimit = 100

// NormalizePaginationParams ensures pagination params are within valid bounds
func NormalizePaginationParams(params *PaginationParams) {
	if params.Limit <= 0 {
		params.Limit = DefaultPaginationLimit
	}
	if params.Limit > MaxPaginationLimit {
		params.Limit = MaxPaginationLimit
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
}
