// +build go1.8

package envconfig

import (
	"errors"
	"net/url"
	"os"
	"testing"
)

type SpecWithURL struct {
	UrlValue   url.URL
	UrlPointer *url.URL
}

func TestParseURL(t *testing.T) {
	var s SpecWithURL

	os.Clearenv()
	os.Setenv("ENV_CONFIG_URLVALUE", "https://github.com/kelseyhightower/envconfig")
	os.Setenv("ENV_CONFIG_URLPOINTER", "https://github.com/kelseyhightower/envconfig")

	err := Process("env_config", &s)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	u, err := url.Parse("https://github.com/kelseyhightower/envconfig")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.UrlValue != *u {
		t.Errorf("expected %q, got %q", u, s.UrlValue.String())
	}

	if *s.UrlPointer != *u {
		t.Errorf("expected %q, got %q", u, s.UrlPointer)
	}
}

func TestParseURLError(t *testing.T) {
	var s SpecWithURL

	os.Clearenv()
	os.Setenv("ENV_CONFIG_URLPOINTER", "http_://foo")

	err := Process("env_config", &s)

	v, ok := err.(*ParseErrorList)
	if !ok {
		t.Fatalf("expected ParseErrorList, got %T %v", err, err)
	}

	perr := (*v)[0]
	if perr.FieldName != "UrlPointer" {
		t.Errorf("expected %s, got %v", "UrlPointer", perr.FieldName)
	}

	expectedUnerlyingError := url.Error{
		Op:  "parse",
		URL: "http_://foo",
		Err: errors.New("first path segment in URL cannot contain colon"),
	}

	if perr.Err.Error() != expectedUnerlyingError.Error() {
		t.Errorf("expected %q, got %q", expectedUnerlyingError, perr.Err)
	}
}
