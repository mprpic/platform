# Feature Specification: Automated CI Debugging Workflows

**Feature Branch**: `001-ci-debug-workflows`
**Created**: 2026-02-17
**Status**: Draft
**Input**: User description: "Automated CI Debugging Workflows for ACP — poll CI platforms for failures, spin up debug sessions, propose fixes, and notify users"

## Clarifications

### Session 2026-02-17

- Q: How long should ProcessedFailure records be retained for deduplication? → A: Retain forever (until the parent workflow is deleted). All failure history is preserved for the lifetime of the workflow.
- Q: Should workflows auto-pause after persistent CI platform errors, or keep retrying? → A: Auto-pause after 10 consecutive failures. User must manually resume after fixing the issue.
- Q: Should there be a maximum number of automated workflows per project? → A: No limit. Users can create as many workflows as needed.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create an Automated CI Watch (Priority: P1)

A platform user navigates to the "Automations" section within their ACP project. They create a new automated workflow by specifying a repository and CI platform (GitHub Actions or GitLab CI). They optionally filter which branches and CI workflows/pipelines to monitor. They configure the action chain: start a debug session, create a PR on approval, and notify via Slack. Once saved, ACP begins polling the CI platform for failures at a configurable interval (default: 5 minutes).

**Why this priority**: This is the foundational capability. Without the ability to define what to watch and what to do, no other functionality works.

**Independent Test**: Can be fully tested by creating an automated workflow in the UI and verifying that the system begins polling the specified CI platform and detecting failures.

**Acceptance Scenarios**:

1. **Given** a user is in the Automations section, **When** they fill out the workflow form (repository, CI platform, branch filters, actions) and save, **Then** the system creates the workflow and begins polling at the configured interval.
2. **Given** an automated workflow exists, **When** the CI pipeline for the watched repository fails on a matching branch, **Then** the system detects the failure within one polling cycle.
3. **Given** an automated workflow with branch filters set to `["main"]`, **When** a CI pipeline fails on a feature branch, **Then** the system does not trigger any action.
4. **Given** a user has not configured CI credentials in ProjectSettings, **When** they try to create an automated workflow, **Then** the system prompts them to configure credentials first.

---

### User Story 2 - Automated Failure Analysis and Fix Proposal (Priority: P1)

When ACP detects a CI failure matching a workflow's filters, it automatically launches a two-stage debug process. First, a short-lived subagent analyzes the raw CI logs and distills them to the relevant errors, failing tests, and stack traces. Then, a main interactive debug session starts with the condensed logs, commit diff, and CI workflow configuration as context. The agent analyzes the root cause and proposes a fix, then pauses and waits for the user to review the analysis before making any code changes.

**Why this priority**: This is the core value proposition — automated debugging that saves developer time. Without this, the feature is just a CI status dashboard.

**Independent Test**: Can be tested by triggering a CI failure on a watched repository and verifying that a debug session is created with the correct context, the agent produces a root cause analysis, and the session pauses for user approval.

**Acceptance Scenarios**:

1. **Given** a CI failure is detected, **When** the system processes it, **Then** a log analysis subagent runs first and produces a condensed failure summary (errors, failing tests, stack traces only).
2. **Given** the log analysis subagent has completed, **When** the main debug session starts, **Then** it receives the condensed logs, the commit diff that triggered the failure, and the CI workflow/pipeline configuration file.
3. **Given** the debug session is running, **When** the agent completes its analysis, **Then** it presents a root cause analysis and proposed fix but does NOT apply code changes automatically.
4. **Given** the agent has proposed a fix, **When** the session pauses for approval, **Then** the user can review the analysis in the ACP session UI, ask follow-up questions, and approve or reject the proposed changes.

---

### User Story 3 - PR Creation on Approval (Priority: P2)

After reviewing the agent's analysis in the ACP session, the user approves the proposed fix. The agent creates a new branch (e.g., `acp/fix-ci-{run-id}`) and opens a pull request with the fix. The PR includes a description summarizing the CI failure, root cause, and the applied fix.

**Why this priority**: Delivering the fix as a PR is the natural conclusion of the debug workflow, but the analysis itself (P1) already provides significant value even without automated PR creation.

**Independent Test**: Can be tested by approving a fix proposal in an active debug session and verifying that a PR is created on the correct repository with the expected branch name and description.

**Acceptance Scenarios**:

1. **Given** a user approves the agent's proposed fix, **When** the agent applies the changes, **Then** it creates a new branch named `acp/fix-ci-{run-id}` and opens a pull request on the source repository.
2. **Given** the agent creates a PR, **When** the PR is opened, **Then** the description includes the CI failure summary, root cause analysis, and details of the applied fix.
3. **Given** the user rejects the proposed fix, **When** they provide feedback, **Then** the agent revises its approach based on the feedback and proposes an updated fix.

---

### User Story 4 - Slack Notifications (Priority: P2)

