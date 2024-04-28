package api

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/danyouknowme/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Missing Basic Auth",
			setupAuth: func(request *http.Request) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Username or Password",
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "wrongPassword")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				AdminUsername: "adminTest",
				AdminPassword: "test!",
			}
			server := NewServer(&cfg, nil)

			server.router.GET("/auth", server.basicAuth(func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			}))

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestAcceptCSVExtension(t *testing.T) {
	testCases := []struct {
		name          string
		filePath      string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			filePath: filepath.Join("..", "testdata", "taxes.csv"),
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "Missing File",
			filePath: "",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "Invalid File Format",
			filePath: filepath.Join("..", "testdata", "taxes.txt"),
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(&config.Config{}, nil)

			server.router.GET("/csv", server.acceptCSVExtension(func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			}))

			if tc.filePath == "" {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest(http.MethodGet, "/csv", nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
				return
			}

			var file *os.File
			file, err := os.Open(tc.filePath)
			assert.NoError(t, err)
			defer file.Close()

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := createFormFile(writer, file.Name())
			require.NoError(t, err)

			_, err = io.Copy(part, file)
			require.NoError(t, err)
			require.NoError(t, writer.Close())

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/csv", body)
			require.NoError(t, err)

			request.Header.Set("Content-Type", writer.FormDataContentType())

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func createFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "taxFile", filename))
	h.Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	return w.CreatePart(h)
}
