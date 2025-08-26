package samm

import (
	"strings"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/samana-group/sammaws/pkg/models"

	"github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammWorkspaceDirectory struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *workspaces.DescribeWorkspaceDirectoriesInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *workspaces.WorkSpaces
}

func NewSammWorkspacesDirectory(svc *workspaces.WorkSpaces, filterConditions []models.FilterCondition, Limit int) SammWorkspaceDirectory {
    return SammWorkspaceDirectory{
        attributes: map[string]interface{} {
            "ActiveDirectoryConfig": []string{},
            "Alias": []*string{},
            "CertificateBasedAuthProperties": []string{},
            "CustomerUserName": []*string{},
            "DirectoryId": []*string{},
            "DirectoryName": []*string{},
            "DirectoryType": []*string{},
            "DnsIpAddresses": []string{},
            "ErrorMessage": []*string{},
            "IamRoleId": []*string{},
            "IpGroupIds": []string{},
            "RegistrationCode": []*string{},
            "SamlProperties": []string{},
            "SelfservicePermissions": []string{},
            "State": []*string{},
            "StreamingProperties": []string{},
            "SubnetIds": []string{},
            "Tenancy": []*string{},
            "UserIdentityType": []*string{},
            "WorkspaceAccessProperties": []string{},
            "WorkspaceCreationProperties": []string{},
            "WorkspaceDirectoryDescription": []*string{},
            "WorkspaceDirectoryName": []*string{},
            "WorkspaceSecurityGroupId": []*string{},
            "WorkspaceType": []*string{},
        },
        defaultFieldList: []string {
            "WorkspaceCreationProperties",
            "DirectoryId",
            "DirectoryName",
            "Alias",
            "CustomerUserName",
            "DirectoryType",
            "DnsIpAddresses",
            "RegistrationCode",
            "State",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammWorkspaceDirectory) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*workspaces.WorkspaceDirectory)
    switch name {
    case "ActiveDirectoryConfig":
        if object.ActiveDirectoryConfig != nil {
            field.Append(object.ActiveDirectoryConfig.String())
        } else {
            field.Append("")
        }
    case "Alias":
        field.Append(object.Alias)
    case "CertificateBasedAuthProperties":
        field.Append(object.CertificateBasedAuthProperties.String())
    case "CustomerUserName":
        field.Append(object.CustomerUserName)
    case "DirectoryId":
        field.Append(object.DirectoryId)
    case "DirectoryName":
        field.Append(object.DirectoryName)
    case "DirectoryType":
        field.Append(object.DirectoryType)
    case "DnsIpAddresses":
        temp := make([]string, len(object.DnsIpAddresses))
        for i, ip := range object.DnsIpAddresses {
            temp[i] = *ip
        }
        field.Append(strings.Join(temp, ","))
    case "ErrorMessage":
        field.Append(object.ErrorMessage)
    case "IamRoleId":
        field.Append(object.IamRoleId)
    case "IpGroupIds":
        temp := make([]string, len(object.IpGroupIds))
        for i, id := range object.IpGroupIds {
            temp[i] = *id
        }
        field.Append(strings.Join(temp, ","))
    case "RegistrationCode":
        field.Append(object.RegistrationCode)
    case "SamlProperties":
        temp, _ := json.Marshal(object.SamlProperties)
        field.Append(string(temp))
    case "SelfservicePermissions":
        temp, _ := json.Marshal(object.SelfservicePermissions)
        field.Append(string(temp))
    case "State":
        field.Append(object.State)
    case "StreamingProperties":
        temp, _ := json.Marshal(object.StreamingProperties)
        field.Append(string(temp))
    case "SubnetIds":
        temp := make([]string, len(object.SubnetIds))
        for i, id := range object.SubnetIds {
            temp[i] = *id
        }
        field.Append(strings.Join(temp, ","))
    case "Tenancy":
        field.Append(object.Tenancy)
    case "UserIdentityType":
        field.Append(object.UserIdentityType)
    case "WorkspaceAccessProperties":
        temp, _ := json.Marshal(object.WorkspaceAccessProperties)
        field.Append(string(temp))
    case "WorkspaceCreationProperties":
        temp, _ := json.Marshal(object.WorkspaceCreationProperties)
        field.Append(string(temp))
    case "WorkspaceDirectoryDescription":
        field.Append(object.WorkspaceDirectoryDescription)
    case "WorkspaceDirectoryName":
        field.Append(object.WorkspaceDirectoryName)
    case "WorkspaceSecurityGroupId":
        field.Append(object.WorkspaceSecurityGroupId)
    case "WorkspaceType":
        field.Append(object.WorkspaceType)
    }
}

func (samm SammWorkspaceDirectory) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammWorkspaceDirectory) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammWorkspaceDirectory) createFilter(NextToken *string) {
    samm.filter = &workspaces.DescribeWorkspaceDirectoriesInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    directoryIds := []*string{}
    directoryNames := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "DirectoryId":
                directoryIds = append(directoryIds, &filterCondition.Value)
        case "DirectoryName":
                directoryNames = append(directoryNames, &filterCondition.Value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    if len(directoryIds) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by DirectoryIds.", "ids_count", len(directoryIds)))
        samm.filter.SetDirectoryIds(directoryIds)
    }
    if len(directoryIds) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by WorkspaceDirectoryNames.", "names_count", len(directoryNames)))
        samm.filter.SetWorkspaceDirectoryNames(directoryNames)
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammWorkspaceDirectory) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammWorkspaceDirectory) Elements() []interface{} {
    return samm.elements
}

func (samm SammWorkspaceDirectory) Len() int {
    return len(samm.elements)
}

func (samm SammWorkspaceDirectory) NextToken() *string {
    return samm.nextToken
}

func (samm *SammWorkspaceDirectory) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *workspaces.DescribeWorkspaceDirectoriesOutput

        awsoutput, err = samm.svc.DescribeWorkspaceDirectories(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Directories {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("workspaces.DescribeWorkspaceDirectories Elements.", "cache_length", strconv.Itoa(len(elements)))

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

func (samm *SammWorkspaceDirectory) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