The user receives Slack notifications at key points in the automated workflow: when the agent's analysis is ready for review, and when a PR has been created. Each notification includes a brief summary and a link to the ACP session or PR.

**Why this priority**: Notifications ensure the user doesn't have to poll ACP to know when action is needed, but the core debugging functionality works without them.

**Independent Test**: Can be tested by configuring Slack notifications on a workflow, triggering a CI failure, and verifying that Slack messages arrive at the configured channel at the correct moments.

**Acceptance Scenarios**:

1. **Given** a workflow is configured with Slack notifications, **When** the agent's analysis is ready, **Then** a Slack message is sent to the configured channel with a failure summary and a link to the ACP session.
2. **Given** a workflow is configured with Slack notifications, **When** a PR is created, **Then** a Slack message is sent to the configured channel with the PR title and a link to the PR.
3. **Given** a workflow has no Slack notification configured, **When** a CI failure is processed, **Then** no Slack messages are sent (notifications are optional).

---

### User Story 5 - Manage and Monitor Automated Workflows (Priority: P2)

Users can list, edit, pause, resume, and delete their automated workflows from the ACP Automations section. A status dashboard shows recent trigger events, active debug sessions, and processed failures with their outcomes (session link, PR link, or error).

**Why this priority**: Ongoing management is essential for a production feature but is secondary to the core create-and-trigger flow.

**Independent Test**: Can be tested by creating multiple workflows, editing filters, pausing/resuming, and verifying the status dashboard reflects accurate state.

**Acceptance Scenarios**:

1. **Given** a user has created automated workflows, **When** they navigate to the Automations section, **Then** they see a list of all workflows with their current status (active, paused), last poll time, and number of triggered sessions.
2. **Given** an active workflow, **When** the user pauses it, **Then** polling stops and no new sessions are created for that workflow until it is resumed.
3. **Given** a workflow, **When** the user edits its filters (branches, workflows/pipelines), **Then** future polling uses the updated filters.
4. **Given** a workflow, **When** the user views its history, **Then** they see a list of processed failures with links to the corresponding debug sessions and PRs.

---

### User Story 6 - Deduplication of CI Failures (Priority: P3)

The system prevents duplicate debug sessions for the same CI failure. If a pipeline fails on the same commit and the system has already created a debug session for it, the failure is skipped. This prevents noise from flaky tests, manual re-runs, or repeated failures on the same broken commit.

**Why this priority**: Deduplication is important for a production-quality experience but the feature is usable (if noisier) without it.

**Independent Test**: Can be tested by triggering multiple CI failures on the same commit and verifying that only one debug session is created.

**Acceptance Scenarios**:

1. **Given** a CI failure was already processed for commit SHA `abc123` on pipeline `tests.yml`, **When** the same pipeline fails again on the same commit (e.g., manual re-run), **Then** no new debug session is created.
2. **Given** a CI failure was processed for commit `abc123`, **When** a new commit `def456` fails on the same pipeline, **Then** a new debug session IS created (different commit).
3. **Given** deduplication is active, **When** the user views the workflow history, **Then** skipped (deduplicated) failures are visible with a "skipped - duplicate" indicator.

---

### Edge Cases

- What happens when the CI platform API is unreachable during a poll cycle? The system logs the error, increments a failure counter, and retries on the next cycle. After 3 consecutive failures, the workflow status condition is updated to reflect the connectivity issue. After 10 consecutive failures, the workflow is automatically paused. The user must manually resume after fixing the issue.
- What happens when the raw CI log exceeds the subagent's context window? The log is chunked and the subagent processes it in segments, merging the extracted errors into a single condensed report.
- What happens when the user deletes a workflow while a debug session is active? The active session continues to completion but no new sessions are triggered.
- What happens when CI credentials are revoked or expire? The poller detects authentication failures and updates the workflow status to indicate invalid credentials. The user is notified in the UI. After 10 consecutive authentication failures, the workflow is automatically paused.
- What happens when multiple CI pipelines fail simultaneously on the same commit? Each pipeline failure is treated as a separate event and gets its own debug session (they have different pipeline IDs).
- What happens when the repository has been deleted or renamed? The poller detects the 404 response and marks the workflow as errored with a descriptive message.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to create automated workflows from the ACP UI, specifying a CI platform (GitHub Actions or GitLab CI), repository, optional branch filters, optional workflow/pipeline filters, polling interval, and an ordered list of actions.
- **FR-002**: System MUST poll configured CI platforms at the specified interval (default: 5 minutes) and detect pipeline/workflow failures matching the workflow's filters.
- **FR-003**: System MUST deduplicate failures using a composite key of commit SHA + pipeline/workflow ID, creating at most one debug session per unique failure.
- **FR-004**: System MUST launch a short-lived batch-mode subagent to analyze and condense raw CI logs before feeding them to the main debug session. The subagent extracts error messages, failing test names, stack traces, and relevant surrounding context.
- **FR-005**: System MUST create an interactive debug session with the condensed CI logs, the commit diff that triggered the failure, and the CI workflow/pipeline configuration file as context.
- **FR-006**: System MUST ensure the debug agent analyzes the root cause and proposes a fix, then pauses for user approval before making any code changes.
- **FR-007**: System MUST create a pull request on the source repository when the user approves the agent's proposed fix, using a branch name of `acp/fix-ci-{run-id}`.
- **FR-008**: System MUST send Slack notifications (when configured) at two points: when the analysis is ready for review, and when a PR is created. Each notification includes a summary and relevant link.
- **FR-009**: System MUST allow users to list, edit, pause, resume, and delete automated workflows from the ACP UI.
- **FR-010**: System MUST display a workflow history showing processed failures, their outcomes (session link, PR link), and any skipped duplicates.
- **FR-011**: System MUST store CI platform credentials (GitHub/GitLab API tokens) and Slack webhook URLs in the ProjectSettings, referencing secrets.
- **FR-012**: System MUST support both GitHub Actions and GitLab CI as CI platforms from the initial release.
- **FR-013**: System MUST handle CI platform API errors gracefully — logging failures, retrying on next cycle, and surfacing persistent errors in the workflow status. After 10 consecutive poll failures, the system MUST automatically pause the workflow. The user must manually resume after resolving the issue.
- **FR-014**: System MUST allow configurable polling intervals per workflow, with a minimum of 1 minute and a default of 5 minutes.
- **FR-015**: System MUST track workflow status conditions (Polling, CredentialsValid, LastPollSucceeded) following standard condition conventions.

