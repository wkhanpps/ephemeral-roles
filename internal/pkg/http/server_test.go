package http

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ewohltman/ephemeral-roles/internal/pkg/mock"
)

const (
	testPort = "8081"
	testURL  = "http://localhost:" + testPort

	serverStartupDelay = 50 * time.Millisecond

	expectedGuildsFile = "testdata/guilds.json"
)

func TestNewServer(t *testing.T) {
	log := mock.NewLogger()

	session, err := mock.NewSession()
	if err != nil {
		t.Fatalf("Error obtaining mock session: %s", err)
	}

	defer mock.SessionClose(t, session)

	testServer := NewServer(log, session, testPort)

	go func() {
		serverErr := testServer.ListenAndServe()
		if !errors.Is(serverErr, http.ErrServerClosed) {
			t.Errorf("Test server error: %s", serverErr)
		}
	}()

	time.Sleep(serverStartupDelay)

	client := NewClient(nil, nil, "")

	testRootEndpoint(t, client)
	testGuildsEndpoint(t, client)

	ctx, cancelContext := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelContext()

	err = testServer.Shutdown(ctx)
	if err != nil {
		t.Errorf("Error closing test server: %s", err)
	}
}

func testRootEndpoint(t *testing.T, client *http.Client) {
	resp, err := doRequest(context.TODO(), client, testURL+rootEndpoint)
	if err != nil {
		t.Fatal(err)
	}

	err = drainCloseResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
}

func testGuildsEndpoint(t *testing.T, client *http.Client) {
	expectedGuildsBytes, err := ioutil.ReadFile(expectedGuildsFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedGuilds := make(sortableGuilds, 0)

	err = json.Unmarshal(expectedGuildsBytes, &expectedGuilds)
	if err != nil {
		t.Fatalf("Error unmarshaling expected guild data: %s", err)
	}

	resp, err := doRequest(context.TODO(), client, testURL+guildsEndpoint)
	if err != nil {
		t.Fatal(err)
	}

	actualGuildsBytes, err := readCloseResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	actualGuilds := make(sortableGuilds, 0)

	err = json.Unmarshal(actualGuildsBytes, &actualGuilds)
	if err != nil {
		t.Fatalf("Error unmarshaling actual guild data: %s", err)
	}

	if !reflect.DeepEqual(actualGuilds, expectedGuilds) {
		t.Errorf(
			"Unexpected response:\nGot:\n%s\n\nExpected:\n%s",
			string(actualGuildsBytes),
			string(expectedGuildsBytes),
		)
	}
}
