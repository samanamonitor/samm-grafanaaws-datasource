package samm

import (
	"time"
	"fmt"
	"strconv"

	"github.com/samana-group/sammaws/pkg/models"

	"github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammWorkspacesConnectionStatus struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *workspaces.DescribeWorkspacesConnectionStatusInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *workspaces.WorkSpaces
}

func NewSammWorkspacesConnectionStatus(svc *workspaces.WorkSpaces, filterConditions []models.FilterCondition, Limit int) SammWorkspacesConnectionStatus {
    return SammWorkspacesConnectionStatus{
        attributes: map[string]interface{} {
            "ConnectionState": []*string{},
            "ConnectionStateCheckTimestamp": []*time.Time{},
            "LastKnownUserConnectionTimestamp": []*time.Time{},
            "WorkspaceId": []*string{},
        },
        defaultFieldList: []string {
            "ConnectionState",
            "ConnectionStateCheckTimestamp",
            "LastKnownUserConnectionTimestamp",
            "WorkspaceId",
        },
		filterConditions: filterConditions,
		limit: Limit,
		svc: svc,
    }
}

func (samm SammWorkspacesConnectionStatus) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*workspaces.WorkspaceConnectionStatus)
    switch name {
    case "ConnectionState":
        field.Append(object.ConnectionState)
    case "ConnectionStateCheckTimestamp":
        field.Append(object.ConnectionStateCheckTimestamp)
    case "LastKnownUserConnectionTimestamp":
        field.Append(object.LastKnownUserConnectionTimestamp)
    case "WorkspaceId":
        field.Append(object.WorkspaceId)
    }
}

func (samm SammWorkspacesConnectionStatus) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammWorkspacesConnectionStatus) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammWorkspacesConnectionStatus) createFilter(NextToken *string) {
    samm.filter = &workspaces.DescribeWorkspacesConnectionStatusInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
	workspaceIds := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "WorkspaceId":
            value := filterCondition.Value
            workspaceIds = append(workspaceIds, &value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
		}
    }
    if len(workspaceIds) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by workspaceId.", "ids_count", len(workspaceIds)))
        samm.filter.SetWorkspaceIds(workspaceIds)
    } 
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammWorkspacesConnectionStatus) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammWorkspacesConnectionStatus) Elements() []interface{} {
    return samm.elements
}

func (samm SammWorkspacesConnectionStatus) Len() int {
    return len(samm.elements)
}

func (samm SammWorkspacesConnectionStatus) NextToken() *string {
    return samm.nextToken
}

func (samm *SammWorkspacesConnectionStatus) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *workspaces.DescribeWorkspacesConnectionStatusOutput

        awsoutput, err = samm.svc.DescribeWorkspacesConnectionStatus(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.WorkspacesConnectionStatus {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("workspaces.DescribeWorkspacesConnectionStatus Elements.", "cache_length", strconv.Itoa(len(elements)))

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

func (samm *SammWorkspacesConnectionStatus) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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

