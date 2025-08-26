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

type SammFleet struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *appstream.DescribeFleetsInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *appstream.AppStream
}

func NewSammFleet(svc *appstream.AppStream, filterConditions []models.FilterCondition, Limit int) SammFleet {
    return SammFleet{
        attributes: map[string]interface{} {
            "Arn":  []*string{},
            "ComputeCapacityStatus": []string{},
            "CreatedTime": time.Time{},
            "Description": []*string{},
            "DisconnectTimeoutInSeconds": []*int64{},
            "DisplayName": []*string{},
            "DomainJoinInfo": []string{},
            "EnableDefaultInternetAccess": []*bool{},
            "FleetType": []*string{},
            "IamRoleArn": []*string{},
            "IdleDisconnectTimeoutInSeconds": []*int64{},
            "ImageArn": []*string{},
            "ImageName": []*string{},
            "InstanceType": []*string{},
            "MaxConcurrentSessions": []*int64{},
            "MaxSessionsPerInstance": []*int64{},
            "MaxUserDurationInSeconds": []*int64{},
            "Name": []*string{},
            "Platform": []*string{},
            "SessionScriptS3Location": []string{},
            "State": []*string{},
            "StreamView": []*string{},
            "UsbDeviceFilterStrings": []string{},
            "VpcConfig": []string{},
        },
        defaultFieldList: []string {
            "Arn",
            "Description",
            "DisplayName",
            "FleetType",
            "IamRoleArn",
            "ImageArn",
            "ImageName",
            "InstanceType",
            "Name",
            "Platform",
            "State",
            "StreamView",
            "VpcConfig",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammFleet) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*appstream.Fleet)
    switch name {
    case "Arn":
        field.Append(object.Arn)
    case "ComputeCapacityStatus":
        temp, err := json.Marshal(object.ComputeCapacityStatus)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    case "CreatedTime":
        field.Append(object.CreatedTime)
    case "Description":
        field.Append(object.Description)
    case "DisconnectTimeoutInSeconds":
        field.Append(object.DisconnectTimeoutInSeconds)
    case "DisplayName":
        field.Append(object.DisplayName)
    case "DomainJoinInfo":
        temp, err := json.Marshal(object.DomainJoinInfo)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    case "EnableDefaultInternetAccess":
        field.Append(object.EnableDefaultInternetAccess)
    case "FleetType":
        field.Append(object.FleetType)
    case "IamRoleArn":
        field.Append(object.IamRoleArn)
    case "IdleDisconnectTimeoutInSeconds":
        field.Append(object.IdleDisconnectTimeoutInSeconds)
    case "ImageArn":
        field.Append(object.ImageArn)
    case "ImageName":
        field.Append(object.ImageName)
    case "InstanceType":
        field.Append(object.InstanceType)
    case "MaxConcurrentSessions":
        field.Append(object.MaxConcurrentSessions)
    case "MaxSessionsPerInstance":
        field.Append(object.MaxSessionsPerInstance)
    case "MaxUserDurationInSeconds":
        field.Append(object.MaxUserDurationInSeconds)
    case "Name":
        field.Append(object.Name)
    case "Platform":
        field.Append(object.Platform)
    case "SessionScriptS3Location":
        temp, err := json.Marshal(object.SessionScriptS3Location)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    case "State":
        field.Append(object.State)
    case "StreamView":
        field.Append(object.StreamView)
    case "UsbDeviceFilterStrings":
        field.Append(object.UsbDeviceFilterStrings)
    case "VpcConfig":
        temp, err := json.Marshal(object.VpcConfig)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    }
}

func (samm SammFleet) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammFleet) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammFleet) createFilter(NextToken *string) {
    samm.filter = &appstream.DescribeFleetsInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    names := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "FleetName":
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

func (samm SammFleet) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammFleet) Elements() []interface{} {
    return samm.elements
}

func (samm SammFleet) Len() int {
    return len(samm.elements)
}

func (samm SammFleet) NextToken() *string {
    return samm.nextToken
}

func (samm *SammFleet) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *appstream.DescribeFleetsOutput

        awsoutput, err = samm.svc.DescribeFleets(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Fleets {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("appstream.DescribeFleets Elements.", "cache_length", len(elements))

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

func (samm *SammFleet) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
