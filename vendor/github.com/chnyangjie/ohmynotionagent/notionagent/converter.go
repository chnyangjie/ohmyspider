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
		default:
			{
				log.Printf("unsupport propertyType: %+v", propertyType)
				return nil, fmt.Errorf("unsupport propertyType: %+v", propertyType)
			}
		}
	}
	return nil, fmt.Errorf("unsupport variable: %+v", variable)
}
func newRichText(content interface{}) (*notionapi.RichTextProperty, error) {
	if d, ok := content.(string); ok {
		return &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: d}}}}, nil
	}
	return nil, fmt.Errorf("unsupport content type: %+v", content)
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
