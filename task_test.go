package qencode

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestAction_GetSession(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/v1/quencode/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	token, _, err := client.Task.SessionToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	fmt.Println(token)
}

func TestAction_Create(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/v1/quencode/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	task, _, err := client.Task.Create(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	params := &TaskParams{
		TaskToken:   task.TaskToken,
		SourcePath:  "41ysdzbWg/1.mp4",
		FinalPath:   "41ysdzbWg/7.mp4",
		Resolutions: []string{Resolutions["540p"], Resolutions["240p"]},
		Payload:     fmt.Sprintf(`%v|%v`, os.Getenv("QENCODE_WEBHOOK_ACCESS"), "41ysdzbWg"),
		StartTime:   "0",
		Duration:    "80",
	}

	encode, _, err := client.Task.Encode(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	fmt.Println(encode)
}

func TestAction_MuiltipleCreate(t *testing.T) {

	setup()
	defer teardown()
	mux.HandleFunc("/v1/quencode/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	//postID := "-0q088aZg"
	postID := "-1Evzw-Zg"
	duration := strings.Split("9|6|1", "|")
	//vids := []string{"intro.mp4", "1.mp4", "2.mp4", "boomerang.mp4"}
	vids := []string{"intro.mp4"}

	for i, vid := range vids {
		task, _, err := client.Task.Create(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		params := &TaskParams{
			TaskToken:  task.TaskToken,
			PostID:     postID,
			Name:       vid,
			SourcePath: fmt.Sprintf("%v/raw/%v", postID, vid),
			//Resolutions: []string{"540p|1500", "240p|600", "web|1500"},
			Resolutions: []string{"web|1500"},
			Payload:     postID,
			StartTime:   "0",
			Duration:    "3",
		}

		if i < 3 {
			params.Duration = duration[i]
		}

		encode, _, err := client.Task.Encode(ctx, params)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		fmt.Println(encode)
	}
}
