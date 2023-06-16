package notionagent

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

func PropertyToVar(property notionapi.Property) (interface{}, error) {
	switch p := property.(type) {
	case *notionapi.RichTextProperty:
		{
			content := []string{}
			for _, item := range p.RichText {
				content = append(content, item.Text.Content)
			}
			return content, nil
		}
	case *notionapi.TitleProperty:
		{
			content := []string{}
			for _, item := range p.Title {
				content = append(content, item.Text.Content)
			}
			return content, nil
		}
	case *notionapi.MultiSelectProperty:
		{
			content := []string{}
			for _, item := range p.MultiSelect {
				content = append(content, item.Name)
			}
			return content, nil
		}
	case *notionapi.URLProperty:
		{
			return p.URL, nil
		}
	case *notionapi.DateProperty:
		{
			content := []time.Time{}
			if p.Date != nil {
				if p.Date.Start != nil {
					content = append(content, time.Time(*p.Date.Start))
				}
				if p.Date.End != nil {
					content = append(content, time.Time(*p.Date.End))
				}
			}
			return content, nil
		}
	case *notionapi.RelationProperty:
		{
			content := []string{}
			for _, item := range p.Relation {
				content = append(content, string(item.ID))
			}
			return content, nil
		}
	}
	return nil, fmt.Errorf("unsupport property: %+v", property)

}

func VarToBlock(variable interface{}, blockType notionapi.BlockType) (notionapi.Block, error) {
	switch blockType {
	case notionapi.BlockTypeParagraph:
		{
			return newParagraph(variable)
		}
	}
	return nil, fmt.Errorf("unsupport variable: %+v", variable)
}
func VarToProperty(variable interface{}, propertyType notionapi.PropertyType) (notionapi.Property, error) {
	if propertyType == "" {
		if _, ok := variable.(string); ok {
			return newRichText(variable)
		}
	} else {
		switch propertyType {
		case notionapi.PropertyTypeTitle:
			{
				return newTitle(variable)
			}
		case notionapi.PropertyTypeRichText:
			{
				return newRichText(variable)
			}
		case notionapi.PropertyTypeMultiSelect:
			{
				return newMultiSelect(variable)
			}
		case notionapi.PropertyTypeURL:
			{
				return newUrl(variable)
			}
		case notionapi.PropertyTypeDate:
			{
				return newDate(variable)
			}
		case notionapi.PropertyTypeRelation:
			{
				return newRelation(variable)
			}
		case notionapi.PropertyTypeNumber:
			{
				return newNumber(variable)
			}
		default:
			{
				log.Printf("unsupport propertyType: %+v", propertyType)
				return nil, fmt.Errorf("unsupport propertyType: %+v", propertyType)
			}
		}
	}
	return nil, fmt.Errorf("unsupport variable: %+v", variable)
}
func newNumber(content interface{}) (*notionapi.NumberProperty, error) {
	if d, ok := content.(float64); ok {
		return &notionapi.NumberProperty{Number: d}, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
}
func genRichTextObj(content []string) []notionapi.RichText {
	result := []notionapi.RichText{}
	for _, item := range content {
		r := []rune(item)
		start := 0
		for start < len(r) {
			if len(r) <= start+1500 {
				result = append(result, notionapi.RichText{Text: &notionapi.Text{Content: string(r[start:])}})
			} else {
				result = append(result, notionapi.RichText{Text: &notionapi.Text{Content: string(r[start : start+1500])}})
			}
			start += 1500
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
	return &result, fmt.Errorf("unsupport content type: %+v", content)
}
func newParagraph(content interface{}) (*notionapi.ParagraphBlock, error) {
	result := notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Type:   notionapi.BlockTypeParagraph,
			Object: notionapi.ObjectTypeBlock,
		},
	}
	raw := []string{}
	if d, ok := content.(string); ok {
		raw = append(raw, d)
	} else if d, ok := content.([]string); ok {
		raw = append(raw, d...)
	}
	result.Paragraph.RichText = append(result.Paragraph.RichText, genRichTextObj(raw)...)
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
