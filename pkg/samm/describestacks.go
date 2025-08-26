package samm

import (
    "fmt"
    "encoding/json"
    "time"

    "github.com/samana-group/sammaws/pkg/models"

    "github.com/aws/aws-sdk-go/service/appstream"

    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammStack struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *appstream.DescribeStacksInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *appstream.AppStream
}

func NewSammStack(svc *appstream.AppStream, filterConditions []models.FilterCondition, Limit int) SammStack {
    return SammStack{
        attributes: map[string]interface{} {
			"AccessEndpoints": []string{},
			"ApplicationSettings": []string{},
			"Arn": []*string{},
			"CreatedTime": []*time.Time{},
			"Description": []*string{},
			"DisplayName": []*string{},
			"EmbedHostDomains": []string{},
			"FeedbackURL": []*string{},
			"Name": []*string{},
			"RedirectURL": []*string{},
			"StackErrors": []string{},
			"StorageConnectors": []string{},
			"StreamingExperienceSettings": []string{},
			"UserSettings": []string{},
        },
        defaultFieldList: []string {
			"AccessEndpoints",
			"ApplicationSettings",
			"Arn",
			"CreatedTime",
			"Description",
			"DisplayName",
			"EmbedHostDomains",
			"FeedbackURL",
			"Name",
			"RedirectURL",
			"StackErrors",
			"StorageConnectors",
			"StreamingExperienceSettings",
			"UserSettings",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammStack) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*appstream.Stack)
    switch name {
	case "AccessEndpoints":
		temp, err := json.Marshal(object.AccessEndpoints)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "ApplicationSettings":
		temp, err := json.Marshal(object.ApplicationSettings)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "Arn":
		field.Append(object.Arn)
	case "CreatedTime":
		field.Append(object.CreatedTime)
	case "Description":
		field.Append(object.Description)
	case "DisplayName":
		field.Append(object.DisplayName)
	case "EmbedHostDomains":
		temp, err := json.Marshal(object.EmbedHostDomains)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "FeedbackURL":
		field.Append(object.FeedbackURL)
	case "Name":
		field.Append(object.Name)
	case "RedirectURL":
		field.Append(object.RedirectURL)
	case "StackErrors":
		temp, err := json.Marshal(object.StackErrors)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "StorageConnectors":
		temp, err := json.Marshal(object.StorageConnectors)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "StreamingExperienceSettings":
		temp, err := json.Marshal(object.StreamingExperienceSettings)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "UserSettings":
		temp, err := json.Marshal(object.UserSettings)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    }
}

func (samm SammStack) At(index int) interface{} {
    if index > 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammStack) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammStack) createFilter(NextToken *string) {
    samm.filter = &appstream.DescribeStacksInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    names := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "StackName":
            value := filterCondition.Value
            names = append(names, &value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    if len(names) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by names.", "names_count", len(names)))
        samm.filter.SetNames(names)
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammStack) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammStack) Elements() []interface{} {
    return samm.elements
}

func (samm SammStack) Len() int {
    return len(samm.elements)
}

func (samm SammStack) NextToken() *string {
    return samm.nextToken
}

func (samm *SammStack) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *appstream.DescribeStacksOutput

        awsoutput, err = samm.svc.DescribeStacks(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Stacks {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("appstream.DescribeStacks Elements.", "cache_length", len(elements))

        NextToken = awsoutput.NextToken
        if NextToken == nil {
            return elements, nil, nil
        } else {
            samm.filter.SetNextToken(*NextToken)
        }
        if samm.limit > 0 && len(elements) >= samm.limit {
            log.DefaultLogger.Info("Limit Reached")
            return elements, NextToken, nil
        }
    }
    return elements, nil, err
}

func (samm *SammStack) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
    var err error

    elements := cachedElements.([]interface{})
    NextToken := nextToken

    if cacheIsValid {
        samm.elements = elements
        return nil
    }

    /* Process Filters */
    samm.createFilter(NextToken)
    /* End Process Filters */

    /* Collect Data */
    elements, samm.nextToken, err = samm.Query(elements)
    /* End Collect Data */

    samm.elements = elements
    log.DefaultLogger.Debug("UpdateElements", "len(elements)", len(elements), "NextToken", NextToken)
    return err
}
