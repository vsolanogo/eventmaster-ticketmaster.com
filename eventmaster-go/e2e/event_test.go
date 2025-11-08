package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventEndpoints(t *testing.T) {
	runSubtest(t, "list events", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/events", nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var list EventListResponse
		decodeJSON(t, resp.Body, &list)

		if list.Events == nil {
			t.Fatalf("expected events field to be present")
		}
	})

	runSubtest(t, "pagination", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/events?page=1&limit=5", nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var list EventListResponse
		decodeJSON(t, resp.Body, &list)

		if len(list.Events) > 5 {
			t.Fatalf("expected at most 5 events, got %d", len(list.Events))
		}
	})

	runSubtest(t, "sort asc", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/events?sortBy=eventDate&sortOrder=ASC", nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var list EventListResponse
		decodeJSON(t, resp.Body, &list)

		if len(list.Events) > 1 {
			first := list.Events[0].EventDate
			second := list.Events[1].EventDate
			if first.After(*second) {
				t.Fatalf("expected events to be sorted ascending")
			}
		}
	})

	runSubtest(t, "sort desc", func(t *testing.T) {
		resp := doRequest(t, http.MethodGet, "/events?sortBy=eventDate&sortOrder=DESC", nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var list EventListResponse
		decodeJSON(t, resp.Body, &list)

		if len(list.Events) > 1 {
			first := list.Events[0].EventDate
			second := list.Events[1].EventDate
			if first.Before(*second) {
				t.Fatalf("expected events to be sorted descending")
			}
		}
	})

	runSubtest(t, "get by id", func(t *testing.T) {
		listResp := doRequest(t, http.MethodGet, "/events", nil, nil)
		defer listResp.Body.Close()

		var list EventListResponse
		decodeJSON(t, listResp.Body, &list)

		if len(list.Events) == 0 {
			t.Skip("no events available to fetch")
		}

		id := list.Events[0].ID
		resp := doRequest(t, http.MethodGet, fmt.Sprintf("/events/%s", id), nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var event EventResponse
		decodeJSON(t, resp.Body, &event)

		if event.ID != id {
			t.Fatalf("expected event id %s, got %s", id, event.ID)
		}
	})

	runSubtest(t, "missing event returns 404", func(t *testing.T) {
		fakeID := uuid.NewString()
		resp := doRequest(t, http.MethodGet, fmt.Sprintf("/events/%s", fakeID), nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	runSubtest(t, "create event requires auth", func(t *testing.T) {
		resetCookies(t)
		payload := map[string]any{
			"title":       "Test Event",
			"description": "Description",
			"organizer":   "Org",
			"latitude":    40.0,
			"longitude":   -70.0,
			"eventDate":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}

		resp := doRequest(t, http.MethodPost, "/events", payload, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}
