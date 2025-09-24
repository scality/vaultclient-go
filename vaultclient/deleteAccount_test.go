package vaultclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	. "github.com/smartystreets/goconvey/convey"
)

func mockDeleteAccountResponseBody(req *http.Request, t *testing.T) mockValue {
	return mockValue{}
}

type deleteAccountTest struct {
	accountName *string
	err         error
	description string
}

func deleteAccountErrorMaker(errs []smithy.InvalidParamError) error {
	invalidParams := &smithy.InvalidParamsError{Context: "DeleteAccountInput"}
	for _, err := range errs {
		invalidParams.Add(err)
	}
	return invalidParams
}

var listDeleteAccountTests = []deleteAccountTest{
	{description: "Should pass with valid accountName", accountName: &mockName, err: nil},

	{description: "Should fail if accountName is empty", accountName: aws.String(""), err: deleteAccountErrorMaker([]smithy.InvalidParamError{NewErrParamMinLen("AccountName", 1)})},
	{description: "Should fail if accountName is not set", err: deleteAccountErrorMaker([]smithy.InvalidParamError{smithy.NewErrParamRequired("AccountName")})},
}

func TestDeleteAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resBody := mockDeleteAccountResponseBody(req, t)
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

	Convey("Test DeleteAccount", t, func() {
		for _, tc := range listDeleteAccountTests {
			description := tc.description
			Convey(description, func() {
				ctx := context.Background()
				svc := New("foo", "bar", "", server.URL, "us-east-1")
				params := &DeleteAccountInput{}
				if tc.accountName != nil {
					params.SetAccountName(*tc.accountName)
				}
				_, err := svc.DeleteAccount(ctx, params)
				if tc.err != nil {
					So(err.Error(), ShouldEqual, tc.err.Error())
				} else {
					So(err, ShouldBeNil)
				}
			})
		}
	})
}
