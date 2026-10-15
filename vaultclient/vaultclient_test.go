package vaultclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewHTTPClient(t *testing.T) {
	Convey("NewHTTPClient returns a client with the default 60s timeout", t, func() {
		c := NewHTTPClient()
		So(c, ShouldNotBeNil)
		So(c.Timeout, ShouldEqual, 60*time.Second)
	})
}

func TestNewOptions(t *testing.T) {
	Convey("New", t, func() {
		Convey("uses a default 60s HTTP client when no option is given", func() {
			v := New("ak", "sk", "", "http://vault", "us-east-1")
			So(v.httpClient, ShouldNotBeNil)
			So(v.httpClient.Timeout, ShouldEqual, 60*time.Second)
		})

		Convey("WithHTTPClient overrides the default client", func() {
			custom := &http.Client{Timeout: 5 * time.Second}
			v := New("ak", "sk", "", "http://vault", "us-east-1", WithHTTPClient(custom))
			So(v.httpClient, ShouldEqual, custom)
		})

		Convey("WithHTTPClient(nil) keeps the default client", func() {
			v := New("ak", "sk", "", "http://vault", "us-east-1", WithHTTPClient(nil))
			So(v.httpClient, ShouldNotBeNil)
			So(v.httpClient.Timeout, ShouldEqual, 60*time.Second)
		})
	})
}

type reqUIDKey struct{}

var signedHeadersRe = regexp.MustCompile(`SignedHeaders=([^,]+)`)

func signedHeaders(req *http.Request) []string {
	m := signedHeadersRe.FindStringSubmatch(req.Header.Get("Authorization"))
	if m == nil {
		return nil
	}
	return strings.Split(m[1], ";")
}

func TestPreSignHook(t *testing.T) {
	Convey("makeAWSRequest with pre-sign hooks", t, func() {
		var calls atomic.Int32
		var received *http.Request
		status := http.StatusOK
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			received = r
			w.WriteHeader(status)
		}))
		defer server.Close()

		ctx := context.WithValue(context.Background(), reqUIDKey{}, "req-123")
		setReqUID := func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Scal-Request-Uids", ctx.Value(reqUIDKey{}).(string))
			return nil
		}

		Convey("sends no extra header without a hook", func() {
			v := New("ak", "sk", "", server.URL, "us-east-1")
			resp, err := v.makeAWSRequest(ctx, "ListAccounts", url.Values{})
			So(err, ShouldBeNil)
			So(resp.Body.Close(), ShouldBeNil)
			So(received.Header.Get("X-Scal-Request-Uids"), ShouldBeEmpty)
			So(signedHeaders(received), ShouldNotContain, "x-scal-request-uids")
		})

		Convey("signs headers set by the hook", func() {
			v := New("ak", "sk", "", server.URL, "us-east-1", WithPreSignHook(setReqUID))
			resp, err := v.makeAWSRequest(ctx, "ListAccounts", url.Values{})
			So(err, ShouldBeNil)
			So(resp.Body.Close(), ShouldBeNil)
			So(received.Header.Get("X-Scal-Request-Uids"), ShouldEqual, "req-123")
			So(signedHeaders(received), ShouldContain, "x-scal-request-uids")
		})

		Convey("runs hooks in registration order", func() {
			var order []string
			hook := func(name string) PreSignHook {
				return func(_ context.Context, req *http.Request) error {
					order = append(order, name)
					req.Header.Set("X-Scal-Hook", name)
					return nil
				}
			}
			v := New("ak", "sk", "", server.URL, "us-east-1",
				WithPreSignHook(hook("first")), WithPreSignHook(nil), WithPreSignHook(hook("second")))
			resp, err := v.makeAWSRequest(ctx, "ListAccounts", url.Values{})
			So(err, ShouldBeNil)
			So(resp.Body.Close(), ShouldBeNil)
			So(order, ShouldResemble, []string{"first", "second"})
			So(received.Header.Get("X-Scal-Hook"), ShouldEqual, "second")
		})

		Convey("aborts the request when a hook fails", func() {
			hookErr := errors.New("boom")
			v := New("ak", "sk", "", server.URL, "us-east-1", WithPreSignHook(
				func(context.Context, *http.Request) error { return hookErr }))
			resp, err := v.makeAWSRequest(ctx, "ListAccounts", url.Values{})
			So(resp, ShouldBeNil)
			So(errors.Is(err, hookErr), ShouldBeTrue)
			So(calls.Load(), ShouldEqual, 0)
		})

		Convey("runs hooks again on retries", func() {
			status = http.StatusServiceUnavailable
			var hookCalls int
			v := New("ak", "sk", "", server.URL, "us-east-1", WithPreSignHook(
				func(ctx context.Context, req *http.Request) error {
					hookCalls++
					if hookCalls > 1 {
						status = http.StatusOK
					}
					return setReqUID(ctx, req)
				}))
			resp, err := v.makeAWSRequest(ctx, "ListAccounts", url.Values{})
			So(err, ShouldBeNil)
			So(resp.Body.Close(), ShouldBeNil)
			So(hookCalls, ShouldEqual, 2)
			So(calls.Load(), ShouldEqual, 2)
			So(signedHeaders(received), ShouldContain, "x-scal-request-uids")
		})
	})
}
