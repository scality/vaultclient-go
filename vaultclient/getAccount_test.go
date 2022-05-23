package vaultclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	. "github.com/smartystreets/goconvey/convey"
)

type GetAccountTest struct {
	arn         *string
	canonicalId *string
	email       *string
	id          *string
	name        *string

	err         error
	description string
}

func GetAccountErrorMaker(errs []request.ErrInvalidParam) error {
	return invalidParamsErrorMaker(errs, "GetAccountInput")
}

var listGetAccountTests = []GetAccountTest{
	{description: "Should pass with valid arn", arn: &mockArn, err: nil},
	{description: "Should pass with valid canonicalId", canonicalId: &mockCanonicalID, err: nil},
	{description: "Should pass with valid email", email: &mockEmail, err: nil},
	{description: "Should pass with valid id", id: &mockID, err: nil},
	{description: "Should pass with valid name", name: &mockName, err: nil},

	{description: "Should fail if arn is empty", arn: aws.String(""), err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamMinLen("Arn", 1)})},
	{description: "Should fail if canonicalId is empty", canonicalId: aws.String(""), err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamMinLen("CanonicalId", 1)})},
	{description: "Should fail if email is empty", email: aws.String(""), err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamMinLen("Email", 1)})},
	{description: "Should fail if id is empty", id: aws.String(""), err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamMinLen("ID", 1)})},
	{description: "Should fail if name is empty", name: aws.String(""), err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamMinLen("Name", 1)})},

	{description: "Should fail if no field is set", err: GetAccountErrorMaker([]request.ErrInvalidParam{request.NewErrParamRequired("Arn, ID, Name, Email or CanonicalId")})},
}

func mockGetAccountResponseBody(req *http.Request, t *testing.T) mockValue {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	v, err := url.ParseQuery(string(buf.Bytes()))
	if err != nil {
		t.Error(err)
	}

	if !v.Has("accountArn") && !v.Has("canonicalId") && !v.Has("emailAddress") &&
		!v.Has("accountId") && !v.Has("accountName") {
		t.Error(errors.New("No parameter in query"))
	}

	if (v.Has("accountArn") && v.Get("accountArn") != mockArn) ||
		(v.Has("canonicalId") && v.Get("canonicalId") != mockCanonicalID) ||
		(v.Has("emailAddress") && v.Get("emailAddress") != mockEmail) ||
		(v.Has("accountId") && v.Get("accountId") != mockID) ||
		(v.Has("accountName") && v.Get("accountName") != mockName) {
		return mockValue{}
	}

	return mockValue{
		"id":           mockID,
		"emailAddress": mockEmail,
		"name":         mockName,
		"quotaMax":     nil,
		"arn":          mockArn,
		"canonicalId":  mockCanonicalID,
		"createDate":   mockCreateDate,
	}
}

func TestGetAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resBody := mockGetAccountResponseBody(req, t)
		if len(resBody) == 0 {
			res.WriteHeader(http.StatusNotFound)
		}

		rjson, err := json.Marshal(resBody)
		if err != nil {
			t.Error(err)
		}
		res.Write(rjson)
	}))
	defer server.Close()

	Convey("Test GetAccount", t, func() {
		for _, tc := range listGetAccountTests {
			description := tc.description
			Convey(description, func() {
				ctx := context.Background()
				sess := session.Must(session.NewSession(&aws.Config{
					Endpoint:    aws.String(server.URL),
					Region:      aws.String("us-east-1"),
					HTTPClient:  server.Client(),
					Credentials: credentials.NewStaticCredentials("foo", "bar", "000"),
				}))
				svc := New(sess)
				params := &GetAccountInput{}
				if tc.arn != nil {
					params.SetArn(*tc.arn)
				}
				if tc.name != nil {
					params.SetName(*tc.name)
				}
				if tc.id != nil {
					params.SetID(*tc.id)
				}
				if tc.email != nil {
					params.SetEmail(*tc.email)
				}
				if tc.canonicalId != nil {
					params.SetCanonicalId(*tc.canonicalId)
				}
				res, err := svc.GetAccount(ctx, params)
				if tc.err != nil {
					So(err.Error(), ShouldEqual, tc.err.Error())
				} else {
					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)

					So(*res.Email, ShouldEqual, mockEmail)
					So(*res.Name, ShouldEqual, mockName)
					So(*res.ID, ShouldEqual, mockID)
					So(*res.Arn, ShouldEqual, mockArn)
					So(*res.CanonicalID, ShouldEqual, mockCanonicalID)
					So(*res.CreateDate, ShouldEqual, mockTime)
				}
			})
		}
	})
}
