package samm

import (
    "github.com/samana-group/sammaws/pkg/models"

    "github.com/aws/aws-sdk-go/service/appstream"

    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammAssociatedStacks struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter appstream.ListAssociatedStacksInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *appstream.AppStream
}

func NewSammAssociatedStacks(svc *appstream.AppStream, filterConditions []models.FilterCondition, Limit int) SammAssociatedStacks {
    return SammAssociatedStacks{
        attributes: map[string]interface{} {
			"Names": []*string{},
        },
        defaultFieldList: []string {
			"Names",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammAssociatedStacks) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*string)
	field.Append(object)
}

func (samm SammAssociatedStacks) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammAssociatedStacks) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammAssociatedStacks) createFilter(NextToken *string) {
    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "FleetName":
			samm.filter.SetFleetName(filterCondition.Value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammAssociatedStacks) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammAssociatedStacks) Elements() []interface{} {
    return samm.elements
}

func (samm SammAssociatedStacks) Len() int {
    return len(samm.elements)
}

func (samm SammAssociatedStacks) NextToken() *string {
    return samm.nextToken
}

func (samm *SammAssociatedStacks) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *appstream.ListAssociatedStacksOutput

        awsoutput, err = samm.svc.ListAssociatedStacks(&samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Names {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("appstream.ListAssociatedStacks Elements.", "cache_length", len(elements))

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

func (samm *SammAssociatedStacks) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
