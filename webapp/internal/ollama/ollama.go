package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func Embed(text string) ([]float32, error) {
	base := os.Getenv("OLLAMA_API_URL")
	uri := fmt.Sprintf("%s/api/embeddings", base)

	requestBody := EmbedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := cli.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned %s", resp.Status)
	}

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var embedResp EmbedResponse

	err = json.Unmarshal(respBytes, &embedResp)

	if err != nil {
		return nil, err
	}

	return embedResp.Embedding, nil

}
