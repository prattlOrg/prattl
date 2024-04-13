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

func TranscribeWhisperApi(file multipart.File) WhisperResponse {
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
		// file, err := os.Open("public/assets/test.mp3")
		// if err != nil {
		// 	pw.CloseWithError(err)
		// 	return
		// }
		// defer file.Close()
		part3, err := writer.CreateFormFile("file", "fileMade.wav")
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

	// // open
	// f, h, err := req.FormFile("q")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer f.Close()

	// // for your information
	// fmt.Println("\nfile:", f, "\nheader:", h, "\nerr", err)

	// // read
	// bs, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// s = string(bs)

	request, _ := http.NewRequest("POST", speechToTextUrl, pr)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	var whisperResponse WhisperResponse
	json.Unmarshal(responseBody, &whisperResponse)
	return whisperResponse
}
