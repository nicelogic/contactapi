package contactapi

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

const (
	ContactRelationshipNone           string = "none"
	ContactRelationshipIAddedHim      string = "iAddedHim"
	ContactRelationshipHeAddedMe      string = "heAddedMe"
	ContactRelationshipAddedEachOther string = "addedEachOther"
)

type contactRelationshipResponse struct {
	ContactRelationship string
}

type Contact struct {
	ID           string
	RemarkName   string
	Relationship string
}

type ContactApiClient struct {
	client *graphql.Client
}

func (client *ContactApiClient) Init(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("contact api client endpoint is empty")
	}
	client.client = graphql.NewClient(endpoint)
	return nil
}

func (client *ContactApiClient) ContactRelationship(ctx context.Context, userID string, token string) (string, error) {
	if client.client == nil {
		return "", fmt.Errorf("ContactApiClient not init")
	}

	req := graphql.NewRequest(`
	query Query($userId: ID!) {
		contactRelationship(userId: $userId)
	  }
`)
	req.Var("userId", userID)
	req.Header.Set("authorization", "Bearer "+token)
	var response contactRelationshipResponse
	if err := client.client.Run(ctx, req, &response); err != nil {
		return "", err
	}
	return response.ContactRelationship, nil
}

func (client *ContactApiClient) Contact(ctx context.Context, contactId string, token string) (*Contact, error) {
	if client.client == nil {
		return nil, fmt.Errorf("ContactApiClient not init")
	}

	req := graphql.NewRequest(`
	query contact($contactId: ID!){
		contact(contactId: $contactId){
		  id
		  remarkName
		  relationship
		}
	  }
`)
	req.Var("contactId", contactId)
	req.Header.Set("authorization", "Bearer "+token)
	var response map[string]interface{}
	if err := client.client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	contact := &Contact{
		ID: contactId,
		RemarkName: "",
		Relationship: "",
	}
	contactMap, ok := response["contact"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response parse error")
	}
	contact.RemarkName, ok = contactMap["remarkName"].(string)
	if !ok {
		return nil, fmt.Errorf("response parse error")
	}
	contact.Relationship = contactMap["relationship"].(string)
	if !ok {
		return nil, fmt.Errorf("response parse error")
	}
	return contact, nil
}
