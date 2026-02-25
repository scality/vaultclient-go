package vaultclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	. "github.com/smartystreets/goconvey/convey"
)

type mockValue map[string]interface{}

var (
	mockName        = "myname"
	mockEmail       = "email@email.com"
	mockQuotaMax    = int64(1)
	mockID          = "893701217479"
	mockCanonicalID = "cdc9948f9124efae674ed122d52ce4d83d18c53ed05dcbf3765db56a051d7496"
	mockCreateDate  = "2020-04-20T01:54:54Z"
	mockArn         = "arn:arn:aws:iam::893701217479:/name/"
	mockTime, _     = time.Parse(time.RFC3339, mockCreateDate)
)

func mockResponseBody(v url.Values, t *testing.T) mockValue {
	var quotaMax int64
	if v.Get("quotaMax") != "" {
		var err error
		quotaMax, err = strconv.ParseInt(v.Get("quotaMax"), 10, 64)
		if err != nil {
			t.Error(err)
		}
	}
	return mockValue{
		"account": mockValue{
			"data": mockValue{
				"id":           mockID,
				"emailAddress": v.Get("emailAddress"),
				"name":         v.Get("name"),
				"quotaMax":     quotaMax,
				"arn":          mockArn,
				"canonicalId":  mockCanonicalID,
				"createDate":   mockCreateDate,
			},
		},
	}
}

type createAccountTest struct {
	name              *string
	email             *string
	quotaMax          *int64
	externalAccountID *string
	err               error
	description       string
}

func createAccountErrorMaker(errs []smithy.InvalidParamError) error {
	invalidParams := &smithy.InvalidParamsError{Context: "CreateAccountInput"}
	for _, err := range errs {
		invalidParams.Add(err)
	}
	return invalidParams
}

var listCreateAccountTests = []createAccountTest{
	{description: "Should pass with valid name and email", name: &mockName, email: &mockEmail, err: nil},
	{description: "Should pass with valid quotaMax", name: &mockName, email: &mockEmail, quotaMax: &mockQuotaMax, err: nil},
	{description: "Should pass with valid externalAccountID", name: &mockName, email: &mockEmail, externalAccountID: &mockID, err: nil},

	{description: "Should fail if name is empty", name: aws.String(""), email: &mockEmail, quotaMax: &mockQuotaMax, err: createAccountErrorMaker([]smithy.InvalidParamError{NewErrParamMinLen("Name", 1)})},
	{description: "Should fail if name is not set", email: &mockEmail, quotaMax: &mockQuotaMax, err: createAccountErrorMaker([]smithy.InvalidParamError{smithy.NewErrParamRequired("Name")})},
	{description: "Should fail if email is empty", name: &mockName, email: aws.String(""), quotaMax: &mockQuotaMax, err: createAccountErrorMaker([]smithy.InvalidParamError{NewErrParamMinLen("Email", 1)})},
	{description: "Should fail if email is not set", name: &mockName, quotaMax: &mockQuotaMax, err: createAccountErrorMaker([]smithy.InvalidParamError{smithy.NewErrParamRequired("Email")})},
	{description: "Should fail if quotaMax is set to 0", name: &mockName, email: &mockEmail, quotaMax: aws.Int64(0), err: createAccountErrorMaker([]smithy.InvalidParamError{NewErrParamMinValue("QuotaMax", 1)})},
	{description: "Should fail if name and email are not set", err: createAccountErrorMaker([]smithy.InvalidParamError{smithy.NewErrParamRequired("Name"), smithy.NewErrParamRequired("Email")})},
	{description: "Should fail if externalAccountID is empty", name: &mockName, email: &mockEmail, externalAccountID: aws.String(""), err: createAccountErrorMaker([]smithy.InvalidParamError{NewErrParamMinLen("ExternalAccountID", 1)})},
}

func TestCreateAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			t.Error(err)
			return
		}

		if req.Form.Get("Action") != opCreateAccount {
			t.Errorf("Expected Action=CreateAccount, got %s", req.Form.Get("Action"))
		}

		// Send response to be tested
		resBody := mockResponseBody(req.Form, t)
		rjson, err := json.Marshal(resBody)
		if err != nil {
			t.Error(err)
		}

		_, err = res.Write(rjson)
		if err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	Convey("Test CreateAccount", t, func() {
		for _, tc := range listCreateAccountTests {
			description := tc.description
			Convey(description, func() {
				ctx := context.Background()
				svc := New("foo", "bar", "", server.URL, "us-east-1")
				params := &CreateAccountInput{}
				if tc.name != nil {
					params.SetName(*tc.name)
				}
				if tc.email != nil {
					params.SetEmail(*tc.email)
				}
				if tc.quotaMax != nil {
					params.SetQuotaMax(*tc.quotaMax)
				}
				if tc.externalAccountID != nil {
					params.SetExternalAccountID(*tc.externalAccountID)
				}
				res, err := svc.CreateAccount(ctx, params)
				if tc.err != nil {
					So(err.Error(), ShouldEqual, tc.err.Error())
				} else {
					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)

					So(*res.GetAccount().Email, ShouldEqual, *tc.email)
					So(*res.GetAccount().Name, ShouldEqual, *tc.name)
					So(*res.GetAccount().ID, ShouldEqual, mockID)
					So(*res.GetAccount().Arn, ShouldEqual, mockArn)
					So(*res.GetAccount().CanonicalID, ShouldEqual, mockCanonicalID)
					So(*res.GetAccount().CreateDate, ShouldEqual, mockTime)
					// optional property
					if tc.quotaMax == nil {
						So(*res.GetAccount().QuotaMax, ShouldEqual, 0)
					} else {
						So(*res.GetAccount().QuotaMax, ShouldEqual, *tc.quotaMax)
					}
				}
			})
		}
	})
}

func TestAccountDataUnmarshalJSON(t *testing.T) {
	Convey("Test AccountData UnmarshalJSON", t, func() {
		Convey("Should unmarshal quotaMax from JSON number", func() {
			jsonData := `{"quotaMax": 100}`
			var data AccountData
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldNotBeNil)
			So(*data.QuotaMax, ShouldEqual, 100)
		})

		Convey("Should unmarshal quotaMax from JSON string", func() {
			jsonData := `{"quotaMax": "200"}`
			var data AccountData
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldNotBeNil)
			So(*data.QuotaMax, ShouldEqual, 200)
		})

		Convey("Should unmarshal quotaMax as nil when null", func() {
			jsonData := `{"quotaMax": null}`
			var data AccountData
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldBeNil)
		})

		Convey("Should unmarshal quotaMax as nil when field is missing", func() {
			jsonData := `{"name": "test"}`
			var data AccountData
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldBeNil)
		})

		Convey("Should fail when quotaMax is an invalid string", func() {
			jsonData := `{"quotaMax": "not-a-number"}`
			var data AccountData
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldNotBeNil)
		})
	})
}
