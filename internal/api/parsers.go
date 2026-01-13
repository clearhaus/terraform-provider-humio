package api

import (
	"context"
	"fmt"
)

// ParserTestEvent represents a test event for a parser
type ParserTestEvent struct {
	RawString string
}

// ParserTestCase represents a test case for a parser
type ParserTestCase struct {
	Event ParserTestEvent
}

// Parser represents a Humio parser
type Parser struct {
	ID          string
	Name        string
	Script      string
	FieldsToTag []string
	TestCases   []ParserTestCase
}

// Parsers provides operations for managing parsers
type Parsers struct {
	client *Client
}

const listParsersQuery = `
query ListParsers($RepositoryName: String!) {
  repository(name: $RepositoryName) {
    parsers {
      id
      name
      isBuiltIn
    }
  }
}
`

const getParserQuery = `
query GetParser($RepositoryName: String!, $ParserName: String!) {
  repository(name: $RepositoryName) {
    parser(name: $ParserName) {
      id
      name
      script
      testCases {
        event {
          rawString
        }
      }
      fieldsToTag
    }
  }
}
`

const createParserMutation = `
mutation CreateParser(
  $RepositoryName: RepoOrViewName!
  $Name: String!
  $Script: String!
  $TestCases: [ParserTestCaseInput!]!
  $FieldsToTag: [String!]!
  $FieldsToBeRemovedBeforeParsing: [String!]!
) {
  createParserV2(input: {
    repositoryName: $RepositoryName
    name: $Name
    script: $Script
    testCases: $TestCases
    fieldsToTag: $FieldsToTag
    fieldsToBeRemovedBeforeParsing: $FieldsToBeRemovedBeforeParsing
  }) {
    id
    name
  }
}
`

const updateParserMutation = `
mutation UpdateParser(
  $RepositoryName: RepoOrViewName!
  $ID: String!
  $Name: String!
  $Script: UpdateParserScriptInput!
  $TestCases: [ParserTestCaseInput!]!
  $FieldsToTag: [String!]!
  $FieldsToBeRemovedBeforeParsing: [String!]!
) {
  updateParserV2(input: {
    repositoryName: $RepositoryName
    id: $ID
    name: $Name
    script: $Script
    testCases: $TestCases
    fieldsToTag: $FieldsToTag
    fieldsToBeRemovedBeforeParsing: $FieldsToBeRemovedBeforeParsing
  }) {
    id
    name
  }
}
`

const deleteParserMutation = `
mutation DeleteParser($RepositoryName: RepoOrViewName!, $ParserID: String!) {
  deleteParser(input: {
    repositoryName: $RepositoryName
    id: $ParserID
  }) {
    __typename
  }
}
`

// listParsersResponse represents the response from list parsers query
type listParsersResponse struct {
	Repository struct {
		Parsers []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			IsBuiltIn bool   `json:"isBuiltIn"`
		} `json:"parsers"`
	} `json:"repository"`
}

// getParserResponse represents the response from get parser query
type getParserResponse struct {
	Repository struct {
		Parser *struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Script    string `json:"script"`
			TestCases []struct {
				Event struct {
					RawString string `json:"rawString"`
				} `json:"event"`
			} `json:"testCases"`
			FieldsToTag []string `json:"fieldsToTag"`
		} `json:"parser"`
	} `json:"repository"`
}

// createParserResponse represents the response from create parser mutation
type createParserResponse struct {
	CreateParserV2 struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createParserV2"`
}

// updateParserResponse represents the response from update parser mutation
type updateParserResponse struct {
	UpdateParserV2 struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"updateParserV2"`
}

// List returns all parsers for the given repository
func (p *Parsers) List(repository string) ([]Parser, error) {
	var resp listParsersResponse
	err := p.client.Query(context.Background(), listParsersQuery, map[string]interface{}{
		"RepositoryName": repository,
	}, &resp)
	if err != nil {
		return nil, err
	}

	parsers := make([]Parser, 0, len(resp.Repository.Parsers))
	for _, parser := range resp.Repository.Parsers {
		if !parser.IsBuiltIn {
			parsers = append(parsers, Parser{
				ID:   parser.ID,
				Name: parser.Name,
			})
		}
	}

	return parsers, nil
}

// Get returns a parser by name
func (p *Parsers) Get(repository, name string) (*Parser, error) {
	var resp getParserResponse
	err := p.client.Query(context.Background(), getParserQuery, map[string]interface{}{
		"RepositoryName": repository,
		"ParserName":     name,
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Repository.Parser == nil {
		return nil, fmt.Errorf("parser not found: %s", name)
	}

	rawParser := resp.Repository.Parser
	testCases := make([]ParserTestCase, len(rawParser.TestCases))
	for i, tc := range rawParser.TestCases {
		testCases[i] = ParserTestCase{
			Event: ParserTestEvent{
				RawString: tc.Event.RawString,
			},
		}
	}

	return &Parser{
		ID:          rawParser.ID,
		Name:        rawParser.Name,
		Script:      rawParser.Script,
		FieldsToTag: rawParser.FieldsToTag,
		TestCases:   testCases,
	}, nil
}

// Add creates a new parser or updates an existing one
func (p *Parsers) Add(repository string, parser *Parser, force bool) (*Parser, error) {
	testCases := make([]map[string]interface{}, len(parser.TestCases))
	for i, tc := range parser.TestCases {
		testCases[i] = map[string]interface{}{
			"event": map[string]interface{}{
				"rawString": tc.Event.RawString,
			},
		}
	}

	fieldsToTag := parser.FieldsToTag
	if fieldsToTag == nil {
		fieldsToTag = []string{}
	}

	variables := map[string]interface{}{
		"RepositoryName":                 repository,
		"Name":                           parser.Name,
		"Script":                         parser.Script,
		"TestCases":                      testCases,
		"FieldsToTag":                    fieldsToTag,
		"FieldsToBeRemovedBeforeParsing": []string{},
	}

	var resp createParserResponse
	err := p.client.Query(context.Background(), createParserMutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	parser.ID = resp.CreateParserV2.ID
	return parser, nil
}

// Update updates an existing parser
func (p *Parsers) Update(repository string, parser *Parser) (*Parser, error) {
	testCases := make([]map[string]interface{}, len(parser.TestCases))
	for i, tc := range parser.TestCases {
		testCases[i] = map[string]interface{}{
			"event": map[string]interface{}{
				"rawString": tc.Event.RawString,
			},
		}
	}

	fieldsToTag := parser.FieldsToTag
	if fieldsToTag == nil {
		fieldsToTag = []string{}
	}

	variables := map[string]interface{}{
		"RepositoryName": repository,
		"ID":             parser.ID,
		"Name":           parser.Name,
		"Script": map[string]interface{}{
			"script": parser.Script,
		},
		"TestCases":                      testCases,
		"FieldsToTag":                    fieldsToTag,
		"FieldsToBeRemovedBeforeParsing": []string{},
	}

	var resp updateParserResponse
	err := p.client.Query(context.Background(), updateParserMutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	parser.ID = resp.UpdateParserV2.ID
	return parser, nil
}

// Delete deletes a parser by name
func (p *Parsers) Delete(repository, name string) error {
	// First get the parser to find its ID
	parser, err := p.Get(repository, name)
	if err != nil {
		return err
	}

	return p.client.Query(context.Background(), deleteParserMutation, map[string]interface{}{
		"RepositoryName": repository,
		"ParserID":       parser.ID,
	}, nil)
}
