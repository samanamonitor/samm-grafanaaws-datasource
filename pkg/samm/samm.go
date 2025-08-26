package samm

import (
    "github.com/grafana/grafana-plugin-sdk-go/data"
)

type SammElement interface {
    AppendData(elementIndex int, field *data.Field, name string)
    Len() int
    AttributeType(string) (interface{}, bool)
    DefaultFieldList() []string
    Elements() []interface{}
    At(int) interface{}
    NextToken() *string
}
