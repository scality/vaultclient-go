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

var (
	mockIsTruncated = true
	mockMarker      = "562385153604"
	mockMaxItems    = aws.Int64(10)
)

func mockListAccountsResponseBody(req *http.Request, t *testing.T) mockValue {
	return mockValue{
		"isTruncated": mockIsTruncated,
		"marker":      mockMarker,
		"accounts": []mockValue{
			{
				"arn":          mockArn,
				"id":           mockID,
				"name":         mockName,
				"createDate":   mockCreateDate,
				"emailAddress": mockEmail,
				"canonicalId":  mockCanonicalID,
				"quota":        mockQuotaMax,
			},
		},
	}
}

type listAccountsTest struct {
	maxItems    *int64
	marker      *string
	err         error
	description string
}

func listAccountsErrorMaker(errs []smithy.InvalidParamError) error {
	invalidParams := &smithy.InvalidParamsError{Context: "ListAccountsInput"}
	for _, err := range errs {
		invalidParams.Add(err)
	}
	return invalidParams
}

var listListAccountsTests = []listAccountsTest{
	{description: "Should pass with no optional parameter", err: nil},
	{description: "Should pass with valid maxItems", maxItems: mockMaxItems, err: nil},
	{description: "Should pass with valid marker", marker: &mockMarker, err: nil},

	{description: "Should fail with invalid maxItems", maxItems: aws.Int64(0), err: listAccountsErrorMaker([]smithy.InvalidParamError{NewErrParamMinValue("MaxItems", 1)})},
	{description: "Should fail with invalid marker", marker: aws.String(""), err: listAccountsErrorMaker([]smithy.InvalidParamError{NewErrParamMinLen("Marker", 1)})},
}

func TestListAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			t.Error(err)
		}

		if action := req.Form.Get("Action"); action != opListAccounts {
			t.Errorf("Expected Action=ListAccounts, got Action=%s", action)
		}

		// Send response to be tested
		resBody := mockListAccountsResponseBody(req, t)
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

	Convey("Test ListAccounts", t, func() {
		for _, tc := range listListAccountsTests {
			description := tc.description
			Convey(description, func() {
				ctx := context.Background()
				svc := New("foo", "bar", "", server.URL, "us-east-1")
				params := &ListAccountsInput{}
				if tc.marker != nil {
					params.SetMarker(*tc.marker)
				}
				if tc.maxItems != nil {
					params.SetMaxItems(*tc.maxItems)
				}
				res, err := svc.ListAccounts(ctx, params)
				if tc.err != nil {
					So(err.Error(), ShouldEqual, tc.err.Error())
				} else {
					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)

					So(*res.IsTruncated, ShouldEqual, mockIsTruncated)
					So(*res.Marker, ShouldEqual, mockMarker)
					So(len(res.Accounts), ShouldEqual, 1)
					account := res.Accounts[0]
					So(*account.Email, ShouldEqual, mockEmail)
					So(*account.Name, ShouldEqual, mockName)
					So(*account.ID, ShouldEqual, mockID)
					So(*account.Arn, ShouldEqual, mockArn)
					So(*account.CanonicalID, ShouldEqual, mockCanonicalID)
					So(*account.CreateDate, ShouldEqual, mockTime)
				}
			})
		}
	})
}

func TestAccountFromListUnmarshalJSON(t *testing.T) {
	Convey("Test AccountFromList UnmarshalJSON", t, func() {
		Convey("Should unmarshal quota from JSON number", func() {
			jsonData := `{"quota": 100}`
			var data AccountFromList
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldNotBeNil)
			So(*data.QuotaMax, ShouldEqual, 100)
		})

		Convey("Should unmarshal quota from JSON string", func() {
			jsonData := `{"quota": "200"}`
			var data AccountFromList
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldNotBeNil)
			So(*data.QuotaMax, ShouldEqual, 200)
		})

		Convey("Should unmarshal quota as nil when null", func() {
			jsonData := `{"quota": null}`
			var data AccountFromList
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldBeNil)
		})

		Convey("Should unmarshal quota as nil when field is missing", func() {
			jsonData := `{"name": "test"}`
			var data AccountFromList
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldBeNil)
			So(data.QuotaMax, ShouldBeNil)
		})

		Convey("Should fail when quota is an invalid string", func() {
			jsonData := `{"quota": "not-a-number"}`
			var data AccountFromList
			err := json.Unmarshal([]byte(jsonData), &data)
			So(err, ShouldNotBeNil)
		})
	})
}
