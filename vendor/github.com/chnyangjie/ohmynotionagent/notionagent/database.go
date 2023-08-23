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
	for k, v := range content {
		if v == nil {
			delete(content, k)
		}
	}
	for i, block := range blocks {
		if block == nil {
			blocks = append(blocks[:i], blocks[i+1:]...)
		}
	}
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
	if icon, ok := content["Icon"]; ok {
		if s, ok := icon.(*notionapi.URLProperty); ok {
			req.Icon = &notionapi.Icon{
				Type: notionapi.FileTypeExternal,
				External: &notionapi.FileObject{
					URL: s.URL,
				},
			}
		}
	}
	resp, err := client.Page.Create(context.Background(), &req)
	if err != nil {
		log.Printf("Failed to create page: %v", err)
		return nil, err
	}
	log.Printf("Created page: %v", resp)
	return resp, nil
}
