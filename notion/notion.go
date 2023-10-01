package notion

import "github.com/dstotijn/go-notion"

type Notion struct {
	Client *notion.Client
	PageId string
}

func NewNotion(secret, pageId string) *Notion {
	return &Notion{
		Client: notion.NewClient(secret),
		PageId: pageId,
	}
}
