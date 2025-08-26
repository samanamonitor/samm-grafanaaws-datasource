package models

import (
    "encoding/json"
    "net/url"
    "fmt"
    "strings"
    "strconv"
)

type FilterCondition struct {
    Property      string `json:"outProperty"`
    Value         string `json:"value"`
}

type QueryModel struct{
    Service       string `json:"service"`
    ServiceQuery  string `json:"service_query"`
    TextProp      string `json:"text_prop,omitempty"`
    ValueProp     string `json:"value_prop,omitempty"`
    Limit         int    `json:"Limit,omitempty"`
    FieldList     []string `json:"fieldList,omitempty"`
    FilterConditions []FilterCondition `json:"filterConditions,omitempty"`
}

func NewQueryModelFromJSON(jsondata []byte) (QueryModel, error) {
    queryData := QueryModel{Limit: 100}

    err := json.Unmarshal(jsondata, &queryData)
    return queryData, err
}

func NewQueryModelFromUrl(urlstr string) (QueryModel, error) {
    var err error
    req_url, _ := url.Parse(urlstr)
    params     := req_url.Query()

    q := QueryModel{Limit: 100}
    q.Service      = params.Get("service")
    q.ServiceQuery = params.Get("service_query")
    q.TextProp     = params.Get("text_prop")
    q.ValueProp    = params.Get("value_prop")
    q.Limit, err   = strconv.Atoi(params.Get("Limit")) 
    if err != nil {
        q.Limit = -1
    }

    if filterConditions, ok := params["filterConditions"]; ok {
        for _, filterCondition := range filterConditions {
            before, after, found := strings.Cut(filterCondition, ":")
            if found {
                q.FilterConditions = append(q.FilterConditions, FilterCondition{before, after})
            }
        }
    }

    if q.Service == "" {
        return QueryModel{}, fmt.Errorf("Parameter 'service' is mandatory")
    }
    if q.ServiceQuery == "" {
        return QueryModel{}, fmt.Errorf("Parameter 'service_query' is mandatory")
    }
    if q.TextProp == "" {
        return QueryModel{}, fmt.Errorf("Parameter 'text_prop' is mandatory")
    }
    return q, nil
}
