package vaultclient

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/smithy-go"
)

const opGetAccount = "GetAccount"

type GetAccountInput struct {
	Arn         *string `locationName:"accountArn"`
	CanonicalID *string `locationName:"canonicalId"`
	Email       *string `locationName:"emailAddress"`
	ID          *string `locationName:"accountId"`
	Name        *string `locationName:"accountName"`
}

// String returns the string representation
func (s GetAccountInput) String() string {
	return fmt.Sprintf("GetAccountInput{Arn: %v, CanonicalID: %v, Email: %v, ID: %v, Name: %v}",
		s.Arn, s.CanonicalID, s.Email, s.ID, s.Name)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *GetAccountInput) Validate() error {
	invalidParams := &smithy.InvalidParamsError{Context: "GetAccountInput"}

	if s.Arn == nil && s.ID == nil && s.Name == nil && s.Email == nil && s.CanonicalID == nil {
		invalidParams.Add(smithy.NewErrParamRequired("Arn, ID, Name, Email or CanonicalId"))
	}

	if s.Arn != nil && len(*s.Arn) < 1 {
		invalidParams.Add(NewErrParamMinLen("Arn", 1))
	}

	if s.CanonicalID != nil && len(*s.CanonicalID) < 1 {
		invalidParams.Add(NewErrParamMinLen("CanonicalId", 1))
	}

	if s.Email != nil && len(*s.Email) < 1 {
		invalidParams.Add(NewErrParamMinLen("Email", 1))
	}

	if s.ID != nil && len(*s.ID) < 1 {
		invalidParams.Add(NewErrParamMinLen("ID", 1))
	}

	if s.Name != nil && len(*s.Name) < 1 {
		invalidParams.Add(NewErrParamMinLen("Name", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetArn sets the Arn field's value.
func (s *GetAccountInput) SetArn(v string) *GetAccountInput {
	s.Arn = &v
	return s
}

// SetCanonicalId sets the CanonicalId field's value.
func (s *GetAccountInput) SetCanonicalId(v string) *GetAccountInput {
	s.CanonicalID = &v
	return s
}

// SetEmail sets the Email field's value.
func (s *GetAccountInput) SetEmail(v string) *GetAccountInput {
	s.Email = &v
	return s
}

// SetId sets the Id field's value.
func (s *GetAccountInput) SetID(v string) *GetAccountInput {
	s.ID = &v
	return s
}

// SetName sets the Name field's value.
func (s *GetAccountInput) SetName(v string) *GetAccountInput {
	s.Name = &v
	return s
}

func (s *GetAccountInput) getUrlValues() url.Values {
	formData := url.Values{}
	if s.Arn != nil {
		formData.Set("accountArn", *s.Arn)
	}
	if s.CanonicalID != nil {
		formData.Set("canonicalId", *s.CanonicalID)
	}
	if s.Email != nil {
		formData.Set("emailAddress", *s.Email)
	}
	if s.ID != nil {
		formData.Set("accountId", *s.ID)
	}
	if s.Name != nil {
		formData.Set("accountName", *s.Name)
	}
	return formData
}

// GetAccount API operation gets details about a Vault account
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) GetAccount(ctx context.Context, input *GetAccountInput) (*GetAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.makeAWSRequest(ctx, opGetAccount, input.getUrlValues())
	if err != nil {
		return nil, err
	}

	var output GetAccountOutput
	if err := c.handleAWSResponse(resp, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

// GetAccountOutput contains the response to a successful GetAccount request.
type GetAccountOutput = AccountData

// String returns the string representation
func (s GetAccountOutput) String() string {
	return fmt.Sprintf("GetAccountOutput{Arn:%v, Name:%v, ID:%v}", s.Arn, s.Name, s.ID)
}
