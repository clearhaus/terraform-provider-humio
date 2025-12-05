package api

import (
	"context"
	"fmt"
)

// Action type constants
const (
	ActionTypeEmail            = "EmailAction"
	ActionTypeHumioRepo        = "HumioRepoAction"
	ActionTypeOpsGenie         = "OpsGenieAction"
	ActionTypePagerDuty        = "PagerDutyAction"
	ActionTypeSlack            = "SlackAction"
	ActionTypeSlackPostMessage = "SlackPostMessageAction"
	ActionTypeVictorOps        = "VictorOpsAction"
	ActionTypeWebhook          = "WebhookAction"
)

// Action represents a Humio action
type Action struct {
	Type                   string
	ID                     string
	Name                   string
	EmailAction            EmailAction
	HumioRepoAction        HumioRepoAction
	OpsGenieAction         OpsGenieAction
	PagerDutyAction        PagerDutyAction
	SlackAction            SlackAction
	SlackPostMessageAction SlackPostMessageAction
	VictorOpsAction        VictorOpsAction
	WebhookAction          WebhookAction
}

// EmailAction represents an email action
type EmailAction struct {
	Recipients      []string
	SubjectTemplate string
	BodyTemplate    string
}

// HumioRepoAction represents a Humio repo action
type HumioRepoAction struct {
	IngestToken string
}

// OpsGenieAction represents an OpsGenie action
type OpsGenieAction struct {
	ApiUrl   string
	GenieKey string
}

// PagerDutyAction represents a PagerDuty action
type PagerDutyAction struct {
	RoutingKey string
	Severity   string
}

// SlackFieldEntryInput represents a Slack field entry
type SlackFieldEntryInput struct {
	FieldName string
	Value     string
}

// SlackAction represents a Slack webhook action
type SlackAction struct {
	Url    string
	Fields []SlackFieldEntryInput
}

// SlackPostMessageAction represents a Slack post message action
type SlackPostMessageAction struct {
	ApiToken string
	Channels []string
	Fields   []SlackFieldEntryInput
	UseProxy bool
}

// VictorOpsAction represents a VictorOps action
type VictorOpsAction struct {
	MessageType string
	NotifyUrl   string
}

// HttpHeaderEntryInput represents an HTTP header entry
type HttpHeaderEntryInput struct {
	Header string
	Value  string
}

// WebhookAction represents a webhook action
type WebhookAction struct {
	Method       string
	Url          string
	Headers      []HttpHeaderEntryInput
	BodyTemplate string
}

// Actions provides operations for managing actions
type Actions struct {
	client *Client
}

const listActionsQuery = `
query ListActions($SearchDomainName: String!) {
  searchDomain(name: $SearchDomainName) {
    actions {
      __typename
      id
      name
      ... on EmailAction {
        recipients
        subjectTemplate
        bodyTemplate
      }
      ... on HumioRepoAction {
        ingestToken
      }
      ... on OpsGenieAction {
        apiUrl
        genieKey
      }
      ... on PagerDutyAction {
        routingKey
        severity
      }
      ... on SlackAction {
        url
        fields {
          fieldName
          value
        }
      }
      ... on SlackPostMessageAction {
        apiToken
        channels
        fields {
          fieldName
          value
        }
        useProxy
      }
      ... on VictorOpsAction {
        messageType
        notifyUrl
      }
      ... on WebhookAction {
        method
        url
        headers {
          header
          value
        }
        bodyTemplate
      }
    }
  }
}
`

const deleteActionMutation = `
mutation DeleteAction($SearchDomainName: String!, $ActionID: String!) {
  deleteAction(input: {
    viewName: $SearchDomainName
    id: $ActionID
  })
}
`

const createEmailActionMutation = `
mutation CreateEmailAction(
  $SearchDomainName: String!
  $Name: String!
  $Recipients: [String!]!
  $SubjectTemplate: String
  $BodyTemplate: String
) {
  createEmailAction(input: {
    viewName: $SearchDomainName
    name: $Name
    recipients: $Recipients
    subjectTemplate: $SubjectTemplate
    bodyTemplate: $BodyTemplate
  }) {
    id
    name
  }
}
`

const createHumioRepoActionMutation = `
mutation CreateHumioRepoAction(
  $SearchDomainName: String!
  $Name: String!
  $IngestToken: String!
) {
  createHumioRepoAction(input: {
    viewName: $SearchDomainName
    name: $Name
    ingestToken: $IngestToken
  }) {
    id
    name
  }
}
`

