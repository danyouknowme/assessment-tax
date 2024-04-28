//go:build integration

package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAdminSettingPersonalDeduction(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	reqBody := `{"amount":100000.0}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/admin/deductions/personal", serverPort), strings.NewReader(reqBody))
	require.NoError(t, err)

	req.SetBasicAuth(testServer.config.AdminUsername, testServer.config.AdminPassword)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	expected := `{"personalDeduction":100000}`

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, expected, strings.TrimSpace(string(byteBody)))
}

func TestIntegrationAdminSettingKReceiptDeduction(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	reqBody := `{"amount": 20000.0}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/admin/deductions/k-receipt", serverPort), strings.NewReader(reqBody))
	require.NoError(t, err)

	req.SetBasicAuth(testServer.config.AdminUsername, testServer.config.AdminPassword)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	expected := `{"kReceipt":20000}`

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, expected, strings.TrimSpace(string(byteBody)))
}
