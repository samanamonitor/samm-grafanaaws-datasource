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

type SammDirectoryConfigs struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter appstream.DescribeDirectoryConfigsInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *appstream.AppStream
}

func NewSammDirectoryConfigs(svc *appstream.AppStream, filterConditions []models.FilterCondition, Limit int) SammDirectoryConfigs {
    return SammDirectoryConfigs{
        attributes: map[string]interface{} {
			"CertificateBasedAuthProperties": []string{},
			"CreatedTime": []*time.Time{},
			"DirectoryName": []*string{},
			"OrganizationalUnitDistinguishedNames": []string{},
			"ServiceAccountCredentials": []string{},
        },
        defaultFieldList: []string {
			"CertificateBasedAuthProperties",
			"CreatedTime",
			"DirectoryName",
			"OrganizationalUnitDistinguishedNames",
			"ServiceAccountCredentials",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammDirectoryConfigs) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*appstream.DirectoryConfig)
    switch name {
	case "CertificateBasedAuthProperties":
        temp, err := json.Marshal(object.CertificateBasedAuthProperties)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "CreatedTime":
		field.Append(object.CreatedTime)
	case "DirectoryName":
		field.Append(object.DirectoryName)
	case "OrganizationalUnitDistinguishedNames":
        temp, err := json.Marshal(object.OrganizationalUnitDistinguishedNames)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
	case "ServiceAccountCredentials":
        temp, err := json.Marshal(object.ServiceAccountCredentials)
        if err != nil {
            log.DefaultLogger.Warn("Unable to convert to json.", "error", err.Error(), "attribute", name)
            return
        }
        field.Append(string(temp))
    }
}

func (samm SammDirectoryConfigs) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammDirectoryConfigs) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammDirectoryConfigs) createFilter(NextToken *string) {

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    directorynames := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "Name":
            value := filterCondition.Value
            directorynames = append(directorynames, &value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    if len(directorynames) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by directorynames.", "ndirectorynames_count", len(directorynames)))
        samm.filter.SetDirectoryNames(directorynames)
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammDirectoryConfigs) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammDirectoryConfigs) Elements() []interface{} {
    return samm.elements
}

func (samm SammDirectoryConfigs) Len() int {
    return len(samm.elements)
}

func (samm SammDirectoryConfigs) NextToken() *string {
    return samm.nextToken
}

func (samm *SammDirectoryConfigs) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *appstream.DescribeDirectoryConfigsOutput

        awsoutput, err = samm.svc.DescribeDirectoryConfigs(&samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.DirectoryConfigs {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("appstream.DescribeDirectoryConfigs Elements.", "cache_length", len(elements))

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

func (samm *SammDirectoryConfigs) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
