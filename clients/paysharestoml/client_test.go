package paysharestoml

import (
	"net/http"
	"strings"
	"testing"

	"github.com/payshares/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientURL(t *testing.T) {
	//HACK:  we're testing an internal method rather than setting up a http client
	//mock.

	c := &Client{UseHTTP: false}
	assert.Equal(t, "https://payshares.org/.well-known/payshares.toml", c.url("payshares.org"))

	c = &Client{UseHTTP: true}
	assert.Equal(t, "http://payshares.org/.well-known/payshares.toml", c.url("payshares.org"))
}

func TestClient(t *testing.T) {
	h := httptest.NewClient()
	c := &Client{HTTP: h}

	// happy path
	h.
		On("GET", "https://payshares.org/.well-known/payshares.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation"`,
		)
	stoml, err := c.GetPaysharesToml("payshares.org")
	require.NoError(t, err)
	assert.Equal(t, "https://localhost/federation", stoml.FederationServer)

	// payshares.toml exceeds limit
	h.
		On("GET", "https://toobig.org/.well-known/payshares.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation`+strings.Repeat("0", PaysharesTomlMaxSize)+`"`,
		)
	stoml, err = c.GetPaysharesToml("toobig.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "payshares.toml response exceeds")
	}

	// not found
	h.
		On("GET", "https://missing.org/.well-known/payshares.toml").
		ReturnNotFound()
	stoml, err = c.GetPaysharesToml("missing.org")
	assert.EqualError(t, err, "http request failed with non-200 status code")

	// invalid toml
	h.
		On("GET", "https://json.org/.well-known/payshares.toml").
		ReturnJSON(http.StatusOK, map[string]string{"hello": "world"})
	stoml, err = c.GetPaysharesToml("json.org")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "toml decode failed")
	}
}
