package samm

import (
	"strings"
	"fmt"

	"github.com/samana-group/sammaws/pkg/models"

	"github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammWorkspace struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *workspaces.DescribeWorkspacesInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *workspaces.WorkSpaces
}

func NewSammWorkspace(svc *workspaces.WorkSpaces, filterConditions []models.FilterCondition, Limit int) SammWorkspace {
    return SammWorkspace{
        attributes: map[string]interface{} {
            "BundleId":  []*string{},
            "ComputerName": []*string{},
            "DataReplicationSettings":  []string{},
            "DirectoryId": []*string{},
            "ErrorCode": []*string{},
            "ErrorMessage": []*string{},
            "IpAddress": []*string{},
            "ModificationStates":  []string{},
            "RelatedWorkspaces":  []string{},
            "RootVolumeEncryptionEnabled": []*bool{},
            "StandbyWorkspacesProperties":  []string{},
            "State": []*string{},
            "SubnetId": []*string{},
            "UserName": []*string{},
            "UserVolumeEncryptionEnabled": []*bool{},
            "VolumeEncryptionKey": []*string{},
            "WorkspaceId": []*string{},
            "WorkspaceName": []*string{},
            "WorkspaceProperties":  []string{},
        },
        defaultFieldList: []string {
            "WorkspaceId",
            "UserName",
            "ComputerName",
            "DirectoryId",
            "IpAddress",
            "State",
            "BundleId",
            "SubnetId",
            "ErrorCode",
            "ErrorMessage",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammWorkspace) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*workspaces.Workspace)
    switch name {
    case "BundleId":
        field.Append(object.BundleId)
    case "ComputerName":
        field.Append(object.ComputerName)
    case "DataReplicationSettings":
        field.Append(object.DataReplicationSettings.String())
    case "DirectoryId":
        field.Append(object.DirectoryId)
    case "ErrorCode":
        field.Append(object.ErrorCode)
    case "ErrorMessage":
        field.Append(object.ErrorMessage)
    case "IpAddress":
        field.Append(object.IpAddress)
    case "ModificationStates":
        temp := make([]string, len(object.ModificationStates))
        for i, state := range object.ModificationStates {
            temp[i] = state.String()
        }
        field.Append(strings.Join(temp, ","))
    case "RelatedWorkspaces":
        temp := make([]string, len(object.RelatedWorkspaces))
        for i, state := range object.RelatedWorkspaces {
            temp[i] = state.String()
        }
        field.Append(strings.Join(temp, ","))
    case "RootVolumeEncryptionEnabled":
        field.Append(object.RootVolumeEncryptionEnabled)
    case "StandbyWorkspacesProperties":
        temp := make([]string, len(object.StandbyWorkspacesProperties))
        for i, state := range object.StandbyWorkspacesProperties {
            temp[i] = state.String()
        }
        field.Append(strings.Join(temp, ","))
    case "State":
        field.Append(object.State)
    case "SubnetId":
        field.Append(object.SubnetId)
    case "UserName":
        field.Append(object.UserName)
    case "UserVolumeEncryptionEnabled":
        field.Append(object.UserVolumeEncryptionEnabled)
    case "VolumeEncryptionKey":
        field.Append(object.VolumeEncryptionKey)
    case "WorkspaceId":
        field.Append(object.WorkspaceId)
    case "WorkspaceName":
        field.Append(object.WorkspaceName)
    case "WorkspaceProperties":
        field.Append(object.WorkspaceProperties.String())
    }
}

func (samm SammWorkspace) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammWorkspace) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammWorkspace) createFilter(NextToken *string) {
    samm.filter = &workspaces.DescribeWorkspacesInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    workspaceIds := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "BundleId":
            samm.filter.SetBundleId(filterCondition.Value)
        case "DirectoryId":
            samm.filter.SetDirectoryId(filterCondition.Value)
        case "UserName":
            samm.filter.SetUserName(filterCondition.Value)
        case "WorkspaceName":
            samm.filter.SetWorkspaceName(filterCondition.Value)
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

func (samm SammWorkspace) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammWorkspace) Elements() []interface{} {
    return samm.elements
}

func (samm SammWorkspace) Len() int {
    return len(samm.elements)
}

func (samm SammWorkspace) NextToken() *string {
    return samm.nextToken
}

func (samm *SammWorkspace) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *workspaces.DescribeWorkspacesOutput

        awsoutput, err = samm.svc.DescribeWorkspaces(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Workspaces {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("workspaces.DescribeWorkspaces Elements.", "cache_length", len(elements))

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

func (samm *SammWorkspace) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
