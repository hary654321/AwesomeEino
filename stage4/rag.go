package stage4

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/cloudwego/eino/schema"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// DemoInsertToMilvus is a small runnable demo showing how to create a collection
// and insert simple documents (id, content, metadata, vector) using the
// official Milvus Go SDK. It uses environment variable MILVUS_ADDRESS (e.g.
// "localhost:19530").
func DemoInsertToMilvus() {
	ctx := context.Background()

	addr := os.Getenv("MILVUS_ADDRESS")
	if addr == "" {
		addr = "localhost:19530"
	}

	// Create client
	client, err := milvus.NewGrpcClient(ctx, addr)
	if err != nil {
		log.Fatalf("failed to create milvus client: %v", err)
	}
	defer client.Close()

	collectionName := "AwesomeEino"

	// Define schema: id (varchar PK), vector (float vector), content, metadata
	dim := 1536
	fields := []*entity.Field{
		{
			Name:       "id",
			DataType:   entity.FieldTypeVarChar,
			PrimaryKey: true,
			TypeParams: map[string]string{"max_length": "255"},
		},
		{
			Name:     "vector",
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": fmt.Sprintf("%d", dim),
			},
		},
		{
			Name:       "content",
			DataType:   entity.FieldTypeVarChar,
			TypeParams: map[string]string{"max_length": "4096"},
		},
		{
			Name:       "metadata",
			DataType:   entity.FieldTypeVarChar,
			TypeParams: map[string]string{"max_length": "4096"},
		},
	}

	// If collection exists, drop it for demo cleanliness (optional)
	exists, err := client.HasCollection(ctx, collectionName)
	if err != nil {
		log.Fatalf("HasCollection error: %v", err)
	}
	if exists {
		if err := client.DropCollection(ctx, collectionName); err != nil {
			log.Fatalf("failed to drop existing collection: %v", err)
		}
		fmt.Println("Dropped existing collection for demo")
	}

	// return
	schema1 := &entity.Schema{
		CollectionName: collectionName,
		AutoID:         false,
		Fields:         fields,
	}

	if err := client.CreateCollection(ctx, schema1, 2); err != nil {
		log.Fatalf("CreateCollection error: %v", err)
	}
	fmt.Println("Created collection:", collectionName)

	// Prepare some sample documents
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

	// Build entities for insertion
	ids := make([]string, 0, len(docs))
	contents := make([]string, 0, len(docs))
	metas := make([]string, 0, len(docs))
	vectors := make([][]float32, 0, len(docs))

	rand.Seed(time.Now().UnixNano())
	for _, d := range docs {
		ids = append(ids, d.ID)
		contents = append(contents, d.Content)
		mj, _ := json.Marshal(d.MetaData)
		metas = append(metas, string(mj))

		// For demo purposes generate random vectors. In real use replace with
		// embedding vectors from your embedder (e.g., Doubao embedding).
		vec := make([]float32, dim)
		for i := 0; i < dim; i++ {
			vec[i] = rand.Float32()
		}
		vectors = append(vectors, vec)
	}

	// Convert to milvus entity columns
	idCol := entity.NewColumnVarChar("id", ids)
	contentCol := entity.NewColumnVarChar("content", contents)
	metaCol := entity.NewColumnVarChar("metadata", metas)
	vecCol := entity.NewColumnFloatVector("vector", dim, vectors)

	// Insert
	res, err := client.Insert(ctx, collectionName, "", idCol, vecCol, contentCol, metaCol)
	// 调试版（最详细）
	fmt.Printf("insert result ---->\n%#v\n", res)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

	// Flush to make data persistent
	if err := client.Flush(ctx, collectionName, false); err != nil {
		log.Fatalf("Flush error: %v", err)
	}

	// 1. 创建索引（在 Flush 之后、LoadCollection 之前）
	// 为浮点向量字段建索引
	idx := entity.NewGenericIndex(
		"ivec", // 索引名随意
		entity.IvfFlat,
		map[string]string{
			"nlist":       "128",
			"metric_type": "L2", // 或 "IP"/"COSINE"，按业务选
		})
	if err := client.CreateIndex(ctx, collectionName, "vector", idx, false); err != nil {
		log.Fatalf("CreateIndex error: %v", err)
	}
	if err := client.LoadCollection(ctx, collectionName, true); err != nil {
		log.Fatalf("LoadCollection error: %v", err)
	}
	fmt.Println("Inserted demo documents into collection:", collectionName)
}
