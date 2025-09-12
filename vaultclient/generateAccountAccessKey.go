package vaultclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/smithy-go"
)

const opGenerateAccountAccessKey = "GenerateAccountAccessKey"

type GenerateAccountAccessKeyInput struct {
	AccountName       *string
	ExternalAccessKey *string `locationName:"externalAccessKey"`
	ExternalSecretKey *string `locationName:"externalSecretKey"`
}

// String returns the string representation
func (s GenerateAccountAccessKeyInput) String() string {
	return fmt.Sprintf("GenerateAccountAccessKeyInput{AccountName: %v, ExternalAccessKey: %v, ExternalSecretKey: %v}",
		s.AccountName, s.ExternalAccessKey, s.ExternalSecretKey)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *GenerateAccountAccessKeyInput) Validate() error {
	invalidParams := &smithy.InvalidParamsError{Context: "GenerateAccountAccessKeyInput"}

	if s.AccountName == nil {
		invalidParams.Add(smithy.NewErrParamRequired("AccountName"))
	} else if len(*s.AccountName) < 1 {
		invalidParams.Add(NewErrParamMinLen("AccountName", 1))
	}

	if s.ExternalAccessKey != nil && len(*s.ExternalAccessKey) < 1 {
		invalidParams.Add(NewErrParamMinLen("ExternalAccessKey", 1))
	}

	if s.ExternalSecretKey != nil && len(*s.ExternalSecretKey) < 1 {
		invalidParams.Add(NewErrParamMinLen("ExternalSecretKey", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAccountName sets the AccountName field's value.
func (s *GenerateAccountAccessKeyInput) SetAccountName(v string) *GenerateAccountAccessKeyInput {
	s.AccountName = &v
	return s
}

// SetExternalAccessKey sets the ExternalAccessKey field's value.
func (s *GenerateAccountAccessKeyInput) SetExternalAccessKey(v string) *GenerateAccountAccessKeyInput {
	s.ExternalAccessKey = &v
	return s
}

// SetExternalSecretKey sets the ExternalSecretKey field's value.
func (s *GenerateAccountAccessKeyInput) SetExternalSecretKey(v string) *GenerateAccountAccessKeyInput {
	s.ExternalSecretKey = &v
	return s
}

func (s *GenerateAccountAccessKeyInput) getUrlValues() url.Values {
	formData := url.Values{}
	formData.Set("AccountName", *s.AccountName)
	if s.ExternalAccessKey != nil {
		formData.Set("externalAccessKey", *s.ExternalAccessKey)
	}
	if s.ExternalSecretKey != nil {
		formData.Set("externalSecretKey", *s.ExternalSecretKey)
	}
	return formData
}

// GenerateAccountAccessKey API operation generates a new access key for the account
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) GenerateAccountAccessKey(ctx context.Context, input *GenerateAccountAccessKeyInput) (*GenerateAccountAccessKeyOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.makeAWSRequest(ctx, opGenerateAccountAccessKey, input.getUrlValues())
	if err != nil {
		return nil, err
	}

	var output GenerateAccountAccessKeyOutput
	if err := c.handleAWSResponse(resp, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

type GeneratedKey struct {
	ID           *string    `locationName:"id" json:"id"`
	Value        *string    `locationName:"value" json:"value"`
	CreateDate   *time.Time `locationName:"createDate" json:"createDate"`
	LastUsedDate *time.Time `locationName:"lastUsedDate" json:"lastUsedDate"`
	Status       *string    `locationName:"status" json:"status"`
	UserID       *string    `locationName:"userId" json:"userId"`
}

// GenerateAccountAccessKeyOutput contains the response to a successful GenerateAccountAccessKey request.
type GenerateAccountAccessKeyOutput struct {
	GeneratedKey *GeneratedKey `type:"structure" locationName:"data" json:"data"`
}

// String returns the string representation
func (s GenerateAccountAccessKeyOutput) String() string {
	return fmt.Sprintf("GenerateAccountAccessKeyOutput{GeneratedKey: %v}", s.GeneratedKey)
}
