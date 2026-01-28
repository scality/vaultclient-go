package vaultclient

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/smithy-go"
)

const opListAccounts = "ListAccounts"

type ListAccountsInput struct {
	Marker   *string
	MaxItems *int64
}

// String returns the string representation
func (s ListAccountsInput) String() string {
	return fmt.Sprintf("ListAccountsInput{Marker: %v, MaxItems: %v}", s.Marker, s.MaxItems)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ListAccountsInput) Validate() error {
	invalidParams := &smithy.InvalidParamsError{Context: "ListAccountsInput"}

	if s.Marker != nil && len(*s.Marker) < 1 {
		invalidParams.Add(NewErrParamMinLen("Marker", 1))
	}

	if s.MaxItems != nil {
		if *s.MaxItems < 1 {
			invalidParams.Add(NewErrParamMinValue("MaxItems", 1))
		}
		if *s.MaxItems > 1000 {
			invalidParams.Add(NewErrParamMaxValue("MaxItems", 1000))
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetMarker sets the Marker field's value.
func (s *ListAccountsInput) SetMarker(v string) *ListAccountsInput {
	s.Marker = &v
	return s
}

// SetMaxItems sets the MaxItems field's value.
func (s *ListAccountsInput) SetMaxItems(v int64) *ListAccountsInput {
	s.MaxItems = &v
	return s
}

func (s *ListAccountsInput) getUrlValues() url.Values {
	formData := url.Values{}
	if s.Marker != nil {
		formData.Set("Marker", *s.Marker)
	}
	if s.MaxItems != nil {
		formData.Set("MaxItems", strconv.FormatInt(*s.MaxItems, 10))
	}
	return formData
}

// ListAccounts API operation lists Vault accounts
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) ListAccounts(ctx context.Context, input *ListAccountsInput) (*ListAccountsOutput, error) {
	if input == nil {
		input = &ListAccountsInput{}
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.makeAWSRequest(ctx, opListAccounts, input.getUrlValues())
	if err != nil {
		return nil, err
	}

	var output ListAccountsOutput
	if err := c.handleAWSResponse(resp, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

type AccountFromList struct {
	Arn         *string    `locationName:"arn" json:"arn"`
	Name        *string    `locationName:"name" json:"name"`
	Email       *string    `locationName:"emailAddress" json:"emailAddress"`
	ID          *string    `locationName:"id" json:"id"`
	QuotaMax    *big.Int   `locationName:"quota" json:"quota"`
	CreateDate  *time.Time `locationName:"createDate" json:"createDate"`
	CanonicalID *string    `locationName:"canonicalId" json:"canonicalId"`
}

// ListAccountsOutput contains the response to a successful ListAccounts request.
type ListAccountsOutput struct {
	Accounts    []*AccountFromList `locationName:"accounts" json:"accounts"`
	IsTruncated *bool              `locationName:"isTruncated" json:"isTruncated"`
	Marker      *string            `locationName:"marker" json:"marker"`
}

// String returns the string representation
func (s ListAccountsOutput) String() string {
	return fmt.Sprintf("ListAccountsOutput{Accounts: %d accounts, IsTruncated: %v, Marker: %v}",
		len(s.Accounts), s.IsTruncated, s.Marker)
}
