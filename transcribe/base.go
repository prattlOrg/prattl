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
	api_key, api_key_ok := os.LookupEnv("OPENAI_API_KEY")

	if !api_key_ok {
		return "OpenAI API key not found"
	}

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
	res, error := client.Do(req)

	if error != nil {
		return fmt.Sprintf("Error making request: %s", error)
	}

	return fmt.Sprintf("%v", res.Request.Body)
}
