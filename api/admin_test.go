package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/danyouknowme/assessment-tax/config"
	"github.com/danyouknowme/assessment-tax/db"
	mockdb "github.com/danyouknowme/assessment-tax/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminSetPersonalDeductionAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          map[string]float64
		setupAuth     func(request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]float64{
				"amount": 70000.0,
			},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDeductionByType(gomock.Any(), gomock.Any(), db.UpdateDeductionParams{Amount: 70000.0}).
					Times(1).
					Return(&db.Deduction{Type: "personal", Amount: 70000.0}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				expected := `{"personalDeduction":70000}`
				require.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(recorder.Body.String()))
			},
		},
		{
			name: "Invalid Body(Missing Amount)",
			body: map[string]float64{},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(Negative Amount)",
			body: map[string]float64{
				"amount": -70000.0,
			},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found Personal Deduction",
			body: map[string]float64{
				"amount": 70000.0,
			},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDeductionByType(gomock.Any(), gomock.Any(), db.UpdateDeductionParams{Amount: 70000.0}).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Failed to Update Personal Deduction",
			body: map[string]float64{
				"amount": 70000.0,
			},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDeductionByType(gomock.Any(), gomock.Any(), db.UpdateDeductionParams{Amount: 70000.0}).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			cfg := config.Config{
				AdminUsername: "adminTest",
				AdminPassword: "test!",
			}

			server := NewServer(&cfg, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/admin/deductions/personal"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestAdminSetKReceiptDeductionAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          map[string]float64
		setupAuth     func(request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]float64{
				"amount": 70000.0,
			},
			setupAuth: func(request *http.Request) {
				request.SetBasicAuth("adminTest", "test!")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDeductionByType(gomock.Any(), gomock.Any(), db.UpdateDeductionParams{Amount: 70000.0}).
					Times(1).
					Return(&db.Deduction{Type: "k-receipt", Amount: 70000.0}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				expected := `{"kReceipt":70000}`
				require.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(recorder.Body.String()))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			cfg := config.Config{
				AdminUsername: "adminTest",
				AdminPassword: "test!",
			}

			server := NewServer(&cfg, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/admin/deductions/k-receipt"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
