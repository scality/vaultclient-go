package vaultclient

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/smithy-go"
)

const opDeleteAccount = "DeleteAccount"

type DeleteAccountInput struct {
	AccountName *string
}

// String returns the string representation
func (s DeleteAccountInput) String() string {
	return fmt.Sprintf("DeleteAccountInput{AccountName: %v}", s.AccountName)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteAccountInput) Validate() error {
	invalidParams := &smithy.InvalidParamsError{Context: "DeleteAccountInput"}

	if s.AccountName == nil {
		invalidParams.Add(smithy.NewErrParamRequired("AccountName"))
	} else if len(*s.AccountName) < 1 {
		invalidParams.Add(NewErrParamMinLen("AccountName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetAccountName sets the AccountName field's value.
func (s *DeleteAccountInput) SetAccountName(v string) *DeleteAccountInput {
	s.AccountName = &v
	return s
}

func (s *DeleteAccountInput) getUrlValues() url.Values {
	formData := url.Values{}
	formData.Set("AccountName", *s.AccountName)

	return formData
}

// DeleteAccount API operation deletes Vault account
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) DeleteAccount(ctx context.Context, input *DeleteAccountInput) (*DeleteAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.makeAWSRequest(ctx, opDeleteAccount, input.getUrlValues())
	if err != nil {
		return nil, err
	}

	if err := c.handleAWSResponse(resp, nil); err != nil {
		return nil, err
	}

	return &DeleteAccountOutput{}, nil
}

// DeleteAccountOutput contains the response to a successful DeleteAccount request.
type DeleteAccountOutput struct{}
