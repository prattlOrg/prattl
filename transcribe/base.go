package transcribe

import (
	// "bytes"
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

func TranscribeWhisperApi(form *multipart.Form) WhisperResponse {
	fmt.Println("transcribing...")

	const speechToTextUrl string = "https://api.openai.com/v1/audio/transcriptions"
	api_key, api_key_ok := os.LookupEnv("OPENAI_API_KEY")
	if !api_key_ok {
		fmt.Println("OpenAI API key not found")
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		_ = writer.WriteField("model", "whisper-1")
		for _, files := range form.File {
			for _, file_handler := range files {
				fmt.Printf("File: %s\n", file_handler.Filename)
				// part, _ := writer.CreateFormFile("file", file_handler.Filename)
				file, _ := file_handler.Open()

				dst, err := os.Create(file_handler.Filename)
				if err != nil {
					fmt.Println("error creating file", err)
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer dst.Close()
				if _, err := io.Copy(dst, file); err != nil {
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Printf("uploaded file\n")

				defer file.Close()
				// _, _ = io.Copy(part, file)
			}
		}
		pw.CloseWithError(writer.Close())
	}()

	fmt.Println("Out of routine")

	request, _ := http.NewRequest("POST", speechToTextUrl, pr)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	fmt.Println(fmt.Sprintf("Sending req %s", request.Form))
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Println("Got response")
	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(
			fmt.Sprintf("Error reading response: %s", error),
		)
	}
	defer response.Body.Close()

	var whisperResponse WhisperResponse
	e := json.Unmarshal(responseBody, &whisperResponse)

	if e != nil {
		fmt.Println(
			fmt.Sprintf("Error umarshalling json response: %s", error),
		)
	}
	return whisperResponse
}
