package types

// AgenticSession represents the structure of our custom resource
// @Description Kubernetes custom resource representing an AI-powered agentic session
type AgenticSession struct {
	APIVersion string                 `json:"apiVersion" example:"vteam.ambient-code/v1alpha1"`
	Kind       string                 `json:"kind" example:"AgenticSession"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       AgenticSessionSpec     `json:"spec"`
	Status     *AgenticSessionStatus  `json:"status,omitempty"`
}

// AgenticSessionSpec defines the desired state of an agentic session
// @Description Configuration for an agentic session including prompt, repos, and LLM settings
type AgenticSessionSpec struct {
	InitialPrompt        string             `json:"initialPrompt,omitempty"`
	Interactive          bool               `json:"interactive,omitempty"`
	DisplayName          string             `json:"displayName"`
	LLMSettings          LLMSettings        `json:"llmSettings"`
	Timeout              int                `json:"timeout"`
	UserContext          *UserContext       `json:"userContext,omitempty"`
	BotAccount           *BotAccountRef     `json:"botAccount,omitempty"`
	ResourceOverrides    *ResourceOverrides `json:"resourceOverrides,omitempty"`
	EnvironmentVariables map[string]string  `json:"environmentVariables,omitempty"`
	Project              string             `json:"project,omitempty"`
	// Multi-repo support
	Repos []SimpleRepo `json:"repos,omitempty"`
	// Active workflow for dynamic workflow switching
	ActiveWorkflow *WorkflowSelection `json:"activeWorkflow,omitempty"`
}

// SimpleRepo represents a simplified repository configuration
type SimpleRepo struct {
	URL    string  `json:"url"`
	Branch *string `json:"branch,omitempty"`
}

type AgenticSessionStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Phase              string              `json:"phase,omitempty"`
	StartTime          *string             `json:"startTime,omitempty"`
	CompletionTime     *string             `json:"completionTime,omitempty"`
	ReconciledRepos    []ReconciledRepo    `json:"reconciledRepos,omitempty"`
	ReconciledWorkflow *ReconciledWorkflow `json:"reconciledWorkflow,omitempty"`
	SDKSessionID       string              `json:"sdkSessionId,omitempty"`
	SDKRestartCount    int                 `json:"sdkRestartCount,omitempty"`
	Conditions         []Condition         `json:"conditions,omitempty"`
}

// CreateAgenticSessionRequest represents the request body for creating an agentic session
// @Description Request payload for creating a new agentic session with AI agent configuration
type CreateAgenticSessionRequest struct {
	InitialPrompt   string       `json:"initialPrompt,omitempty" example:"Implement feature XYZ with tests"`
	DisplayName     string       `json:"displayName,omitempty" example:"Feature XYZ Implementation"`
	LLMSettings     *LLMSettings `json:"llmSettings,omitempty"`
	Timeout         *int         `json:"timeout,omitempty" example:"3600"`
	Interactive     *bool        `json:"interactive,omitempty" example:"false"`
	ParentSessionID string       `json:"parent_session_id,omitempty" example:"parent-session-123"`
	// Multi-repo support
	Repos                []SimpleRepo      `json:"repos,omitempty"`
	AutoPushOnComplete   *bool             `json:"autoPushOnComplete,omitempty" example:"true"`
	UserContext          *UserContext      `json:"userContext,omitempty"`
	EnvironmentVariables map[string]string `json:"environmentVariables,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	Annotations          map[string]string `json:"annotations,omitempty"`
}

// CloneSessionRequest represents the request body for cloning a session
// @Description Request to clone an existing session to a new project
type CloneSessionRequest struct {
	TargetProject  string `json:"targetProject" binding:"required" example:"my-new-project"`
	NewSessionName string `json:"newSessionName" binding:"required" example:"cloned-session"`
}

// UpdateAgenticSessionRequest represents the request body for updating a session
// @Description Request payload for updating an existing agentic session
type UpdateAgenticSessionRequest struct {
	InitialPrompt *string      `json:"initialPrompt,omitempty" example:"Updated prompt text"`
	DisplayName   *string      `json:"displayName,omitempty" example:"Updated Display Name"`
	Timeout       *int         `json:"timeout,omitempty" example:"7200"`
	LLMSettings   *LLMSettings `json:"llmSettings,omitempty"`
}

type CloneAgenticSessionRequest struct {
	TargetProject     string `json:"targetProject,omitempty"`
	TargetSessionName string `json:"targetSessionName,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	InitialPrompt     string `json:"initialPrompt,omitempty"`
}

// WorkflowSelection represents a workflow to load into the session
type WorkflowSelection struct {
	GitURL string `json:"gitUrl" binding:"required"`
	Branch string `json:"branch,omitempty"`
	Path   string `json:"path,omitempty"`
}

// ReconciledRepo captures reconciliation state for a repository
type ReconciledRepo struct {
	URL      string  `json:"url"`
	Branch   string  `json:"branch"`
	Name     string  `json:"name,omitempty"`
	Status   string  `json:"status,omitempty"`
	ClonedAt *string `json:"clonedAt,omitempty"`
}

// ReconciledWorkflow captures reconciliation state for the active workflow
type ReconciledWorkflow struct {
	GitURL    string  `json:"gitUrl"`
	Branch    string  `json:"branch"`
	Path      string  `json:"path,omitempty"`
	Status    string  `json:"status,omitempty"`
	AppliedAt *string `json:"appliedAt,omitempty"`
}

// Condition mirrors metav1.Condition for API transport
type Condition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	ObservedGeneration int64  `json:"observedGeneration,omitempty"`
}
