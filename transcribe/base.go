package transcribe

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type WhisperResponse struct {
	// optional fields
	WhisperError   WhisperError `json:"error"`
	WhisperSuccess string       `json:"text"`
}

type WhisperError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

func TranscribeWhisperApi() WhisperResponse {
	fmt.Println("transcribing...")

	const speechToTextUrl string = "https://api.openai.com/v1/audio/transcriptions"
	api_key, api_key_ok := os.LookupEnv("OPENAI_API_KEY")
	if !api_key_ok {
		fmt.Println("OpenAI API key not found")
	}

	// go routine for copying file direclty to request instead of loading entire file in memory
	// https://stackoverflow.com/questions/77091845/golang-uploading-big-file-to-external-api-with-multipart-how-to-avoid-io-copy
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		_ = writer.WriteField("model", "whisper-1")
		file, err := os.Open("templates/assets/test.mp3")
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer file.Close()
		part3, err := writer.CreateFormFile("file", "test.mp3")
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part3, file)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.CloseWithError(writer.Close())
	}()

	request, _ := http.NewRequest("POST", speechToTextUrl, pr)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
	}
	defer response.Body.Close()

	var whisperResponse WhisperResponse
	json.Unmarshal(responseBody, &whisperResponse)
	return whisperResponse
}
