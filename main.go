package main

import (
	"AwesomeEino/embed"
	"AwesomeEino/stage4"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load() // 加载环境变量
	if err != nil {
		log.Fatal("Error loading .env file") // 处理加载错误
	}
}
func main() {
	stage4.DemoInsertToMilvus()
}

func TestIndexerRAG() {
	println("=== Indexer RAG Example ===")
	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: "This is the content of document 1.",
			MetaData: map[string]any{
				"source": "source1",
			},
		},
		{
			ID:      "doc2",
			Content: "This is the content of document 2.",
			MetaData: map[string]any{
				"source": "source2",
			},
		},
	}
	stage4.IndexerRAG(docs)
}

func TestEmbedStrings() {
	// Example 1: Embed text strings
	fmt.Println("=== Text Embedding Example ===")
	result, err := embed.EmbedStrings([]string{"Hello, world!", "Multimodal embeddings are powerful."})
	if err != nil {
		log.Fatalf("Text embedding error: %v", err)
	}

	//写入文件
	jsonIndentWrite("result.json", result)

	// Example 2: Show specific fields
	fmt.Println("\n=== Key Fields ===")
	fmt.Printf("Model: %s\n", result.Model)
	fmt.Printf("Created: %d\n", result.Created)
	fmt.Printf("Object: %s\n", result.Object)
	fmt.Printf("ID: %s\n", result.ID)
	fmt.Printf("Usage - Total Tokens: %d\n", result.Usage.TotalTokens)
}

func jsonIndentWrite(filename string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0644)
}
