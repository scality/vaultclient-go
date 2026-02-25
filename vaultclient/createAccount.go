package vaultclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/smithy-go"
)

const opCreateAccount = "CreateAccount"

type CreateAccountInput struct {
	Name              *string `locationName:"name"`
	Email             *string `locationName:"emailAddress"`
	QuotaMax          *int64  `locationName:"quotaMax"`
	ExternalAccountID *string `locationName:"externalAccountId"`
}

// String returns the string representation
func (s CreateAccountInput) String() string {
	return fmt.Sprintf("CreateAccountInput{Name: %v, Email: %v, QuotaMax: %v, ExternalAccountID: %v}",
		s.Name, s.Email, s.QuotaMax, s.ExternalAccountID)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateAccountInput) Validate() error {
	invalidParams := &smithy.InvalidParamsError{Context: "CreateAccountInput"}
	if s.Name == nil {
		invalidParams.Add(smithy.NewErrParamRequired("Name"))
	} else if len(*s.Name) < 1 {
		invalidParams.Add(NewErrParamMinLen("Name", 1))
	}

	if s.Email == nil {
		invalidParams.Add(smithy.NewErrParamRequired("Email"))
	} else if len(*s.Email) < 1 {
		invalidParams.Add(NewErrParamMinLen("Email", 1))
	}

	if s.QuotaMax != nil && *s.QuotaMax < 1 {
		invalidParams.Add(NewErrParamMinValue("QuotaMax", 1))
	}

	if s.ExternalAccountID != nil && len(*s.ExternalAccountID) < 1 {
		invalidParams.Add(NewErrParamMinLen("ExternalAccountID", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetName sets the Name field's value.
func (s *CreateAccountInput) SetName(v string) *CreateAccountInput {
	s.Name = &v
	return s
}

// SetEmail sets the Email field's value.
func (s *CreateAccountInput) SetEmail(v string) *CreateAccountInput {
	s.Email = &v
	return s
}

// SetQuotaMax sets the QuotaMax field's value.
func (s *CreateAccountInput) SetQuotaMax(v int64) *CreateAccountInput {
	s.QuotaMax = &v
	return s
}

// SetExternalAccountID sets the ExternalAccountID field's value.
func (s *CreateAccountInput) SetExternalAccountID(v string) *CreateAccountInput {
	s.ExternalAccountID = &v
	return s
}

func (s *CreateAccountInput) getUrlValues() url.Values {
	formData := url.Values{}
	formData.Set("name", *s.Name)
	formData.Set("emailAddress", *s.Email)
	if s.QuotaMax != nil {
		formData.Set("quotaMax", strconv.FormatInt(*s.QuotaMax, 10))
	}
	if s.ExternalAccountID != nil {
		formData.Set("externalAccountId", *s.ExternalAccountID)
	}
	return formData
}

// CreateAccount API operation creates a new Vault account
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) CreateAccount(ctx context.Context, input *CreateAccountInput) (*CreateAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.makeAWSRequest(ctx, opCreateAccount, input.getUrlValues())
	if err != nil {
		return nil, err
	}

	var output CreateAccountOutput
	if err := c.handleAWSResponse(resp, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

// Account contains information about a Vault account.
type Account struct {
	AccountData *AccountData `type:"structure" locationName:"data" json:"data"`
}

type AccountData struct {
	Arn         *string    `locationName:"arn" json:"arn"`
	Name        *string    `locationName:"name" json:"name"`
	Email       *string    `locationName:"emailAddress" json:"emailAddress"`
	ID          *string    `locationName:"id" json:"id"`
	QuotaMax    *int64     `locationName:"quotaMax" json:"quotaMax"`
	CreateDate  *time.Time `locationName:"createDate" json:"createDate"`
	CanonicalID *string    `locationName:"canonicalId" json:"canonicalId"`
	AliasList   []*string  `locationName:"aliasList" json:"aliasList"`
}

// UnmarshalJSON implements custom unmarshaling for AccountData to handle
// quotaMax being either a JSON number or a JSON string.
func (a *AccountData) UnmarshalJSON(data []byte) error {
	// Use an alias to avoid infinite recursion
	type Alias AccountData
	aux := &struct {
		QuotaMax *json.Number `json:"quotaMax"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.QuotaMax != nil {
		quotaMax, err := aux.QuotaMax.Int64()
		if err != nil {
			return err
		}
		a.QuotaMax = &quotaMax
	}

	return nil
}

// CreateAccountOutput contains the response to a successful CreateAccount request.
type CreateAccountOutput struct {
	Account *Account `type:"structure" locationName:"account" json:"account"`
}

// String returns the string representation
func (s CreateAccountOutput) String() string {
	return fmt.Sprintf("CreateAccountOutput{Account: %v}", s.Account)
}

// GetAccount returns AccountData
func (s CreateAccountOutput) GetAccount() *AccountData {
	return s.Account.AccountData
}
