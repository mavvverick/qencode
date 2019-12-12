package qencode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type QueryParams struct {
	Query *Query `json:"query"`
}

type Query struct {
	Source      string    `json:"source,omitempty"`
	CallbackURL string    `json:"callback_url,omitempty"`
	Format      *[]Format `json:"format,omitempty"`
}

type Format struct {
	Output               string                `json:"output,omitempty"`
	Size                 string                `json:"size,omitempty"`
	Bitrate              int64                 `json:"bitrate,omitempty"`
	VideoCodec           string                `json:"video_codec,omitempty"`
	Destination          *Destination          `json:"destination,omitempty"`
	Framerate            string                `json:"framerate,omitempty"`
	Keyframe             string                `json:"keyframe,omitempty"`
	FileExtension        string                `json:"file_extension,omitempty"`
	VideoCodecParameters *VideoCodecParameters `json:"video_codec_parameters,omitempty"`
	StartTime            string                `json:"start_time,omitempty"`
	Duration             string                `json:"duration,omitempty"`
	AudioBitrate         string                `json:"audio_bitrate,omitempty"`
	AudioSampleRate      string                `json:"audio_sample_rate,omitempty"`
	AudioChannelsNumber  string                `json:"audio_channels_number,omitempty"`
	Logo                 *Logo                 `json:"logo,omitempty"`
}

type VideoCodecParameters struct {
	Vprofile   string `json:"vprofile,omitempty"`
	Level      string `json:"level,omitempty"`
	Coder      string `json:"coder,omitempty"`
	Flags2     string `json:"flags2,omitempty"`
	Partitions string `json:"partitions,omitempty"`
	Directpred string `json:"directpred,omitempty"`
	MeMethod   string `json:"me_method,omitempty"`
	BStrategy  string `json:"b_strategy,omitempty"`
}

type Destination struct {
	Key         string `json:"key,omitempty"`
	Secret      string `json:"secret,omitempty"`
	URL         string `json:"url,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

type Logo struct {
	Source string `json:"source,omitempty"`
	X      string `json:"x,omitempty"`
	Y      string `json:"y,omitempty"`
}

func QueryBuilder(params *TaskParams, t *TaskServiceOp) (string, error) {

	q := &QueryParams{
		Query: &Query{
			Source:      fmt.Sprintf("https://%v.storage.googleapis.com/%v", t.client.Bucket, params.SourcePath),
			CallbackURL: t.client.CallbackURL,
			Format:      &[]Format{},
		},
	}

	for _, resolution := range params.Resolutions {
		data := strings.Split(resolution, "|")
		reso := data[0]
		bitrate, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			return "", err
		}

		form := Format{
			Output:              "mp4",
			Bitrate:             bitrate,
			Size:                Resolutions[reso],
			VideoCodec:          "libx265",
			Framerate:           "30",
			FileExtension:       "mp4",
			StartTime:           params.StartTime,
			Duration:            params.Duration,
			AudioBitrate:        "64",
			AudioSampleRate:     "44100",
			AudioChannelsNumber: "2",
			Destination: &Destination{
				Key:         t.client.StorageKey,
				Secret:      t.client.StorageSecret,
				URL:         fmt.Sprintf("s3://storage.googleapis.com/%v/%v/yo%v/%v", t.client.Bucket, params.PostID, reso, params.Name),
				Permissions: "public-read",
			},
			VideoCodecParameters: &VideoCodecParameters{
				Vprofile:   "high",
				Coder:      "0",
				Level:      "31",
				BStrategy:  "2",
				Flags2:     "-bpyramid+fastpskip-dct8x8",
				Partitions: "+parti8x8+parti4x4+partp8x8+partb8x8",
				Directpred: "2",
				MeMethod:   "hex",
			},
		}

		if params.IsH264 {
			form.VideoCodec = "libx264"
		}

		if reso == "web" {
			form.VideoCodec = "libx264"
			form.Size = Resolutions["540p"]
			form.Destination.URL = fmt.Sprintf("s3://storage.googleapis.com/%v/%v/%v", t.client.Bucket, params.PostID, params.Name)
			form.Logo = &Logo{
				Source: t.client.WATERMARK,
				X:      "12",
				Y:      "12",
			}
		}

		*q.Query.Format = append(*q.Query.Format, form)
	}

	data, err := json.Marshal(q)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
