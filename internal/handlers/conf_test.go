// Test Dummy endpoint configuration handlers
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chamilto/dummy/internal/errors"
	"github.com/chamilto/dummy/internal/models/dummy"
)

func mapToJSON(t *testing.T, m *map[string]interface{}) []byte {
	b, err := json.Marshal(m)

	if err != nil {
		t.Fatal(err)
	}

	return b
}

func dummyEndpointToMap(t *testing.T, de *dummy.DummyEndpoint) map[string]interface{} {
	var m map[string]interface{}

	b, err := json.Marshal(de)

	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(b, &m)

	if err != nil {
		t.Fatal(err)
	}

	return m
}

func getDefaultDummyEndpoint(t *testing.T) *dummy.DummyEndpoint {
	data := getTestData(t, "testdata/dummy/defaultDummyPost.json")
	var newEndpoint dummy.DummyEndpoint
	err := json.Unmarshal(data, &newEndpoint)

	if err != nil {
		t.Fatal(err)
	}

	return &newEndpoint
}

func postEndpoint(t *testing.T, de *dummy.DummyEndpoint, ctlr *HandlerController, expectedStatus int) func() {
	data, err := json.Marshal(de)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req := newRequest(t, "POST", "/dummy-config/endpoints", bytes.NewBuffer(data))
	ctlr.Router.ServeHTTP(rr, req)
	assertStatus(t, rr, expectedStatus)

	return func() {
		err = de.Delete(ctlr.DB)

		if err != nil {
			t.Fatalf("Failed to delete endpoint: %v", err)
		}
	}
}

func TestCreateDummyEndpointSaveError(t *testing.T) {
	REDIS_CALL_TRACKER = mockRedisCallData{}
	ctlr := newTestingHandlerController()
	type testCase struct {
		name   string
		failOn int
	}

	cases := []testCase{
		testCase{name: "fail on name save", failOn: 0},
		testCase{name: "fail on dummy data save", failOn: 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctlr.DB = &MockRedisClient{failOnCallCount: tc.failOn}
			de := getDefaultDummyEndpoint(t)
			postEndpoint(t, de, ctlr, 500)
		})
	}
}

func TestCreateDummyEndpointConflict(t *testing.T) {
	ctlr := newTestingHandlerController()
	de := getDefaultDummyEndpoint(t)
	defer postEndpoint(t, de, ctlr, 201)()

	type testCase struct {
		name     string
		endpoint func() *map[string]interface{}
	}

	cases := []testCase{
		testCase{
			name: "duplicate name",
			endpoint: func() *map[string]interface{} {
				m := dummyEndpointToMap(t, de)
				m["pathPattern"] = "/not/the/same"

				return &m
			},
		},
		testCase{
			name: "duplicate pathPattern",
			endpoint: func() *map[string]interface{} {
				m := dummyEndpointToMap(t, de)
				m["name"] = "a different name"

				return &m
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := tc.endpoint()
			rr := httptest.NewRecorder()
			req := newRequest(t, "POST", "/dummy-config/endpoints", bytes.NewBuffer(mapToJSON(t, body)))
			ctlr.Router.ServeHTTP(rr, req)

			assertStatus(t, rr, http.StatusConflict)

			var got errors.ErrorMessage
			json.Unmarshal(getResponseBytes(t, rr), &got)

			assertEquals(t, got.ErrorType, errors.ConflictError)
		})
	}
}

// TestCreateDummyEndpointErrantPayload uses each json payload in
// testdata/dummy/errant/ to test our dummy endpoint creation validation
func TestCreateDummyEndpointPayloadValidation(t *testing.T) {
	ctlr := newTestingHandlerController()

	for _, tp := range *getAllPayloadsForDir(t, "testdata/dummy/errant/") {
		t.Run(tp.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := newRequest(t, "POST", "/dummy-config/endpoints", bytes.NewBuffer(tp.data))
			ctlr.Router.ServeHTTP(rr, req)

			assertStatus(t, rr, http.StatusBadRequest)

			var got errors.ErrorMessage
			json.Unmarshal(getResponseBytes(t, rr), &got)

			assertEquals(t, got.ErrorType, errors.ValidationError)
		})
	}
}

// TestUpdateDummyEndpointErrantPayload uses each json payload in
// testdata/dummy/errant/ to test our dummy endpoint update validation
func TestUpdateDummyEndpointPayloadValidation(t *testing.T) {
	ctlr := newTestingHandlerController()

	for _, tp := range *getAllPayloadsForDir(t, "testdata/dummy/errant/") {
		t.Run(tp.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			// 'myendpoint': validation should raise 400 before we raise 404
			req := newRequest(t, "PUT", "/dummy-config/endpoints/myendpoint", bytes.NewBuffer(tp.data))
			ctlr.Router.ServeHTTP(rr, req)

			assertStatus(t, rr, http.StatusBadRequest)

			var got errors.ErrorMessage
			json.Unmarshal(getResponseBytes(t, rr), &got)

			assertEquals(t, got.ErrorType, errors.ValidationError)
		})
	}
}
