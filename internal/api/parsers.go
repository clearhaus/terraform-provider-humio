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
      sourceCode
      testData
      tagFields
    }
  }
}
`

const createParserMutation = `
mutation CreateParser(
  $RepositoryName: String!
  $Name: String!
  $SourceCode: String!
  $TestData: [String!]
  $TagFields: [String!]
  $Force: Boolean
) {
  createParser(input: {
    repositoryName: $RepositoryName
    name: $Name
    sourceCode: $SourceCode
    testData: $TestData
    tagFields: $TagFields
    force: $Force
  }) {
    id
    name
  }
}
`

const deleteParserMutation = `
mutation DeleteParser($RepositoryName: String!, $ParserID: String!) {
  deleteParser(input: {
    repositoryName: $RepositoryName
    id: $ParserID
  })
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
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			SourceCode string   `json:"sourceCode"`
			TestData   []string `json:"testData"`
			TagFields  []string `json:"tagFields"`
		} `json:"parser"`
	} `json:"repository"`
}

// createParserResponse represents the response from create parser mutation
type createParserResponse struct {
	CreateParser struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createParser"`
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
	testCases := make([]ParserTestCase, len(rawParser.TestData))
	for i, testData := range rawParser.TestData {
		testCases[i] = ParserTestCase{
			Event: ParserTestEvent{
				RawString: testData,
			},
		}
	}

	return &Parser{
		ID:          rawParser.ID,
		Name:        rawParser.Name,
		Script:      rawParser.SourceCode,
		FieldsToTag: rawParser.TagFields,
		TestCases:   testCases,
	}, nil
}

// Add creates a new parser or updates an existing one
func (p *Parsers) Add(repository string, parser *Parser, force bool) (*Parser, error) {
	testData := make([]string, len(parser.TestCases))
	for i, tc := range parser.TestCases {
		testData[i] = tc.Event.RawString
	}

	variables := map[string]interface{}{
		"RepositoryName": repository,
		"Name":           parser.Name,
		"SourceCode":     parser.Script,
		"TagFields":      parser.FieldsToTag,
		"Force":          force,
	}

	if len(testData) > 0 {
		variables["TestData"] = testData
	}

	var resp createParserResponse
	err := p.client.Query(context.Background(), createParserMutation, variables, &resp)
	if err != nil {
		return nil, err
	}

	parser.ID = resp.CreateParser.ID
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
