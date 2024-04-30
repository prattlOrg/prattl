package transcribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
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
	println(request)
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
	e := json.Unmarshal(responseBody, &whisperResponse)

	if e != nil {
		fmt.Println(
			fmt.Sprintf("Error umarshalling json response: %s", error),
		)
	}
	return whisperResponse
}

func TranscribeLocal() {
	// cmd := exec.Command("ffmpeg")
	cmd := exec.Command("python3", "transcribe/transcribe.py")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
