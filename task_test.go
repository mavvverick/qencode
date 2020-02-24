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
		PostID:      "41ysdzbWg",
		Name:        "1.mp4",
		SourcePath:  fmt.Sprintf("%v/raw/%v", "41ysdzbWg", "1.mp4"),
		Resolutions: strings.Split(os.Getenv("QUENCODE_RESOLUTIONS"), ","),
		Payload:     "41ysdzbWg|1",
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
	postIDS := []string{"_-Cap6aZH"}
	vids := []string{"1.mp4", "2.mp4", "intro.mp4", "boomerang.mp4"}

	fmt.Println(os.Getenv("QUENCODE_RESOLUTIONS"))
	for _, postID := range postIDS {
		for _, vid := range vids {
			task, _, err := client.Task.Create(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			params := &TaskParams{
				TaskToken:   task.TaskToken,
				PostID:      postID,
				Name:        vid,
				SourcePath:  fmt.Sprintf("%v/raw/%v", postID, vid),
				Resolutions: strings.Split(os.Getenv("QUENCODE_RESOLUTIONS"), ","),
				Payload:     postID,
			}
			encode, _, err := client.Task.Encode(ctx, params)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			fmt.Println(encode)
		}

	}

}

func TestAction_Mvf(t *testing.T) {

	setup()
	defer teardown()
	mux.HandleFunc("/v1/quencode/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	fmt.Println(os.Getenv("QUENCODE_RESOLUTIONS"))
	sum := 0
	for i := 1; i <= 16; i++ {
		sum += i

		task, _, err := client.Task.Create(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		name := fmt.Sprintf("video_%v.mp4", i)
		folder := "mvf/2"

		params := &TaskParams{
			TaskToken:   task.TaskToken,
			PostID:      folder,
			Name:        name,
			SourcePath:  fmt.Sprintf("%v/%v", folder, name),
			Resolutions: strings.Split(os.Getenv("QUENCODE_RESOLUTIONS"), ","),
			Payload:     name,
			// StartTime:   "0",
			// Duration:    "3",
		}

		//fmt.Println(params)

		encode, _, err := client.Task.Encode(ctx, params)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		fmt.Println(encode)

	}

}
