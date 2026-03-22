package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseJsonParser_ValidJSON(t *testing.T) {
	result, err := ResponseJsonParser([]byte(`{"key": "value", "num": 42}`))
	require.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestResponseJsonParser_InvalidJSON(t *testing.T) {
	result, err := ResponseJsonParser([]byte(`not valid json`))
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestResponseJsonParser_EmptyObject(t *testing.T) {
	result, err := ResponseJsonParser([]byte(`{}`))
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestResponseJsonParser_NestedJSON(t *testing.T) {
	result, err := ResponseJsonParser([]byte(`{"outer": {"inner": "val"}}`))
	require.NoError(t, err)
	outer, ok := result["outer"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "val", outer["inner"])
}

func TestCallPost_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`)) //nolint:errcheck
	}))
	defer ts.Close()

	body, status, err := CallPost(ts.URL, map[string]string{"key": "value"}, "user", "pass")
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Contains(t, string(body), "id")
}

func TestCallPost_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`)) //nolint:errcheck
	}))
	defer ts.Close()

	body, status, err := CallPost(ts.URL, nil, "", "")
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.NotNil(t, body)
}

func TestCallPost_BasicAuth(t *testing.T) {
	var receivedUser, receivedPass string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUser, receivedPass, _ = r.BasicAuth()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	CallPost(ts.URL, nil, "testuser", "testpass") //nolint:errcheck
	assert.Equal(t, "testuser", receivedUser)
	assert.Equal(t, "testpass", receivedPass)
}

func TestCallPost_ConnectionFailed(t *testing.T) {
	body, status, err := CallPost("http://127.0.0.1:1", nil, "", "")
	assert.Error(t, err)
	assert.Equal(t, 500, status)
	assert.Nil(t, body)
}
