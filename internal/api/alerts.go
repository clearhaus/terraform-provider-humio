package api

import (
	"context"
	"fmt"
)

// Alert represents a Humio alert
type Alert struct {
	ID                 string
	Name               string
	Description        string
	QueryString        string
	QueryStart         string
	ThrottleField      string
	ThrottleTimeMillis int
	Enabled            bool
	Actions            []string
	Labels             []string
	RunAsUserID        string
	QueryOwnershipType string
}

// Alerts provides operations for managing alerts
type Alerts struct {
	client *Client
}

const listAlertsQuery = `
query ListAlerts($SearchDomainName: String!) {
  searchDomain(name: $SearchDomainName) {
    alerts {
      id
      name
      description
      queryString
      queryStart
      throttleField
      throttleTimeMillis
      enabled
      actions
      labels
      queryOwnership {
        __typename
        id
      }
    }
  }
}
`

const createAlertMutation = `
mutation CreateAlert(
  $SearchDomainName: String!
  $Name: String!
  $Description: String
  $QueryString: String!
  $QueryStart: String!
  $ThrottleTimeMillis: Long!
  $ThrottleField: String
  $Enabled: Boolean!
  $Actions: [String!]!
  $Labels: [String!]
  $RunAsUserID: String
  $QueryOwnershipType: QueryOwnershipType
) {
  createAlert(input: {
    viewName: $SearchDomainName
    name: $Name
    description: $Description
    queryString: $QueryString
    queryStart: $QueryStart
    throttleTimeMillis: $ThrottleTimeMillis
    throttleField: $ThrottleField
    enabled: $Enabled
    actions: $Actions
    labels: $Labels
    runAsUserId: $RunAsUserID
    queryOwnershipType: $QueryOwnershipType
  }) {
    id
    name
  }
}
`

const deleteAlertMutation = `
mutation DeleteAlert($SearchDomainName: String!, $AlertID: String!) {
  deleteAlert(input: {
    viewName: $SearchDomainName
    id: $AlertID
  })
}
`

// alertResponse is the response structure for alert queries
type alertResponse struct {
	SearchDomain struct {
		Alerts []struct {
			ID                 string   `json:"id"`
			Name               string   `json:"name"`
			Description        string   `json:"description"`
			QueryString        string   `json:"queryString"`
			QueryStart         string   `json:"queryStart"`
			ThrottleField      string   `json:"throttleField"`
			ThrottleTimeMillis int      `json:"throttleTimeMillis"`
			Enabled            bool     `json:"enabled"`
			Actions            []string `json:"actions"`
			Labels             []string `json:"labels"`
			QueryOwnership     struct {
				Typename string `json:"__typename"`
				ID       string `json:"id"`
			} `json:"queryOwnership"`
		} `json:"alerts"`
	} `json:"searchDomain"`
}

// createAlertResponse is the response structure for creating an alert
type createAlertResponse struct {
	CreateAlert struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createAlert"`
}

// List returns all alerts for the given repository
func (a *Alerts) List(repository string) ([]Alert, error) {
	var resp alertResponse
	err := a.client.Query(context.Background(), listAlertsQuery, map[string]interface{}{
		"SearchDomainName": repository,
	}, &resp)
	if err != nil {
		return nil, err
	}

	alerts := make([]Alert, len(resp.SearchDomain.Alerts))
	for i, alert := range resp.SearchDomain.Alerts {
		ownershipType := "Organization"
		runAsUserID := ""
		if alert.QueryOwnership.Typename == "UserOwnership" {
			ownershipType = "User"
			runAsUserID = alert.QueryOwnership.ID
		}

		alerts[i] = Alert{
			ID:                 alert.ID,
			Name:               alert.Name,
			Description:        alert.Description,
			QueryString:        alert.QueryString,
			QueryStart:         alert.QueryStart,
			ThrottleField:      alert.ThrottleField,
			ThrottleTimeMillis: alert.ThrottleTimeMillis,
			Enabled:            alert.Enabled,
			Actions:            alert.Actions,
			Labels:             alert.Labels,
			RunAsUserID:        runAsUserID,
			QueryOwnershipType: ownershipType,
		}
	}

	return alerts, nil
}

// Get returns an alert by name
func (a *Alerts) Get(repository, name string) (*Alert, error) {
	alerts, err := a.List(repository)
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		if alert.Name == name {
			return &alert, nil
		}
	}

	return nil, fmt.Errorf("alert not found: %s", name)
}

// Add creates a new alert
func (a *Alerts) Add(repository string, alert *Alert) (*Alert, error) {
	actions := alert.Actions
	if actions == nil {
		actions = []string{}
	}
	labels := alert.Labels
	if labels == nil {
		labels = []string{}
	}

	variables := map[string]interface{}{
		"SearchDomainName":   repository,
		"Name":               alert.Name,
		"Description":        alert.Description,
		"QueryString":        alert.QueryString,
		"QueryStart":         alert.QueryStart,
		"ThrottleTimeMillis": alert.ThrottleTimeMillis,
		"Enabled":            alert.Enabled,
		"Actions":            actions,
		"Labels":             labels,
	}

	if alert.ThrottleField != "" {
		variables["ThrottleField"] = alert.ThrottleField
	}
	if alert.RunAsUserID != "" {
		variables["RunAsUserID"] = alert.RunAsUserID
	}
	if alert.QueryOwnershipType != "" {
		variables["QueryOwnershipType"] = alert.QueryOwnershipType
	}

	var resp createAlertResponse
	err := a.client.Query(context.Background(), createAlertMutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	alert.ID = resp.CreateAlert.ID
	return alert, nil
}

// Update updates an existing alert by deleting and recreating it
func (a *Alerts) Update(repository string, alert *Alert) (*Alert, error) {
	// Delete the existing alert
	if err := a.Delete(repository, alert.Name); err != nil {
		return nil, fmt.Errorf("failed to delete existing alert: %w", err)
	}

	// Create the new alert
	return a.Add(repository, alert)
}

// Delete deletes an alert by name
func (a *Alerts) Delete(repository, alertName string) error {
	// Look up the alert by name to get its ID
	alert, err := a.Get(repository, alertName)
	if err != nil {
		return err
	}

	return a.client.Query(context.Background(), deleteAlertMutation, map[string]interface{}{
		"SearchDomainName": repository,
		"AlertID":          alert.ID,
	}, nil)
}
