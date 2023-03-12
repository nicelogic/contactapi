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

type ContactApiClient struct {
	client *graphql.Client
}

func (client *ContactApiClient)Init(endpoint string) (error){
	if endpoint == "" {
		return fmt.Errorf("contact api client endpoint is empty")
	}
	client.client = graphql.NewClient(endpoint)
	return nil
}

func (client *ContactApiClient)ContactRelationship(ctx context.Context, userID string, token string) (string, error) {
	if(client.client == nil){
		return "", fmt.Errorf("ContactApiClient not init")
	}

	req := graphql.NewRequest(`
	query Query($userId: ID!) {
		contactRelationship(userId: $userId)
	  }
`)
	req.Var("userId", userID)
	req.Header.Set("authorization", "Bearer " + token)
	var response contactRelationshipResponse 
	if err := client.client.Run(ctx, req, &response); err != nil {
		return "", err
	}
	return response.ContactRelationship, nil
}
