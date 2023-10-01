package main

import (
	"context"
	"github.com/dstotijn/go-notion"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"log"
	"os"
	notionpkg "ryohei-takagi/sqlboiler-notion-sample/notion"
	"ryohei-takagi/sqlboiler-notion-sample/sqlboiler"
)

type Column struct {
	Name     string
	DBType   string
	PK       bool
	Unique   bool
	Nullable bool
	Default  string
	Comment  string
}

func NewColumn(c drivers.Column, p *drivers.PrimaryKey) *Column {
	var pk bool

	if p != nil {
		for _, p := range p.Columns {
			if p == c.Name {
				pk = true
				break
			}
		}
	}

	return &Column{
		PK:       pk,
		Name:     c.Name,
		DBType:   c.DBType,
		Unique:   c.Unique,
		Nullable: c.Nullable,
		Default:  c.Default,
		Comment:  c.Comment,
	}
}

func main() {
	var (
		notionSecret = os.Getenv("NOTION_SECRET")
		notionPageId = os.Getenv("NOTION_DB_DOC_PAGE_ID")
	)

	ctx := context.Background()

	// Notion init
	n := notionpkg.NewNotion(notionSecret, notionPageId)

	// SQLBoiler init
	b, err := sqlboiler.NewBoiler()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, table := range b.Tables {
		// Create Notion database page
		dbPage, err := n.Client.CreateDatabase(ctx, notion.CreateDatabaseParams{
			ParentPageID: n.PageId,
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: table.Name,
					},
				},
			},
			Description: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "", // description here
					},
				},
			},
			Properties: notion.DatabaseProperties{
				"PK": notion.DatabaseProperty{
					Type:     notion.DBPropTypeCheckbox,
					Checkbox: &notion.EmptyMetadata{},
				},
				"Column": notion.DatabaseProperty{
					Type:  notion.DBPropTypeTitle,
					Title: &notion.EmptyMetadata{},
				},
				"DBType": notion.DatabaseProperty{
					Type:     notion.DBPropTypeRichText,
					RichText: &notion.EmptyMetadata{},
				},
				"Unique": notion.DatabaseProperty{
					Type:     notion.DBPropTypeCheckbox,
					Checkbox: &notion.EmptyMetadata{},
				},
				"Nullable": notion.DatabaseProperty{
					Type:     notion.DBPropTypeCheckbox,
					Checkbox: &notion.EmptyMetadata{},
				},
				"Default": notion.DatabaseProperty{
					Type:     notion.DBPropTypeRichText,
					RichText: &notion.EmptyMetadata{},
				},
				"Comment": notion.DatabaseProperty{
					Type:     notion.DBPropTypeRichText,
					RichText: &notion.EmptyMetadata{},
				},
			},
		})
		if err != nil {
			log.Fatal(err.Error())
		}

		// Update Notion database table row
		for i := range table.Columns {
			c := table.Columns[len(table.Columns)-1-i]
			column := NewColumn(c, table.PKey)

			_, err := n.Client.CreatePage(ctx, notion.CreatePageParams{
				ParentType: notion.ParentTypeDatabase,
				ParentID:   dbPage.ID,
				DatabasePageProperties: &notion.DatabasePageProperties{
					"Column": notion.DatabasePageProperty{
						Type: notion.DBPropTypeTitle,
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: column.Name,
								},
							},
						},
					},
					"DBType": notion.DatabasePageProperty{
						Type: notion.DBPropTypeRichText,
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: column.DBType,
								},
							},
						},
					},
					"PK": notion.DatabasePageProperty{
						Type:     notion.DBPropTypeCheckbox,
						Checkbox: notion.BoolPtr(column.PK),
					},
					"Unique": notion.DatabasePageProperty{
						Type:     notion.DBPropTypeCheckbox,
						Checkbox: notion.BoolPtr(column.Unique),
					},
					"Nullable": notion.DatabasePageProperty{
						Type:     notion.DBPropTypeCheckbox,
						Checkbox: notion.BoolPtr(column.Nullable),
					},
					"Default": notion.DatabasePageProperty{
						Type: notion.DBPropTypeRichText,
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: column.Default,
								},
							},
						},
					},
					"Comment": notion.DatabasePageProperty{
						Type: notion.DBPropTypeRichText,
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: column.Comment,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
