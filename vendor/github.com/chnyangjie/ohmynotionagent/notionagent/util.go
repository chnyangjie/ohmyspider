package notionagent

import (
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

func newNumber(content interface{}) (*notionapi.NumberProperty, error) {
	if d, ok := content.(float64); ok {
		return &notionapi.NumberProperty{Number: d}, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}
func newCheckbox(content interface{}) (*notionapi.CheckboxProperty, error) {
	if d, ok := content.(bool); ok {
		return &notionapi.CheckboxProperty{Checkbox: d}, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}
func genRichTextObj(content []string) []notionapi.RichText {
	result := []notionapi.RichText{}
	for _, item := range content {
		if item == "" {
			continue
		}
		r := []rune(item)
		start := 0
		for start < len(r) {
			if len(r) <= start+1900 {
				result = append(result, notionapi.RichText{Text: &notionapi.Text{Content: string(r[start:])}})
			} else {
				result = append(result, notionapi.RichText{Text: &notionapi.Text{Content: string(r[start : start+1900])}})
			}
			start += 1900
		}
	}
	if len(result) == 0 {
		result = append(result, notionapi.RichText{Text: &notionapi.Text{Content: ""}})
	}
	return result
}
func newRichText(content interface{}) (*notionapi.RichTextProperty, error) {
	result := notionapi.RichTextProperty{}
	raw := []string{}
	if d, ok := content.(string); ok {
		raw = append(raw, d)
	} else if d, ok := content.([]string); ok {
		raw = append(raw, d...)
	}
	result.RichText = append(result.RichText, genRichTextObj(raw)...)
	if len(result.RichText) == 0 {
		return nil, fmt.Errorf("unsupport content type: %+v", content)
	}
	return &result, nil
}

func newRelation(content interface{}) (*notionapi.RelationProperty, error) {
	if d, ok := content.([]string); ok {
		r := notionapi.RelationProperty{Relation: []notionapi.Relation{}}
		for _, item := range d {
			r.Relation = append(r.Relation, notionapi.Relation{ID: notionapi.PageID(item)})
		}
		return &r, nil
	} else if d, ok := content.(string); ok {
		r := notionapi.RelationProperty{Relation: []notionapi.Relation{}}
		r.Relation = append(r.Relation, notionapi.Relation{ID: notionapi.PageID(d)})
		return &r, nil
	} else if d, ok := content.(notionapi.ObjectID); ok {
		r := notionapi.RelationProperty{Relation: []notionapi.Relation{}}
		r.Relation = append(r.Relation, notionapi.Relation{ID: notionapi.PageID(d)})
		return &r, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}

func newTitle(content interface{}) (*notionapi.TitleProperty, error) {
	if d, ok := content.(string); ok {
		return &notionapi.TitleProperty{Title: []notionapi.RichText{{Text: &notionapi.Text{Content: d}}}}, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}
func newUrl(content interface{}) (*notionapi.URLProperty, error) {
	if d, ok := content.(string); ok {
		if strings.HasPrefix(strings.ToLower(d), "http") {
			return &notionapi.URLProperty{URL: d}, nil
		}
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}

func newMultiSelect(content interface{}) (*notionapi.MultiSelectProperty, error) {
	if d, ok := content.([]string); ok {
		r := notionapi.MultiSelectProperty{MultiSelect: []notionapi.Option{}}
		for _, item := range d {
			r.MultiSelect = append(r.MultiSelect, notionapi.Option{Name: item})
		}
		return &r, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}

func newDate(content interface{}) (*notionapi.DateProperty, error) {
	if d, ok := content.([]time.Time); ok {
		if len(d) == 0 {
			return nil, fmt.Errorf("unsupport content type: %+v", content)
		}
		s := notionapi.Date(d[0])
		r := &notionapi.DateProperty{Date: &notionapi.DateObject{Start: &s}}
		if len(d) > 1 {
			e := notionapi.Date(d[1])
			r.Date.End = &e
		} else {
			r.Date.End = &s
		}
		return r, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}
func newTableBlock(content interface{}) ([]notionapi.Block, error) {
	if d, ok := content.([][]string); !ok {
		return nil, fmt.Errorf("unsupport content type: %+v", content)
	} else {
		if len(d) == 0 {
			return nil, fmt.Errorf("empty content: %+v", content)
		}
		rows := []notionapi.Block{}
		for _, rawRow := range d {
			row := notionapi.TableRowBlock{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeTableRowBlock,
					Object: notionapi.ObjectTypeBlock,
				},
				TableRow: notionapi.TableRow{
					Cells: [][]notionapi.RichText{},
				},
			}
			for _, rawCell := range rawRow {
				cell := genRichTextObj([]string{rawCell})
				row.TableRow.Cells = append(row.TableRow.Cells, cell)
			}
			rows = append(rows, row)
		}
		block := []notionapi.Block{
			notionapi.TableBlock{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeTableBlock,
					Object: notionapi.ObjectTypeBlock,
				},
				Table: notionapi.Table{
					TableWidth:      len(d[0]),
					HasColumnHeader: true,
					Children:        rows,
				},
			},
		}
		return block, nil
	}
}

func newImage(content interface{}) ([]notionapi.Block, error) {
	if d, ok := content.([]string); !ok {
		return nil, fmt.Errorf("unsupport content type: %+v", content)
	} else {
		return []notionapi.Block{
			notionapi.ImageBlock{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeImage,
					Object: notionapi.ObjectTypeBlock,
				},
				Image: notionapi.Image{
					Type: notionapi.FileTypeExternal,
					External: &notionapi.FileObject{
						URL: d[0],
					},
					Caption: genRichTextObj(d[1:]),
				},
			},
		}, nil
	}
}
func newQuote(content interface{}) ([]notionapi.Block, error) {
	result := []notionapi.QuoteBlock{
		{
			BasicBlock: notionapi.BasicBlock{
				Type:   notionapi.BlockQuote,
				Object: notionapi.ObjectTypeBlock,
			},
		},
	}
	raw := []string{}
	if d, ok := content.(string); ok {
		raw = append(raw, d)
	} else if d, ok := content.([]string); ok {
		raw = append(raw, d...)
	}
	count := 0
	raws := genRichTextObj(raw)
	for _, item := range raws {
		count += 1
		if count > 100 {
			result = append(result, notionapi.QuoteBlock{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeParagraph,
					Object: notionapi.ObjectTypeBlock,
				},
			})
			count = 1
		}
		result[len(result)-1].Quote.RichText = append(result[len(result)-1].Quote.RichText, item)
	}
	r := []notionapi.Block{}
	for _, item := range result {
		r = append(r, notionapi.Block(item))
	}
	return r, nil
}
func newParagraph(content interface{}) ([]notionapi.Block, error) {
	result := []notionapi.ParagraphBlock{
		{
			BasicBlock: notionapi.BasicBlock{
				Type:   notionapi.BlockTypeParagraph,
				Object: notionapi.ObjectTypeBlock,
			},
		},
	}
	raw := []string{}
	if d, ok := content.(string); ok {
		raw = append(raw, d)
	} else if d, ok := content.([]string); ok {
		raw = append(raw, d...)
	}
	count := 0
	raws := genRichTextObj(raw)
	for _, item := range raws {
		count += 1
		if count > 100 {
			result = append(result, notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeParagraph,
					Object: notionapi.ObjectTypeBlock,
				},
			})
			count = 1
		}
		result[len(result)-1].Paragraph.RichText = append(result[len(result)-1].Paragraph.RichText, item)
	}
	r := []notionapi.Block{}
	for _, item := range result {
		r = append(r, notionapi.Block(item))
	}
	return r, nil
}
func newColumn(blocks []notionapi.Block) (notionapi.Block, error) {
	res := notionapi.ColumnListBlock{
		BasicBlock: notionapi.BasicBlock{
			Type:   notionapi.BlockTypeColumnList,
			Object: notionapi.ObjectTypeBlock,
		},
		ColumnList: notionapi.ColumnList{},
	}
	for _, item := range blocks {
		column := notionapi.ColumnBlock{
			BasicBlock: notionapi.BasicBlock{
				Type:   notionapi.BlockTypeColumn,
				Object: notionapi.ObjectTypeBlock,
			},
			Column: notionapi.Column{
				Children: []notionapi.Block{item},
			},
		}
		res.ColumnList.Children = append(res.ColumnList.Children, column)
	}
	return res, nil
}

func newHeading(content interface{}, blockType notionapi.BlockType) ([]notionapi.Block, error) {
	result := []notionapi.Block{}
	if d, ok := content.(string); ok {
		if blockType == notionapi.BlockTypeHeading1 {
			result = append(result, notionapi.Heading1Block{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeHeading1,
					Object: notionapi.ObjectTypeBlock,
				},
				Heading1: notionapi.Heading{
					RichText: genRichTextObj([]string{d}),
				},
			})
		} else if blockType == notionapi.BlockTypeHeading2 {
			result = append(result, notionapi.Heading2Block{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeHeading2,
					Object: notionapi.ObjectTypeBlock,
				},
				Heading2: notionapi.Heading{
					RichText: genRichTextObj([]string{d}),
				},
			})
		} else if blockType == notionapi.BlockTypeHeading3 {
			result = append(result, notionapi.Heading3Block{
				BasicBlock: notionapi.BasicBlock{
					Type:   notionapi.BlockTypeHeading3,
					Object: notionapi.ObjectTypeBlock,
				},
				Heading3: notionapi.Heading{
					RichText: genRichTextObj([]string{d}),
				},
			})
		}

	} else {
		return nil, fmt.Errorf("unsupport content type: %+v", content)
	}
	return result, nil
}
