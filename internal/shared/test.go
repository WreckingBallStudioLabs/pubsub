package shared

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

const (
	// DatabaseName is the database name.
	DatabaseName = Test

	// TableName is the table/collection/index/etc name.
	TableName = Test

	// DefaultTimeout is the default timeout.
	DefaultTimeout = 30 * time.Second

	// DocumentID is the document ID.
	DocumentID = "VFzrpYMBXu5BQSZxo0qX"

	// DocumentIDInvalid is an invalid document ID.
	DocumentIDInvalid = "VFzrpYMBXu5BQSZxo0qY"

	// DocumentName is the document name.
	DocumentName = Test

	// DocumentNameUpdated is the updated document name.
	DocumentNameUpdated = "Test2"

	// DocumentVersion is the document version.
	DocumentVersion = "1.0.0"

	RoleName = "test"

	// Test name.
	Test = "test"
)

// TestDataS is the test data definition.
type TestDataS struct {
	Name    string `json:"name,omitempty" query:"name" db:"name" dbType:"varchar(255)"`
	Version string `json:"version,omitempty" query:"version" db:"version" dbType:"varchar(255)"`
}

// TestDataWithIDS is the test data definition.
type TestDataWithIDS struct {
	ID      string `json:"id,omitempty" query:"id" db:"id" dbType:"varchar(255)" bson:"_id"`
	Name    string `json:"name,omitempty" query:"name" db:"name" dbType:"varchar(255)"`
	Version string `json:"version,omitempty" query:"version" db:"version" dbType:"varchar(255)"`
}

var (
	// TestData is the test data.
	TestData = &TestDataS{
		Name:    DocumentName,
		Version: DocumentVersion,
	}

	// TestDataWithID is the test data.
	TestDataWithID = &TestDataWithIDS{
		ID:      DocumentID,
		Name:    DocumentName,
		Version: DocumentVersion,
	}

	// UpdatedTestData is the updated test data.
	UpdatedTestData = &TestDataS{
		Name:    DocumentNameUpdated,
		Version: DocumentVersion,
	}

	// UpdatedTestDataID is the updated test data with ID.
	UpdatedTestDataID = &TestDataWithIDS{
		ID:      DocumentID,
		Name:    DocumentNameUpdated,
		Version: DocumentVersion,
	}
)

// CreateHTTPTestServer creates a mocked HTTP server. Don't forget to defer
// close it!
//
// SEE: Test for usage.
func CreateHTTPTestServer(
	statusCode int,
	headers map[string]string,
	queryParams map[string]string,
	body string,
) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check headers in request.
		for key, value := range headers {
			if req.Header.Get(key) != value {
				http.Error(
					res,
					fmt.Sprintf("header not found (K: %s  V: %s)", key, value),
					http.StatusBadRequest,
				)

				return
			}
		}

		// Check query params in request.
		for key, value := range queryParams {
			if req.URL.Query().Get(key) != value {
				http.Error(
					res,
					fmt.Sprintf("query param not found (K: %s  V: %s)", key, value),
					http.StatusBadRequest,
				)

				return
			}
		}

		// Set status code.
		res.WriteHeader(statusCode)

		// Set body.
		if _, err := res.Write([]byte(body)); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)

			return
		}
	}))

	// Give enough time to be ready.
	time.Sleep(1 * time.Second)

	return testServer
}

// ErrorContains checks if an error contains a string.
func ErrorContains(err error, text ...string) bool {
	if err == nil {
		return false
	}

	for _, t := range text {
		if strings.Contains(err.Error(), t) {
			return true
		}
	}

	return false
}

// IsEnvironment determines if the current environment is one of the provided.
func IsEnvironment(environments ...string) bool {
	for _, environment := range environments {
		if os.Getenv("ENVIRONMENT") == environment {
			return true
		}
	}

	return false
}
