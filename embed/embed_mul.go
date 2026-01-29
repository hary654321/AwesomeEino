package embed

import (
	"context"
	"fmt"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

// EmbeddingResult represents a structured embedding result
type EmbeddingResult struct {
	ID      string                         `json:"id"`
	Model   string                         `json:"model"`
	Created int64                          `json:"created"`
	Object  string                         `json:"object"`
	Data    model.MultimodalEmbedding      `json:"data"`
	Usage   model.MultimodalEmbeddingUsage `json:"usage"`
}

// DoubaoEmbeddingVision performs multimodal embedding on a vision input and returns structured result
func DoubaoEmbeddingVision() (*EmbeddingResult, error) {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
	)
	ctx := context.Background()

	fmt.Println("----- multimodal embeddings request -----")
	req := model.MultiModalEmbeddingRequest{
		Model: os.Getenv("EMBEDDER"),
		Input: []model.MultimodalEmbeddingInput{
			{
				Type:     model.MultiModalEmbeddingInputTypeImageURL,
				ImageURL: &model.MultimodalEmbeddingImageURL{URL: "https://ark-project.tos-cn-beijing.volces.com/images/view.jpeg"},
			},
		},
	}

	resp, err := client.CreateMultiModalEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("multimodal embeddings error: %w", err)
	}

	result := &EmbeddingResult{
		ID:      resp.Id,
		Model:   resp.Model,
		Created: resp.Created,
		Object:  resp.Object,
		Data:    resp.Data,
		Usage:   resp.Usage,
	}

	return result, nil
}

// EmbedStrings performs embedding on text strings and returns structured results
func EmbedStrings(txts []string) (*EmbeddingResult, error) {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
	)
	ctx := context.Background()

	var inputs []model.MultimodalEmbeddingInput
	for _, txt := range txts {
		textCopy := txt // Create a copy to avoid taking address of loop variable
		inputs = append(inputs, model.MultimodalEmbeddingInput{
			Type: model.MultiModalEmbeddingInputTypeText,
			Text: &textCopy,
		})
	}

	req := model.MultiModalEmbeddingRequest{
		Model: os.Getenv("EMBEDDER"),
		Input: inputs,
	}

	resp, err := client.CreateMultiModalEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("text embeddings error: %w", err)
	}

	result := &EmbeddingResult{
		ID:      resp.Id,
		Model:   resp.Model,
		Created: resp.Created,
		Object:  resp.Object,
		Data:    resp.Data,
		Usage:   resp.Usage,
	}

	return result, nil
}

// EmbedImages performs embedding on image URLs and returns structured results
func EmbedImages(imageUrls []string) (*EmbeddingResult, error) {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
	)
	ctx := context.Background()

	var inputs []model.MultimodalEmbeddingInput
	for _, url := range imageUrls {
		inputs = append(inputs, model.MultimodalEmbeddingInput{
			Type:     model.MultiModalEmbeddingInputTypeImageURL,
			ImageURL: &model.MultimodalEmbeddingImageURL{URL: url},
		})
	}

	req := model.MultiModalEmbeddingRequest{
		Model: os.Getenv("EMBEDDER"),
		Input: inputs,
	}

	resp, err := client.CreateMultiModalEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("image embeddings error: %w", err)
	}

	result := &EmbeddingResult{
		ID:      resp.Id,
		Model:   resp.Model,
		Created: resp.Created,
		Object:  resp.Object,
		Data:    resp.Data,
		Usage:   resp.Usage,
	}

	return result, nil
}
