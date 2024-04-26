package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danyouknowme/assessment-tax/config"

	"github.com/danyouknowme/assessment-tax/db"
	mockdb "github.com/danyouknowme/assessment-tax/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCalculateTaxAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         0.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "donation",
						"amount":        200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDeductions(gomock.Any()).
					Times(1).
					Return([]db.Deduction{
						{Type: "personal", Amount: 60000.0},
						{Type: "donation", Amount: 100000.0},
						{Type: "k-receipt", Amount: 50000.0},
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid Body(Missing TotalIncome)",
			body: map[string]interface{}{
				"wht":        0.0,
				"allowances": []map[string]interface{}{},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(TotalIncome is negative)",
			body: map[string]interface{}{
				"totalIncome": -500000.0,
				"allowances":  []map[string]interface{}{},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(WHT is negative)",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         -100000.0,
				"allowances":  []map[string]interface{}{},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(Wht more than TotalIncome)",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         1000000.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "donation",
						"amount":        200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(Invalid Allowance Type)",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         0.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "invalid",
						"amount":        200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Body(Invalid Allowance Amount)",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         0.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "donation",
						"amount":        -200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found Default Deductions",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         0.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "donation",
						"amount":        200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDeductions(gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Failed to Get Default Deductions",
			body: map[string]interface{}{
				"totalIncome": 500000.0,
				"wht":         0.0,
				"allowances": []map[string]interface{}{
					{
						"allowanceType": "donation",
						"amount":        200000.0,
					},
				},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDeductions(gomock.Any()).
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

			server := NewServer(&config.Config{}, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/tax/calculations"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
