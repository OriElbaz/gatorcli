package rss


import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFeed(t *testing.T) {
	type testCase struct {
		name          string
		mockStatus    int
		mockBody      string
		expectErr     bool
		expectedTitle string
		expectedItems int
	}

	tests := []testCase{
		{
			name:       "Valid Feed",
			mockStatus: http.StatusOK,
			mockBody: `<?xml version="1.0" encoding="UTF-8"?>
				<rss><channel>
					<title>Go Blog</title>
					<item><title>Post 1</title></item>
					<item><title>Post 2</title></item>
				</channel></rss>`,
			expectErr:     false,
			expectedTitle: "Go Blog",
			expectedItems: 2,
		},
		{
			name:       "Malformed XML",
			mockStatus: http.StatusOK,
			mockBody:   `<rss><channel><title>Broken XML`,
			expectErr:  true,
		},
		{
			name:       "Server Error 500",
			mockStatus: http.StatusInternalServerError,
			mockBody:   `Internal Server Error`,
			expectErr:  true,
		},
	}

	fmt.Println("\n--- Starting RSS Feed Tests ---")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.mockStatus)
				w.Write([]byte(tc.mockBody))
			}))
			defer server.Close()

			feed, err := FetchFeed(context.Background(), server.URL)

			// Determine Success/Failure
			errOccurred := (err != nil)
			passed := (errOccurred == tc.expectErr)

			// Logic check for non-error cases
			if !tc.expectErr && passed {
				if feed.Channel.Title != tc.expectedTitle || len(feed.Channel.Item) != tc.expectedItems {
					passed = false
				}
			}

			// Printing the results in a readable format
			statusSymbol := "✅"
			if !passed {
				statusSymbol = "❌"
			}

			fmt.Printf("%s Test Case: %s\n", statusSymbol, tc.name)
			
			// Detail "Got vs Expected" for Titles
			if !tc.expectErr {
				gotTitle := "nil"
				if feed != nil {
					gotTitle = feed.Channel.Title
				}
				fmt.Printf("   - Title:    [Got: %-10s | Expected: %s]\n", gotTitle, tc.expectedTitle)
				
				gotItems := 0
				if feed != nil {
					gotItems = len(feed.Channel.Item)
				}
				fmt.Printf("   - Items:    [Got: %-10d | Expected: %d]\n", gotItems, tc.expectedItems)
			} else {
				fmt.Printf("   - Error:    [Got Error: %-5t | Expected Error: %t]\n", errOccurred, tc.expectErr)
			}

			// Final Assertion for the Test Runner
			if !passed {
				t.Errorf("%s failed", tc.name)
			}
		})
	}
	fmt.Print("-------------------------------\n")
}