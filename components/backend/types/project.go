package types

// AmbientProject represents project management types.
// @Description Kubernetes/OpenShift project (namespace) with metadata
type AmbientProject struct {
	Name              string            `json:"name" example:"my-project"`                  // Kubernetes namespace name
	DisplayName       string            `json:"displayName" example:"My Project"`           // OpenShift display name (empty for k8s)
	Description       string            `json:"description,omitempty" example:"Project for AI automation"` // OpenShift description (empty for k8s)
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp string            `json:"creationTimestamp" example:"2025-01-15T10:30:00Z"`
	Status            string            `json:"status" example:"Active"`
	IsOpenShift       bool              `json:"isOpenShift" example:"true"` // true if running on OpenShift cluster
}

// CreateProjectRequest represents the request body for creating a project
// @Description Request payload for creating a new project (Kubernetes namespace)
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required" example:"my-new-project"`
	DisplayName string `json:"displayName,omitempty" example:"My New Project"` // Optional: only used on OpenShift
	Description string `json:"description,omitempty" example:"Description of the project"` // Optional: only used on OpenShift
}
