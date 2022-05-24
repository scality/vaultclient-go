package vaultclient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
)

const opGetAccount = "GetAccount"

type GetAccountInput struct {
	Arn         *string `locationName:"accountArn"`
	CanonicalId *string `locationName:"canonicalId"`
	Email       *string `locationName:"emailAddress"`
	ID          *string `locationName:"accountId"`
	Name        *string `locationName:"accountName"`
}

// String returns the string representation
func (s GetAccountInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *GetAccountInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "GetAccountInput"}

	if s.Arn != nil && len(*s.Arn) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("Arn", 1))
	}

	if s.CanonicalId != nil && len(*s.CanonicalId) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("CanonicalId", 1))
	}

	if s.Email != nil && len(*s.Email) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("Email", 1))
	}

	if s.ID != nil && len(*s.ID) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("ID", 1))
	}

	if s.Name != nil && len(*s.Name) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("Name", 1))
	}

	if s.Arn == nil && s.ID == nil && s.Name == nil && s.Email == nil && s.CanonicalId == nil {
		invalidParams.Add(request.NewErrParamRequired("Arn, ID, Name, Email or CanonicalId"))
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
	s.CanonicalId = &v
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

// GetAccount API operation gets details about a Vault account
// and adds the ability to pass a context and additional request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *Vault) GetAccount(ctx aws.Context, input *GetAccountInput, opts ...request.Option) (*GetAccountOutput, error) {
	req, out := c.GetAccountRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

// GetAccountRequest generates a "aws/request.Request" representing the
// client's request for the GetAccount operation. The "output" return
// value will be populated with the request's response once the request completes
// successfully.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
func (c *Vault) GetAccountRequest(input *GetAccountInput) (req *request.Request, output *GetAccountOutput) {
	op := &request.Operation{
		Name:       opGetAccount,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &GetAccountInput{}
	}

	output = &GetAccountOutput{}
	req = c.newRequest(op, input, output)
	return
}

// GetAccountOutput contains the response to a successful GetAccount request.
type GetAccountOutput = AccountData

// String returns the string representation
func (s GetAccountOutput) String() string {
	return awsutil.Prettify(s)
}
