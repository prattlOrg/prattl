package transcribe

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

const speechToTextUrl string = "https://api.openai.com/v1/audio/transcriptions"

func Test() string {
	api_key, _ := os.LookupEnv("OPENAI_API_KEY")

	jsonBody := []byte(`
        "data": {
            "file": "@./test.mp3",
            "model": "whisper-1"
          }`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest(http.MethodPost, speechToTextUrl, bodyReader)
	req.Header.Set("Content-Type", "mulitpart/form-data")
	req.Header.Set("Authorization", fmt.Sprintf("BEARER %s", api_key))
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, _ := client.Do(req)
	return fmt.Sprintf("%d", res.Request.Body)
}
