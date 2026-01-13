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
	UseProxy        bool
}

// HumioRepoAction represents a Humio repo action
type HumioRepoAction struct {
	IngestToken string
}

// OpsGenieAction represents an OpsGenie action
type OpsGenieAction struct {
	ApiUrl   string
	GenieKey string
	UseProxy bool
}

// PagerDutyAction represents a PagerDuty action
type PagerDutyAction struct {
	RoutingKey string
	Severity   string
	UseProxy   bool
}

// SlackFieldEntryInput represents a Slack field entry
type SlackFieldEntryInput struct {
	FieldName string
	Value     string
}

// SlackAction represents a Slack webhook action
type SlackAction struct {
	Url      string
	Fields   []SlackFieldEntryInput
	UseProxy bool
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
	UseProxy    bool
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
	IgnoreSSL    bool
	UseProxy     bool
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
        emailBodyTemplate: bodyTemplate
        emailUseProxy: useProxy
      }
      ... on HumioRepoAction {
        ingestToken
      }
      ... on OpsGenieAction {
        apiUrl
        genieKey
        opsGenieUseProxy: useProxy
      }
      ... on PagerDutyAction {
        routingKey
        severity
        pagerDutyUseProxy: useProxy
      }
      ... on SlackAction {
        url
        fields {
          fieldName
          value
        }
        slackUseProxy: useProxy
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
        victorOpsUseProxy: useProxy
      }
      ... on WebhookAction {
        method
        webhookUrl: url
        headers {
          header
          value
        }
        webhookBodyTemplate: bodyTemplate
        ignoreSSL
        webhookUseProxy: useProxy
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
  $UseProxy: Boolean!
) {
  createEmailAction(input: {
    viewName: $SearchDomainName
    name: $Name
    recipients: $Recipients
    subjectTemplate: $SubjectTemplate
    bodyTemplate: $BodyTemplate
    useProxy: $UseProxy
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
  $UseProxy: Boolean!
) {
  createOpsGenieAction(input: {
    viewName: $SearchDomainName
    name: $Name
    apiUrl: $ApiUrl
    genieKey: $GenieKey
    useProxy: $UseProxy
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
  $UseProxy: Boolean!
) {
  createPagerDutyAction(input: {
    viewName: $SearchDomainName
    name: $Name
    routingKey: $RoutingKey
    severity: $Severity
    useProxy: $UseProxy
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
  $UseProxy: Boolean!
) {
  createSlackAction(input: {
    viewName: $SearchDomainName
    name: $Name
    url: $Url
    fields: $Fields
    useProxy: $UseProxy
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
  $UseProxy: Boolean!
) {
  createVictorOpsAction(input: {
    viewName: $SearchDomainName
    name: $Name
    messageType: $MessageType
    notifyUrl: $NotifyUrl
    useProxy: $UseProxy
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
  $IgnoreSSL: Boolean!
  $UseProxy: Boolean!
) {
  createWebhookAction(input: {
    viewName: $SearchDomainName
    name: $Name
    url: $Url
    method: $Method
    headers: $Headers
    bodyTemplate: $BodyTemplate
    ignoreSSL: $IgnoreSSL
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updateEmailActionMutation = `
mutation UpdateEmailAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $Recipients: [String!]!
  $SubjectTemplate: String
  $BodyTemplate: String
  $UseProxy: Boolean!
) {
  updateEmailAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    recipients: $Recipients
    subjectTemplate: $SubjectTemplate
    bodyTemplate: $BodyTemplate
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updateHumioRepoActionMutation = `
mutation UpdateHumioRepoAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $IngestToken: String!
) {
  updateHumioRepoAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    ingestToken: $IngestToken
  }) {
    id
    name
  }
}
`

const updateOpsGenieActionMutation = `
mutation UpdateOpsGenieAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $ApiUrl: String!
  $GenieKey: String!
  $UseProxy: Boolean!
) {
  updateOpsGenieAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    apiUrl: $ApiUrl
    genieKey: $GenieKey
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updatePagerDutyActionMutation = `
mutation UpdatePagerDutyAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $RoutingKey: String!
  $Severity: String!
  $UseProxy: Boolean!
) {
  updatePagerDutyAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    routingKey: $RoutingKey
    severity: $Severity
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updateSlackActionMutation = `
mutation UpdateSlackAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $Url: String!
  $Fields: [SlackFieldEntryInput!]!
  $UseProxy: Boolean!
) {
  updateSlackAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    url: $Url
    fields: $Fields
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updateSlackPostMessageActionMutation = `
mutation UpdateSlackPostMessageAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $ApiToken: String!
  $Channels: [String!]!
  $Fields: [SlackFieldEntryInput!]!
  $UseProxy: Boolean!
) {
  updateSlackPostMessageAction(input: {
    viewName: $SearchDomainName
    id: $ID
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

const updateVictorOpsActionMutation = `
mutation UpdateVictorOpsAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $MessageType: String!
  $NotifyUrl: String!
  $UseProxy: Boolean!
) {
  updateVictorOpsAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    messageType: $MessageType
    notifyUrl: $NotifyUrl
    useProxy: $UseProxy
  }) {
    id
    name
  }
}
`

const updateWebhookActionMutation = `
mutation UpdateWebhookAction(
  $SearchDomainName: String!
  $ID: String!
  $Name: String!
  $Url: String!
  $Method: String!
  $Headers: [HttpHeaderEntryInput!]!
  $BodyTemplate: String!
  $IgnoreSSL: Boolean!
  $UseProxy: Boolean!
) {
  updateWebhookAction(input: {
    viewName: $SearchDomainName
    id: $ID
    name: $Name
    url: $Url
    method: $Method
    headers: $Headers
    bodyTemplate: $BodyTemplate
    ignoreSSL: $IgnoreSSL
    useProxy: $UseProxy
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
			Typename          string   `json:"__typename"`
			ID                string   `json:"id"`
			Name              string   `json:"name"`
			Recipients        []string `json:"recipients,omitempty"`
			SubjectTemplate   string   `json:"subjectTemplate,omitempty"`
			EmailBodyTemplate string   `json:"emailBodyTemplate,omitempty"`
			EmailUseProxy     bool     `json:"emailUseProxy,omitempty"`
			IngestToken       string   `json:"ingestToken,omitempty"`
			ApiUrl            string   `json:"apiUrl,omitempty"`
			GenieKey          string   `json:"genieKey,omitempty"`
			OpsGenieUseProxy  bool     `json:"opsGenieUseProxy,omitempty"`
			RoutingKey        string   `json:"routingKey,omitempty"`
			Severity          string   `json:"severity,omitempty"`
			PagerDutyUseProxy bool     `json:"pagerDutyUseProxy,omitempty"`
			Url               string   `json:"url,omitempty"`
			SlackUseProxy     bool     `json:"slackUseProxy,omitempty"`
			Fields            []struct {
				FieldName string `json:"fieldName"`
				Value     string `json:"value"`
			} `json:"fields,omitempty"`
			ApiToken            string   `json:"apiToken,omitempty"`
			Channels            []string `json:"channels,omitempty"`
			UseProxy            bool     `json:"useProxy,omitempty"`
			MessageType         string   `json:"messageType,omitempty"`
			NotifyUrl           string   `json:"notifyUrl,omitempty"`
			VictorOpsUseProxy   bool     `json:"victorOpsUseProxy,omitempty"`
			Method              string   `json:"method,omitempty"`
			WebhookUrl          string   `json:"webhookUrl,omitempty"`
			WebhookBodyTemplate string   `json:"webhookBodyTemplate,omitempty"`
			IgnoreSSL           bool     `json:"ignoreSSL,omitempty"`
			WebhookUseProxy     bool     `json:"webhookUseProxy,omitempty"`
			Headers             []struct {
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

type updateActionResponse struct {
	UpdateEmailAction            *struct{ ID, Name string } `json:"updateEmailAction,omitempty"`
	UpdateHumioRepoAction        *struct{ ID, Name string } `json:"updateHumioRepoAction,omitempty"`
	UpdateOpsGenieAction         *struct{ ID, Name string } `json:"updateOpsGenieAction,omitempty"`
	UpdatePagerDutyAction        *struct{ ID, Name string } `json:"updatePagerDutyAction,omitempty"`
	UpdateSlackAction            *struct{ ID, Name string } `json:"updateSlackAction,omitempty"`
	UpdateSlackPostMessageAction *struct{ ID, Name string } `json:"updateSlackPostMessageAction,omitempty"`
	UpdateVictorOpsAction        *struct{ ID, Name string } `json:"updateVictorOpsAction,omitempty"`
	UpdateWebhookAction          *struct{ ID, Name string } `json:"updateWebhookAction,omitempty"`
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
				BodyTemplate:    rawAction.EmailBodyTemplate,
				UseProxy:        rawAction.EmailUseProxy,
			}
		case ActionTypeHumioRepo:
			action.HumioRepoAction = HumioRepoAction{
				IngestToken: rawAction.IngestToken,
			}
		case ActionTypeOpsGenie:
			action.OpsGenieAction = OpsGenieAction{
				ApiUrl:   rawAction.ApiUrl,
				GenieKey: rawAction.GenieKey,
				UseProxy: rawAction.OpsGenieUseProxy,
			}
		case ActionTypePagerDuty:
			action.PagerDutyAction = PagerDutyAction{
				RoutingKey: rawAction.RoutingKey,
				Severity:   rawAction.Severity,
				UseProxy:   rawAction.PagerDutyUseProxy,
			}
		case ActionTypeSlack:
			fields := make([]SlackFieldEntryInput, len(rawAction.Fields))
			for j, f := range rawAction.Fields {
				fields[j] = SlackFieldEntryInput{FieldName: f.FieldName, Value: f.Value}
			}
			action.SlackAction = SlackAction{
				Url:      rawAction.Url,
				Fields:   fields,
				UseProxy: rawAction.SlackUseProxy,
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
				UseProxy:    rawAction.VictorOpsUseProxy,
			}
		case ActionTypeWebhook:
			headers := make([]HttpHeaderEntryInput, len(rawAction.Headers))
			for j, h := range rawAction.Headers {
				headers[j] = HttpHeaderEntryInput{Header: h.Header, Value: h.Value}
			}
			action.WebhookAction = WebhookAction{
				Method:       rawAction.Method,
				Url:          rawAction.WebhookUrl,
				Headers:      headers,
				BodyTemplate: rawAction.WebhookBodyTemplate,
				IgnoreSSL:    rawAction.IgnoreSSL,
				UseProxy:     rawAction.WebhookUseProxy,
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
		if action.EmailAction.SubjectTemplate != "" {
			variables["SubjectTemplate"] = action.EmailAction.SubjectTemplate
		}
		if action.EmailAction.BodyTemplate != "" {
			variables["BodyTemplate"] = action.EmailAction.BodyTemplate
		}
		variables["UseProxy"] = action.EmailAction.UseProxy

	case ActionTypeHumioRepo:
		mutation = createHumioRepoActionMutation
		variables["IngestToken"] = action.HumioRepoAction.IngestToken

	case ActionTypeOpsGenie:
		mutation = createOpsGenieActionMutation
		variables["ApiUrl"] = action.OpsGenieAction.ApiUrl
		variables["GenieKey"] = action.OpsGenieAction.GenieKey
		variables["UseProxy"] = action.OpsGenieAction.UseProxy

	case ActionTypePagerDuty:
		mutation = createPagerDutyActionMutation
		variables["RoutingKey"] = action.PagerDutyAction.RoutingKey
		variables["Severity"] = action.PagerDutyAction.Severity
		variables["UseProxy"] = action.PagerDutyAction.UseProxy

	case ActionTypeSlack:
		mutation = createSlackActionMutation
		variables["Url"] = action.SlackAction.Url
		fields := make([]map[string]string, len(action.SlackAction.Fields))
		for i, f := range action.SlackAction.Fields {
			fields[i] = map[string]string{"fieldName": f.FieldName, "value": f.Value}
		}
		variables["Fields"] = fields
		variables["UseProxy"] = action.SlackAction.UseProxy

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
		variables["UseProxy"] = action.VictorOpsAction.UseProxy

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
		variables["IgnoreSSL"] = action.WebhookAction.IgnoreSSL
		variables["UseProxy"] = action.WebhookAction.UseProxy

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

// Update updates an existing action in place
func (a *Actions) Update(repository string, action *Action) (*Action, error) {
	var resp updateActionResponse
	var mutation string
	variables := map[string]interface{}{
		"SearchDomainName": repository,
		"ID":               action.ID,
		"Name":             action.Name,
	}

	switch action.Type {
	case ActionTypeEmail:
		mutation = updateEmailActionMutation
		variables["Recipients"] = action.EmailAction.Recipients
		if action.EmailAction.SubjectTemplate != "" {
			variables["SubjectTemplate"] = action.EmailAction.SubjectTemplate
		}
		if action.EmailAction.BodyTemplate != "" {
			variables["BodyTemplate"] = action.EmailAction.BodyTemplate
		}
		variables["UseProxy"] = action.EmailAction.UseProxy

	case ActionTypeHumioRepo:
		mutation = updateHumioRepoActionMutation
		variables["IngestToken"] = action.HumioRepoAction.IngestToken

	case ActionTypeOpsGenie:
		mutation = updateOpsGenieActionMutation
		variables["ApiUrl"] = action.OpsGenieAction.ApiUrl
		variables["GenieKey"] = action.OpsGenieAction.GenieKey
		variables["UseProxy"] = action.OpsGenieAction.UseProxy

	case ActionTypePagerDuty:
		mutation = updatePagerDutyActionMutation
		variables["RoutingKey"] = action.PagerDutyAction.RoutingKey
		variables["Severity"] = action.PagerDutyAction.Severity
		variables["UseProxy"] = action.PagerDutyAction.UseProxy

	case ActionTypeSlack:
		mutation = updateSlackActionMutation
		variables["Url"] = action.SlackAction.Url
		fields := make([]map[string]string, len(action.SlackAction.Fields))
		for i, f := range action.SlackAction.Fields {
			fields[i] = map[string]string{"fieldName": f.FieldName, "value": f.Value}
		}
		variables["Fields"] = fields
		variables["UseProxy"] = action.SlackAction.UseProxy

	case ActionTypeSlackPostMessage:
		mutation = updateSlackPostMessageActionMutation
		variables["ApiToken"] = action.SlackPostMessageAction.ApiToken
		variables["Channels"] = action.SlackPostMessageAction.Channels
		fields := make([]map[string]string, len(action.SlackPostMessageAction.Fields))
		for i, f := range action.SlackPostMessageAction.Fields {
			fields[i] = map[string]string{"fieldName": f.FieldName, "value": f.Value}
		}
		variables["Fields"] = fields
		variables["UseProxy"] = action.SlackPostMessageAction.UseProxy

	case ActionTypeVictorOps:
		mutation = updateVictorOpsActionMutation
		variables["MessageType"] = action.VictorOpsAction.MessageType
		variables["NotifyUrl"] = action.VictorOpsAction.NotifyUrl
		variables["UseProxy"] = action.VictorOpsAction.UseProxy

	case ActionTypeWebhook:
		mutation = updateWebhookActionMutation
		variables["Url"] = action.WebhookAction.Url
		variables["Method"] = action.WebhookAction.Method
		headers := make([]map[string]string, len(action.WebhookAction.Headers))
		for i, h := range action.WebhookAction.Headers {
			headers[i] = map[string]string{"header": h.Header, "value": h.Value}
		}
		variables["Headers"] = headers
		variables["BodyTemplate"] = action.WebhookAction.BodyTemplate
		variables["IgnoreSSL"] = action.WebhookAction.IgnoreSSL
		variables["UseProxy"] = action.WebhookAction.UseProxy

	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}

	err := a.client.Query(context.Background(), mutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	return action, nil
}

// Delete deletes an action by ID
func (a *Actions) Delete(repository, actionID string) error {
	return a.client.Query(context.Background(), deleteActionMutation, map[string]interface{}{
		"SearchDomainName": repository,
		"ActionID":         actionID,
	}, nil)
}