### Key Entities

- **AutomatedWorkflow**: Represents a CI watch configuration. Contains the source definition (CI platform, repository, filters), polling schedule, deduplication settings, and an ordered action chain. Tracks processing history and current status.
- **ProcessedFailure**: A record of a CI failure that was detected and acted upon. Contains the commit SHA, pipeline/workflow ID, timestamp, and references to the created debug session and PR (if any). Used for deduplication. Records are retained for the lifetime of the parent AutomatedWorkflow and deleted when the workflow is deleted.
- **ActionChain**: An ordered list of actions to execute when a failure is detected. Each action has a type (start-session, create-pr, notify-slack) and action-specific configuration. Actions execute sequentially.
- **CIIntegration**: Credentials and configuration for connecting to a CI platform. Part of ProjectSettings. References a secret containing the API token.
- **SlackIntegration**: Configuration for Slack notifications. Part of ProjectSettings. References a secret containing the webhook URL.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can create and activate an automated CI watch workflow in under 3 minutes through the ACP UI.
- **SC-002**: CI failures are detected within one polling cycle (default: 5 minutes) of the failure occurring.
- **SC-003**: The log analysis subagent reduces raw CI logs to a condensed summary of 500 lines or fewer, retaining all error-relevant information.
- **SC-004**: The debug agent produces a root cause analysis and proposed fix for at least 70% of CI failures without requiring user guidance during the analysis phase.
- **SC-005**: Users receive Slack notifications within 60 seconds of the triggering event (analysis ready or PR created).
- **SC-006**: Duplicate CI failures (same commit + pipeline) never result in more than one debug session.
- **SC-007**: The system correctly handles CI platform API outages without data loss or duplicate processing — resumes normal operation when connectivity is restored.
- **SC-008**: Users can view the complete history of an automated workflow's triggers, sessions, and outcomes from the ACP UI.

## Assumptions

- Users already have ACP projects configured with Git repository access (existing ACP functionality).
- CI platform API tokens provided by users have sufficient permissions to read workflow/pipeline runs, job logs, and commit information.
- The Slack webhook URL provided by users is valid and has permission to post to the configured channel.
- The existing ACP session creation, interactive mode, and PR creation capabilities are functional and do not require modification for the core flow.
- The subagent (log analysis) can be implemented using the existing parent-child session mechanism already supported by ACP.
- Polling-based detection with a 5-minute default interval is acceptable latency for CI failure debugging (users do not require sub-minute reaction times).
- There is no per-project limit on the number of automated workflows. Resource consumption scales linearly with workflow count.

## Scope Boundaries

**In scope**:

- GitHub Actions and GitLab CI as CI platforms
- Pull-based (polling) CI failure detection
- Two-stage log analysis (subagent + main session)
- Configurable ordered action chains (start-session, create-pr, notify-slack)
- Deduplication by commit SHA + pipeline ID
- Workflow CRUD and status dashboard in ACP UI
- Slack notifications (notification + link format)

**Out of scope (deferred)**:

- Jira integration as an action type
- API rate limiting / adaptive throttling
- Automatic failure type classification
- REST API for AutomatedWorkflow management (UI only initially)
- Additional CI platforms (Jenkins, CircleCI, etc.)
- Bidirectional Slack interaction (reply/approve from Slack)
- Push-based (webhook) CI failure detection
