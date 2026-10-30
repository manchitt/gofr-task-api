package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gofr.dev/pkg/gofr/testutil"
)

func TestMain(m *testing.M) {
	os.Setenv("GOFR_TELEMETRY", "false")
	m.Run()
}

func TestIntegration_TaskAPI(t *testing.T) {
	configs := testutil.NewServerConfigs(t)

	go main()
	testutil.WaitForHTTPServer(t, configs.HTTPHost)

	client := http.Client{}

	// 1. Create Task
	createBody := []byte(`{"title":"Write GoFr Tutorial","description":"Complete article for open-source contribution","priority":"high"}`)
	req, err := http.NewRequest(http.MethodPost, configs.HTTPHost+"/tasks", bytes.NewReader(createBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResp struct {
		Data Task `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	require.NoError(t, err)
	assert.Equal(t, "Write GoFr Tutorial", createResp.Data.Title)
	assert.Equal(t, "pending", createResp.Data.Status)
	assert.Equal(t, "high", createResp.Data.Priority)
	taskID := createResp.Data.ID
	assert.Positive(t, taskID)

	// 2. List Tasks
	req, err = http.NewRequest(http.MethodGet, configs.HTTPHost+"/tasks", http.NoBody)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 3. Get Task by ID
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/%d", configs.HTTPHost, taskID), http.NoBody)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 4. Update Task (set status to in_progress)
	updateBody := []byte(`{"status":"in_progress"}`)
	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/%d", configs.HTTPHost, taskID), bytes.NewReader(updateBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 5. Delete Task
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", configs.HTTPHost, taskID), http.NoBody)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 6. Verify Task is deleted (404)
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/%d", configs.HTTPHost, taskID), http.NoBody)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// 7. Verify Health check endpoint
	req, err = http.NewRequest(http.MethodGet, configs.HTTPHost+"/.well-known/health", http.NoBody)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
