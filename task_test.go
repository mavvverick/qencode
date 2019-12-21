package qencode

import (
	"fmt"
	"net/http"
	"os"
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

	//postIDS := []string{"0aoz-3aZa", "xoTk-3-WRa", "ZATkaqaZa", "xooka3aZa", "oboka3-Za", "10Tk-q-Wa", "Y3M7-3-Wa", "uAI4a3-Wa", "aGiVaq-Za", "lZX4-3aZa", "Agl7aqaZa", "vJeV-q-Wa", "evzvaq-Wa", "3MRF-qaZa"}
	postIDS := []string{"vJeV-q-Wa"}
	vids := []string{"intro.mp4", "1.mp4", "2.mp4"}
	//vids := []string{"intro.mp4"}
	for _, postID := range postIDS {
		for _, vid := range vids {
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
				Resolutions: []string{"540p|500", "240p|320"},
				Payload:     postID,
				// StartTime:   "0",
				// Duration:    "3",
			}

			// if i < 3 {
			// 	params.Duration = duration[i]
			// }

			encode, _, err := client.Task.Encode(ctx, params)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			fmt.Println(encode)
		}

	}

}
