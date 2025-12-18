package api

import (
	"context"
	"fmt"
)

// IngestToken represents a Humio ingest token
type IngestToken struct {
	Name           string
	Token          string
	AssignedParser string
}

// IngestTokens provides operations for managing ingest tokens
type IngestTokens struct {
	client *Client
}

const listIngestTokensQuery = `
query ListIngestTokens($RepositoryName: String!) {
  repository(name: $RepositoryName) {
    ingestTokens {
      name
      token
      parser {
        name
      }
    }
  }
}
`

const addIngestTokenMutation = `
mutation AddIngestToken($RepositoryName: String!, $Name: String!, $Parser: String) {
  addIngestTokenV3(input: {
    repositoryName: $RepositoryName
    name: $Name
    parser: $Parser
  }) {
    name
    token
    parser {
      name
    }
  }
}
`

const assignParserMutation = `
mutation AssignParser($RepositoryName: String!, $TokenName: String!, $ParserName: String!) {
  assignParserToIngestToken(input: {
    repositoryName: $RepositoryName
    tokenName: $TokenName
    parserName: $ParserName
  }) {
    name
    token
    parser {
      name
    }
  }
}
`

const unassignParserMutation = `
mutation UnassignParser($RepositoryName: String!, $TokenName: String!) {
  unassignParserFromIngestToken(input: {
    repositoryName: $RepositoryName
    tokenName: $TokenName
  }) {
    name
    token
    parser {
      name
    }
  }
}
`

const removeIngestTokenMutation = `
mutation RemoveIngestToken($RepositoryName: String!, $Name: String!) {
  removeIngestToken(repositoryName: $RepositoryName, name: $Name) {
    __typename
  }
}
`

// ingestTokenResponse represents the response from list ingest tokens query
type ingestTokenResponse struct {
	Repository struct {
		IngestTokens []struct {
			Name   string `json:"name"`
			Token  string `json:"token"`
			Parser *struct {
				Name string `json:"name"`
			} `json:"parser"`
		} `json:"ingestTokens"`
	} `json:"repository"`
}

// addIngestTokenResponse represents the response from add ingest token mutation
type addIngestTokenResponse struct {
	AddIngestTokenV3 struct {
		Name   string `json:"name"`
		Token  string `json:"token"`
		Parser *struct {
			Name string `json:"name"`
		} `json:"parser"`
	} `json:"addIngestTokenV3"`
}

// assignParserResponse represents the response from assign/unassign parser mutation
type assignParserResponse struct {
	AssignParserToIngestToken *struct {
		Name   string `json:"name"`
		Token  string `json:"token"`
		Parser *struct {
			Name string `json:"name"`
		} `json:"parser"`
	} `json:"assignParserToIngestToken,omitempty"`
	UnassignParserFromIngestToken *struct {
		Name   string `json:"name"`
		Token  string `json:"token"`
		Parser *struct {
			Name string `json:"name"`
		} `json:"parser"`
	} `json:"unassignParserFromIngestToken,omitempty"`
}

// List returns all ingest tokens for the given repository
func (t *IngestTokens) List(repository string) ([]IngestToken, error) {
	var resp ingestTokenResponse
	err := t.client.Query(context.Background(), listIngestTokensQuery, map[string]interface{}{
		"RepositoryName": repository,
	}, &resp)
	if err != nil {
		return nil, err
	}

	tokens := make([]IngestToken, len(resp.Repository.IngestTokens))
	for i, token := range resp.Repository.IngestTokens {
		assignedParser := ""
		if token.Parser != nil {
			assignedParser = token.Parser.Name
		}
		tokens[i] = IngestToken{
			Name:           token.Name,
			Token:          token.Token,
			AssignedParser: assignedParser,
		}
	}

	return tokens, nil
}

// Get returns an ingest token by name
func (t *IngestTokens) Get(repository, name string) (*IngestToken, error) {
	tokens, err := t.List(repository)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if token.Name == name {
			return &token, nil
		}
	}

	return nil, fmt.Errorf("ingest token not found: %s", name)
}

// Add creates a new ingest token
func (t *IngestTokens) Add(repository, name, parser string) (*IngestToken, error) {
	variables := map[string]interface{}{
		"RepositoryName": repository,
		"Name":           name,
	}
	if parser != "" {
		variables["Parser"] = parser
	}

	var resp addIngestTokenResponse
	err := t.client.Query(context.Background(), addIngestTokenMutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	assignedParser := ""
	if resp.AddIngestTokenV3.Parser != nil {
		assignedParser = resp.AddIngestTokenV3.Parser.Name
	}

	return &IngestToken{
		Name:           resp.AddIngestTokenV3.Name,
		Token:          resp.AddIngestTokenV3.Token,
		AssignedParser: assignedParser,
	}, nil
}

// Update updates an existing ingest token's parser assignment
func (t *IngestTokens) Update(repository, name, parser string) (*IngestToken, error) {
	var resp assignParserResponse
	var err error

	if parser == "" {
		err = t.client.Query(context.Background(), unassignParserMutation, map[string]interface{}{
			"RepositoryName": repository,
			"TokenName":      name,
		}, &resp)
		if err != nil {
			return nil, err
		}
		if resp.UnassignParserFromIngestToken != nil {
			assignedParser := ""
			if resp.UnassignParserFromIngestToken.Parser != nil {
				assignedParser = resp.UnassignParserFromIngestToken.Parser.Name
			}
			return &IngestToken{
				Name:           resp.UnassignParserFromIngestToken.Name,
				Token:          resp.UnassignParserFromIngestToken.Token,
				AssignedParser: assignedParser,
			}, nil
		}
	} else {
		err = t.client.Query(context.Background(), assignParserMutation, map[string]interface{}{
			"RepositoryName": repository,
			"TokenName":      name,
			"ParserName":     parser,
		}, &resp)
		if err != nil {
			return nil, err
		}
		if resp.AssignParserToIngestToken != nil {
			assignedParser := ""
			if resp.AssignParserToIngestToken.Parser != nil {
				assignedParser = resp.AssignParserToIngestToken.Parser.Name
			}
			return &IngestToken{
				Name:           resp.AssignParserToIngestToken.Name,
				Token:          resp.AssignParserToIngestToken.Token,
				AssignedParser: assignedParser,
			}, nil
		}
	}

	// If we get here, fetch the token to return
	return t.Get(repository, name)
}

// Remove deletes an ingest token
func (t *IngestTokens) Remove(repository, name string) error {
	return t.client.Query(context.Background(), removeIngestTokenMutation, map[string]interface{}{
		"RepositoryName": repository,
		"Name":           name,
	}, nil)
}
