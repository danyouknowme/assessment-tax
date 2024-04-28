//go:build integration

package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestIntegrationCalculateTax(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	reqBody := `{"totalIncome": 500000.0,"wht": 0.0,"allowances": [{"allowanceType": "k-receipt","amount": 200000.0},{"allowanceType": "donation","amount": 100000.0}]}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/tax/calculations", serverPort), strings.NewReader(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	expected := `{"tax":14000,"taxRefund":0,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":14000},{"level":"500,001-1,000,000","tax":0},{"level":"1,000,001-2,000,000","tax":0},{"level":"2,000,001 ขึ้นไป","tax":0}]}`

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, expected, strings.TrimSpace(string(byteBody)))
}

func TestIntegrationCalculateTaxForCSV(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	file, err := os.Open(filepath.Join("..", "testdata", "taxes.csv"))
	require.NoError(t, err)
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := createFormFile(writer, file.Name())
	require.NoError(t, err)

	_, err = io.Copy(part, file)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/tax/calculations/upload-csv", serverPort), body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	expected := `{"taxes":[{"totalIncome":500000,"tax":29000},{"totalIncome":600000,"tax":0},{"totalIncome":750000,"tax":11250}]}`

	fmt.Println("Response:", string(byteBody))

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, expected, strings.TrimSpace(string(byteBody)))
}
