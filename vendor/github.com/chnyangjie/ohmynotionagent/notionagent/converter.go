package notionagent

import (
	"fmt"
	"log"
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

func VarToBlock(variable interface{}, blockType notionapi.BlockType) ([]notionapi.Block, error) {
	switch blockType {
	case notionapi.BlockTypeParagraph:
		{
			return newParagraph(variable)
		}
	case notionapi.BlockQuote:
		{
			return newQuote(variable)
		}
	case notionapi.BlockTypeHeading1, notionapi.BlockTypeHeading2, notionapi.BlockTypeHeading3:
		{
			return newHeading(variable, blockType)
		}
	case notionapi.BlockTypeImage:
		{
			return newImage(variable)
		}
	case notionapi.BlockTypeTableBlock:
		{
			return newTableBlock(variable)
		}
	}
	return nil, fmt.Errorf("unsupport variable: %+v", variable)
}
func BlockToColumn(blocks []notionapi.Block) (notionapi.Block, error) {

	return newColumn(blocks)
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
		case notionapi.PropertyTypeCheckbox:
			{
				return newCheckbox(variable)
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