const createOpsGenieActionMutation = `
mutation CreateOpsGenieAction(
  $SearchDomainName: String!
  $Name: String!
  $ApiUrl: String!
  $GenieKey: String!
) {
  createOpsGenieAction(input: {
    viewName: $SearchDomainName
    name: $Name
    apiUrl: $ApiUrl
    genieKey: $GenieKey
  }) {
    id
    name
  }
}
`

const createPagerDutyActionMutation = `
mutation CreatePagerDutyAction(
  $SearchDomainName: String!
  $Name: String!
  $RoutingKey: String!
  $Severity: String!
) {
  createPagerDutyAction(input: {
    viewName: $SearchDomainName
    name: $Name
    routingKey: $RoutingKey
    severity: $Severity
  }) {
    id
    name
  }
}
`

const createSlackActionMutation = `
mutation CreateSlackAction(
  $SearchDomainName: String!
  $Name: String!
  $Url: String!
  $Fields: [SlackFieldEntryInput!]!
) {
  createSlackAction(input: {
    viewName: $SearchDomainName
    name: $Name
    url: $Url
    fields: $Fields
  }) {
    id
    name
  }
}
`

const createSlackPostMessageActionMutation = `
mutation CreateSlackPostMessageAction(
  $SearchDomainName: String!
  $Name: String!
  $ApiToken: String!
  $Channels: [String!]!
  $Fields: [SlackFieldEntryInput!]!
  $UseProxy: Boolean!
) {
  createSlackPostMessageAction(input: {
    viewName: $SearchDomainName
    name: $Name
    apiToken: $ApiToken
    channels: $Channels
    fields: $Fields
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const createVictorOpsActionMutation = `
mutation CreateVictorOpsAction(
  $SearchDomainName: String!
  $Name: String!
  $MessageType: String!
  $NotifyUrl: String!
) {
  createVictorOpsAction(input: {
    viewName: $SearchDomainName
    name: $Name
    messageType: $MessageType
    notifyUrl: $NotifyUrl
  }) {
    id
    name
  }
}
`

const createWebhookActionMutation = `
mutation CreateWebhookAction(
  $SearchDomainName: String!
  $Name: String!
  $Url: String!
  $Method: String!
  $Headers: [HttpHeaderEntryInput!]!
  $BodyTemplate: String!
) {
  createWebhookAction(input: {
    viewName: $SearchDomainName
    name: $Name
    url: $Url
    method: $Method
    headers: $Headers
    bodyTemplate: $BodyTemplate
  }) {
    id
    name
  }
}
`

// actionResponse represents the response from list actions query
type actionResponse struct {
	SearchDomain struct {
		Actions []struct {
			Typename        string   `json:"__typename"`
			ID              string   `json:"id"`
			Name            string   `json:"name"`
			Recipients      []string `json:"recipients,omitempty"`
			SubjectTemplate string   `json:"subjectTemplate,omitempty"`
			BodyTemplate    string   `json:"bodyTemplate,omitempty"`
			IngestToken     string   `json:"ingestToken,omitempty"`
			ApiUrl          string   `json:"apiUrl,omitempty"`
			GenieKey        string   `json:"genieKey,omitempty"`
			RoutingKey      string   `json:"routingKey,omitempty"`
			Severity        string   `json:"severity,omitempty"`
			Url             string   `json:"url,omitempty"`
			Fields          []struct {
				FieldName string `json:"fieldName"`
				Value     string `json:"value"`
			} `json:"fields,omitempty"`
			ApiToken    string   `json:"apiToken,omitempty"`
			Channels    []string `json:"channels,omitempty"`
			UseProxy    bool     `json:"useProxy,omitempty"`
			MessageType string   `json:"messageType,omitempty"`
			NotifyUrl   string   `json:"notifyUrl,omitempty"`
			Method      string   `json:"method,omitempty"`
			Headers     []struct {
				Header string `json:"header"`
				Value  string `json:"value"`
			} `json:"headers,omitempty"`
		} `json:"actions"`
	} `json:"searchDomain"`
}

type createActionResponse struct {
	CreateEmailAction            *struct{ ID, Name string } `json:"createEmailAction,omitempty"`
	CreateHumioRepoAction        *struct{ ID, Name string } `json:"createHumioRepoAction,omitempty"`
	CreateOpsGenieAction         *struct{ ID, Name string } `json:"createOpsGenieAction,omitempty"`
	CreatePagerDutyAction        *struct{ ID, Name string } `json:"createPagerDutyAction,omitempty"`
	CreateSlackAction            *struct{ ID, Name string } `json:"createSlackAction,omitempty"`
	CreateSlackPostMessageAction *struct{ ID, Name string } `json:"createSlackPostMessageAction,omitempty"`
	CreateVictorOpsAction        *struct{ ID, Name string } `json:"createVictorOpsAction,omitempty"`
	CreateWebhookAction          *struct{ ID, Name string } `json:"createWebhookAction,omitempty"`
}

// List returns all actions for the given repository
func (a *Actions) List(repository string) ([]Action, error) {
	var resp actionResponse
	err := a.client.Query(context.Background(), listActionsQuery, map[string]interface{}{
		"SearchDomainName": repository,
	}, &resp)
	if err != nil {
		return nil, err
	}

	actions := make([]Action, len(resp.SearchDomain.Actions))
	for i, rawAction := range resp.SearchDomain.Actions {
		action := Action{
			Type: rawAction.Typename,
			ID:   rawAction.ID,
			Name: rawAction.Name,
		}

		switch rawAction.Typename {
		case ActionTypeEmail:
			action.EmailAction = EmailAction{
				Recipients:      rawAction.Recipients,
				SubjectTemplate: rawAction.SubjectTemplate,
				BodyTemplate:    rawAction.BodyTemplate,
			}
		case ActionTypeHumioRepo:
			action.HumioRepoAction = HumioRepoAction{
				IngestToken: rawAction.IngestToken,
			}
		case ActionTypeOpsGenie:
			action.OpsGenieAction = OpsGenieAction{
				ApiUrl:   rawAction.ApiUrl,
				GenieKey: rawAction.GenieKey,
			}
		case ActionTypePagerDuty:
			action.PagerDutyAction = PagerDutyAction{
				RoutingKey: rawAction.RoutingKey,
				Severity:   rawAction.Severity,
			}
		case ActionTypeSlack:
			fields := make([]SlackFieldEntryInput, len(rawAction.Fields))
			for j, f := range rawAction.Fields {
				fields[j] = SlackFieldEntryInput{FieldName: f.FieldName, Value: f.Value}
			}
			action.SlackAction = SlackAction{
				Url:    rawAction.Url,
				Fields: fields,
			}
		case ActionTypeSlackPostMessage:
			fields := make([]SlackFieldEntryInput, len(rawAction.Fields))
			for j, f := range rawAction.Fields {
				fields[j] = SlackFieldEntryInput{FieldName: f.FieldName, Value: f.Value}
			}
			action.SlackPostMessageAction = SlackPostMessageAction{
				ApiToken: rawAction.ApiToken,
				Channels: rawAction.Channels,
				Fields:   fields,
				UseProxy: rawAction.UseProxy,
			}
		case ActionTypeVictorOps:
			action.VictorOpsAction = VictorOpsAction{
				MessageType: rawAction.MessageType,
				NotifyUrl:   rawAction.NotifyUrl,
			}
		case ActionTypeWebhook:
			headers := make([]HttpHeaderEntryInput, len(rawAction.Headers))
			for j, h := range rawAction.Headers {
				headers[j] = HttpHeaderEntryInput{Header: h.Header, Value: h.Value}
			}
			action.WebhookAction = WebhookAction{
				Method:       rawAction.Method,
				Url:          rawAction.Url,
				Headers:      headers,
				BodyTemplate: rawAction.BodyTemplate,
			}
		}

		actions[i] = action
	}

	return actions, nil
}

// Get returns an action by name
func (a *Actions) Get(repository, name string) (*Action, error) {
	actions, err := a.List(repository)
	if err != nil {
		return nil, err
	}

	for _, action := range actions {
		if action.Name == name {
			return &action, nil
		}
	}

	return nil, fmt.Errorf("action not found: %s", name)
}

// Add creates a new action
func (a *Actions) Add(repository string, action *Action) (*Action, error) {
	var resp createActionResponse
	var mutation string
	variables := map[string]interface{}{
		"SearchDomainName": repository,
		"Name":             action.Name,
	}

	switch action.Type {
	case ActionTypeEmail:
		mutation = createEmailActionMutation
		variables["Recipients"] = action.EmailAction.Recipients
		variables["SubjectTemplate"] = action.EmailAction.SubjectTemplate
		variables["BodyTemplate"] = action.EmailAction.BodyTemplate

	case ActionTypeHumioRepo:
		mutation = createHumioRepoActionMutation
		variables["IngestToken"] = action.HumioRepoAction.IngestToken

	case ActionTypeOpsGenie:
		mutation = createOpsGenieActionMutation
		variables["ApiUrl"] = action.OpsGenieAction.ApiUrl
		variables["GenieKey"] = action.OpsGenieAction.GenieKey

	case ActionTypePagerDuty:
		mutation = createPagerDutyActionMutation
		variables["RoutingKey"] = action.PagerDutyAction.RoutingKey
		variables["Severity"] = action.PagerDutyAction.Severity

	case ActionTypeSlack:
		mutation = createSlackActionMutation
		variables["Url"] = action.SlackAction.Url
		fields := make([]map[string]string, len(action.SlackAction.Fields))
		for i, f := range action.SlackAction.Fields {
			fields[i] = map[string]string{"fieldName": f.FieldName, "value": f.Value}
		}
		variables["Fields"] = fields

	case ActionTypeSlackPostMessage:
		mutation = createSlackPostMessageActionMutation
		variables["ApiToken"] = action.SlackPostMessageAction.ApiToken
		variables["Channels"] = action.SlackPostMessageAction.Channels
		fields := make([]map[string]string, len(action.SlackPostMessageAction.Fields))
		for i, f := range action.SlackPostMessageAction.Fields {
			fields[i] = map[string]string{"fieldName": f.FieldName, "value": f.Value}
		}
		variables["Fields"] = fields
		variables["UseProxy"] = action.SlackPostMessageAction.UseProxy

	case ActionTypeVictorOps:
		mutation = createVictorOpsActionMutation
		variables["MessageType"] = action.VictorOpsAction.MessageType
		variables["NotifyUrl"] = action.VictorOpsAction.NotifyUrl

	case ActionTypeWebhook:
		mutation = createWebhookActionMutation
		variables["Url"] = action.WebhookAction.Url
		variables["Method"] = action.WebhookAction.Method
		headers := make([]map[string]string, len(action.WebhookAction.Headers))
		for i, h := range action.WebhookAction.Headers {
			headers[i] = map[string]string{"header": h.Header, "value": h.Value}
		}
		variables["Headers"] = headers
		variables["BodyTemplate"] = action.WebhookAction.BodyTemplate

	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}

	err := a.client.Query(context.Background(), mutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	// Extract the ID from the response
	switch action.Type {
	case ActionTypeEmail:
		if resp.CreateEmailAction != nil {
			action.ID = resp.CreateEmailAction.ID
		}
	case ActionTypeHumioRepo:
		if resp.CreateHumioRepoAction != nil {
			action.ID = resp.CreateHumioRepoAction.ID
		}
	case ActionTypeOpsGenie:
		if resp.CreateOpsGenieAction != nil {
			action.ID = resp.CreateOpsGenieAction.ID
		}
	case ActionTypePagerDuty:
		if resp.CreatePagerDutyAction != nil {
			action.ID = resp.CreatePagerDutyAction.ID
		}
	case ActionTypeSlack:
		if resp.CreateSlackAction != nil {
			action.ID = resp.CreateSlackAction.ID
		}
	case ActionTypeSlackPostMessage:
		if resp.CreateSlackPostMessageAction != nil {
			action.ID = resp.CreateSlackPostMessageAction.ID
		}
	case ActionTypeVictorOps:
		if resp.CreateVictorOpsAction != nil {
			action.ID = resp.CreateVictorOpsAction.ID
		}
	case ActionTypeWebhook:
		if resp.CreateWebhookAction != nil {
			action.ID = resp.CreateWebhookAction.ID
		}
	}

	return action, nil
}

// Update updates an existing action by deleting and recreating it
func (a *Actions) Update(repository string, action *Action) (*Action, error) {
	// First get the existing action to get its ID
	existing, err := a.Get(repository, action.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing action: %w", err)
	}

	// Delete the existing action
	if err := a.Delete(repository, existing.ID); err != nil {
		return nil, fmt.Errorf("failed to delete existing action: %w", err)
	}

	// Create the new action
	return a.Add(repository, action)
}

// Delete deletes an action by name or ID
func (a *Actions) Delete(repository, actionNameOrID string) error {
	actionID := actionNameOrID

	// If it doesn't look like an ID, try to find the action by name
	if len(actionNameOrID) < 20 {
		action, err := a.Get(repository, actionNameOrID)
		if err != nil {
			return err
		}
		actionID = action.ID
	}

	return a.client.Query(context.Background(), deleteActionMutation, map[string]interface{}{
		"SearchDomainName": repository,
		"ActionID":         actionID,
	}, nil)
}
