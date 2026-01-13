package api

import (
	"context"
)

// Repository represents a Humio repository
type Repository struct {
	ID            string
	Name          string
	Description   string
	RetentionDays float64
}

// Repositories provides operations for managing repositories
type Repositories struct {
	client *Client
}

const getRepositoryQuery = `
query GetRepository($RepositoryName: String!) {
  repository(name: $RepositoryName) {
    id
    name
    description
    timeBasedRetention
  }
}
`

const createRepositoryMutation = `
mutation CreateRepository($Name: String!) {
  createRepository(name: $Name) {
    __typename
  }
}
`

const updateDescriptionMutation = `
mutation UpdateDescription($RepositoryName: String!, $Description: String!) {
  updateDescriptionForSearchDomain(name: $RepositoryName, newDescription: $Description) {
    __typename
  }
}
`

const updateTimeBasedRetentionMutation = `
mutation UpdateTimeBasedRetention($RepositoryName: String!, $RetentionDays: Float) {
  updateRetention(
    repositoryName: $RepositoryName
    timeBasedRetention: $RetentionDays
  ) {
    repository {
      id
      name
      ... on Repository {
        timeBasedRetention
      }
    }
  }
}
`

const deleteRepositoryMutation = `
mutation DeleteRepository($RepositoryName: String!, $Reason: String) {
  deleteSearchDomain(name: $RepositoryName, deleteMessage: $Reason) {
    result
  }
}
`

const listRepositoriesQuery = `
query ListRepositories {
  repositories {
    id
    name
  }
}
`

// listRepositoriesResponse represents the response from list repositories query
type listRepositoriesResponse struct {
	Repositories []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"repositories"`
}

// getRepositoryResponse represents the response from get repository query
type getRepositoryResponse struct {
	Repository struct {
		ID                 string   `json:"id"`
		Name               string   `json:"name"`
		Description        string   `json:"description"`
		TimeBasedRetention *float64 `json:"timeBasedRetention"`
	} `json:"repository"`
}

// List returns all repositories
func (r *Repositories) List() ([]Repository, error) {
	var resp listRepositoriesResponse
	err := r.client.Query(context.Background(), listRepositoriesQuery, nil, &resp)
	if err != nil {
		return nil, err
	}

	repositories := make([]Repository, len(resp.Repositories))
	for i, repo := range resp.Repositories {
		repositories[i] = Repository{
			ID:   repo.ID,
			Name: repo.Name,
		}
	}
	return repositories, nil
}

// Get returns a repository by name
func (r *Repositories) Get(name string) (Repository, error) {
	var resp getRepositoryResponse
	err := r.client.Query(context.Background(), getRepositoryQuery, map[string]interface{}{
		"RepositoryName": name,
	}, &resp)
	if err != nil {
		return Repository{}, err
	}

	retentionDays := 0.0
	if resp.Repository.TimeBasedRetention != nil {
		retentionDays = *resp.Repository.TimeBasedRetention
	}

	return Repository{
		ID:            resp.Repository.ID,
		Name:          resp.Repository.Name,
		Description:   resp.Repository.Description,
		RetentionDays: retentionDays,
	}, nil
}

// Create creates a new repository
func (r *Repositories) Create(name string) error {
	return r.client.Query(context.Background(), createRepositoryMutation, map[string]interface{}{
		"Name": name,
	}, nil)
}

// UpdateDescription updates the description of a repository
func (r *Repositories) UpdateDescription(name, description string) error {
	return r.client.Query(context.Background(), updateDescriptionMutation, map[string]interface{}{
		"RepositoryName": name,
		"Description":    description,
	}, nil)
}

// UpdateTimeBasedRetention updates the time-based retention for a repository
func (r *Repositories) UpdateTimeBasedRetention(name string, retentionDays float64) error {
	variables := map[string]interface{}{
		"RepositoryName": name,
	}

	// Only set retention if it's > 0, otherwise null means unlimited
	if retentionDays > 0 {
		variables["RetentionDays"] = retentionDays
	}

	return r.client.Query(context.Background(), updateTimeBasedRetentionMutation, variables, nil)
}

// Delete deletes a repository
func (r *Repositories) Delete(name, reason string) error {
	return r.client.Query(context.Background(), deleteRepositoryMutation, map[string]interface{}{
		"RepositoryName": name,
		"Reason":         reason,
	}, nil)
}
