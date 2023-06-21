package notionagent

import (
	"context"
	"log"

	"github.com/jomei/notionapi"
)

func CreateNewPageInDatabase(client *notionapi.Client, database notionapi.DatabaseID, content map[string]notionapi.Property) (*notionapi.Page, error) {
	req := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: database,
		},
		Properties: content,
	}
	resp, err := client.Page.Create(context.Background(), &req)
	if err != nil {
		log.Printf("Failed to create page: %v", err)
		return nil, err
	}
	log.Printf("Created page: %v", resp)
	return resp, nil
}

func CreateNewPageWithBlockInDatabase(client *notionapi.Client, database notionapi.DatabaseID, content map[string]notionapi.Property, blocks []notionapi.Block) (*notionapi.Page, error) {
	req := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: database,
		},
		Properties: content,
	}
	if len(blocks) > 0 {
		for _, block := range blocks {
			if d, ok := block.(notionapi.ImageBlock); ok {
				req.Cover = &d.Image
				req.Cover.Caption = []notionapi.RichText{}
				break
			}
		}
		req.Children = blocks
	}
	resp, err := client.Page.Create(context.Background(), &req)
	if err != nil {
		log.Printf("Failed to create page: %v", err)
		return nil, err
	}
	log.Printf("Created page: %v", resp)
	return resp, nil
}
