package vaultclient

import (
	"net/http"
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
