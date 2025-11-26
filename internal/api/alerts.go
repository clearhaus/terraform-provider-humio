// Copyright © 2024 Clearhaus
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
        id
        ... on QueryOwnershipTypeUser {
          user {
            id
          }
        }
        ... on QueryOwnershipTypeOrganization {
          id
        }
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
				ID   string `json:"id"`
				User *struct {
					ID string `json:"id"`
				} `json:"user,omitempty"`
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
		if alert.QueryOwnership.User != nil {
			ownershipType = "User"
			runAsUserID = alert.QueryOwnership.User.ID
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
	variables := map[string]interface{}{
		"SearchDomainName":   repository,
		"Name":               alert.Name,
		"Description":        alert.Description,
		"QueryString":        alert.QueryString,
		"QueryStart":         alert.QueryStart,
		"ThrottleTimeMillis": alert.ThrottleTimeMillis,
		"ThrottleField":      alert.ThrottleField,
		"Enabled":            alert.Enabled,
		"Actions":            alert.Actions,
		"Labels":             alert.Labels,
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
	// First get the existing alert to get its ID
	existing, err := a.Get(repository, alert.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing alert: %w", err)
	}

	// Delete the existing alert
	if err := a.Delete(repository, existing.ID); err != nil {
		return nil, fmt.Errorf("failed to delete existing alert: %w", err)
	}

	// Create the new alert
	return a.Add(repository, alert)
}

// Delete deletes an alert by ID
func (a *Alerts) Delete(repository, alertID string) error {
	// If alertID looks like a name, get the actual ID first
	if len(alertID) < 20 { // IDs are typically UUIDs
		alert, err := a.Get(repository, alertID)
		if err != nil {
			return err
		}
		alertID = alert.ID
	}

	return a.client.Query(context.Background(), deleteAlertMutation, map[string]interface{}{
		"SearchDomainName": repository,
		"AlertID":          alertID,
	}, nil)
}
