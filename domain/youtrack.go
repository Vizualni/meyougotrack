package domain

import (
	"encoding/xml"

	"time"
)

type YouTrackIssue struct {
	Issue  xml.Name `xml:"issue"`
	Fields []Field  `xml:"field"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

func (i *YouTrackIssue) FindField(fieldName string) string {
	for _, field := range i.Fields {
		if field.Name == fieldName {
			return field.Value
		}
	}

	return ""
}

type IssueWorkLog struct {
	IssueId     string
	Description string
	Type        string
	Duration    int
	Date        time.Time
}
