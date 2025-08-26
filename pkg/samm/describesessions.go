package samm

import (
    "encoding/json"
    "time"
    "strconv"

    "github.com/samana-group/sammaws/pkg/models"

    "github.com/aws/aws-sdk-go/service/appstream"

    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammSession struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *appstream.DescribeSessionsInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *appstream.AppStream
}

func NewSammSession(svc *appstream.AppStream, filterConditions []models.FilterCondition, Limit int) SammSession {
    return SammSession{
        attributes: map[string]interface{} {
            "AuthenticationType": []*string{},
            "ConnectionState": []*string{},
            "FleetName": []*string{},
            "Id": []*string{},
            "InstanceId": []*string{},
            "MaxExpirationTime": []*time.Time{},
            "NetworkAccessConfiguration": []string{},
            "StackName": []*string{},
            "StartTime": []*time.Time{},
            "State": []*string{},
            "UserId": []*string{},
        },
        defaultFieldList: []string {
            "AuthenticationType",
            "ConnectionState",
            "FleetName",
            "Id",
            "InstanceId",
            "MaxExpirationTime",
            "NetworkAccessConfiguration",
            "StackName",
            "StartTime",
            "State",
            "UserId",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammSession) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*appstream.Session)
    switch name {
    case "AuthenticationType":
        field.Append(object.AuthenticationType)
    case "ConnectionState":
        field.Append(object.ConnectionState)
    case "FleetName":
        field.Append(object.FleetName)
    case "Id":
        field.Append(object.Id)
    case "InstanceId":
        field.Append(object.InstanceId)
    case "MaxExpirationTime":
        field.Append(object.MaxExpirationTime)
    case "NetworkAccessConfiguration":
        temp, err := json.Marshal(object.NetworkAccessConfiguration)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    case "StackName":
        field.Append(object.StackName)
    case "StartTime":
        field.Append(object.StartTime)
    case "State":
        field.Append(object.State)
    case "UserId":
        field.Append(object.UserId)
    }
}

func (samm SammSession) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammSession) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammSession) createFilter(NextToken *string) {
    samm.filter = &appstream.DescribeSessionsInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "AuthenticationType":
            samm.filter.SetAuthenticationType(filterCondition.Value)
        case "FleetName":
            samm.filter.SetFleetName(filterCondition.Value)
        case "InstanceId":
            samm.filter.SetInstanceId(filterCondition.Value)
        case "Limit":
            temp, err := strconv.Atoi(filterCondition.Value)
            if err != nil {
                log.DefaultLogger.Warn("Unable to convert to int.", "error", err.Error(), "property", property)
                return
            }
            samm.filter.SetLimit(int64(temp))
        case "StackName":
            samm.filter.SetStackName(filterCondition.Value)
        case "UserId":
            samm.filter.SetUserId(filterCondition.Value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammSession) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammSession) Elements() []interface{} {
    return samm.elements
}

func (samm SammSession) Len() int {
    return len(samm.elements)
}

func (samm SammSession) NextToken() *string {
    return samm.nextToken
}

func (samm *SammSession) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *appstream.DescribeSessionsOutput

        awsoutput, err = samm.svc.DescribeSessions(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Sessions {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("appstream.DescribeSessions Elements.", "cache_length", len(elements))

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

func (samm *SammSession) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
