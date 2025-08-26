package samm

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/samana-group/sammaws/pkg/models"

	"github.com/aws/aws-sdk-go/service/workspaces"

	"github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

)

type SammWorkspaceBundle struct {
    attributes map[string]interface{}
    defaultFieldList []string
    elements []interface{}
    filter *workspaces.DescribeWorkspaceBundlesInput
    filterConditions []models.FilterCondition
    limit int
    nextToken *string
    svc *workspaces.WorkSpaces
}

func NewSammWorkspaceBundle(svc *workspaces.WorkSpaces, filterConditions []models.FilterCondition, Limit int) SammWorkspaceBundle {
    return SammWorkspaceBundle{
        attributes: map[string]interface{}{
            "BundleId": []*string{},
            "BundleType": []*string{},
            "ComputeType": []string{},
            "CreationTime": []*time.Time{},
            "Description": []*string{},
            "ImageId": []*string{},
            "LastUpdatedTime": []*time.Time{},
            "Name": []*string{},
            "Owner": []*string{},
            "RootStorage": []string{},
            "State": []*string{},
            "UserStorage": []string{},
        },
        defaultFieldList: []string {
            "BundleId",
            "BundleType",
            "CreationTime",
            "Description",
            "ImageId",
            "LastUpdatedTime",
            "Name",
            "Owner",
            "State",
        },
        filterConditions: filterConditions,
        limit: Limit,
        svc: svc,
    }
}

func (samm SammWorkspaceBundle) AppendData(elementIndex int, field *data.Field, name string) {
    object := samm.elements[elementIndex].(*workspaces.WorkspaceBundle)
    switch name {
    case "BundleId":
        field.Append(object.BundleId)
    case "BundleType":
        field.Append(object.BundleType)
    case "ComputeType":
        temp, _ := json.Marshal(object.ComputeType)
        field.Append(string(temp))
    case "CreationTime":
        field.Append(object.CreationTime)
    case "Description":
        field.Append(object.Description)
    case "ImageId":
        field.Append(object.ImageId)
    case "LastUpdatedTime":
        field.Append(object.LastUpdatedTime)
    case "Name":
        field.Append(object.Name)
    case "Owner":
        field.Append(object.Owner)
    case "RootStorage":
        temp, _ := json.Marshal(object.RootStorage)
        field.Append(string(temp))
    case "State":
        field.Append(object.State)
    case "UserStorage":
        temp, _ := json.Marshal(object.UserStorage)
        field.Append(string(temp))
    }
}

func (samm SammWorkspaceBundle) At(index int) interface{} {
    if index >= 0 && index < len(samm.elements) {
        return samm.elements[index]
    }
    return nil
}

func (samm SammWorkspaceBundle) AttributeType(attributeName string) (interface{}, bool) {
    attr, ok := samm.attributes[attributeName]
    return attr, ok
}

func (samm *SammWorkspaceBundle) createFilter(NextToken *string) {
    samm.filter = &workspaces.DescribeWorkspaceBundlesInput{}

    if NextToken != nil {
        samm.filter.SetNextToken(*NextToken)
    }
    bundleIds := []*string{}
    for _, filterCondition := range samm.filterConditions {
        switch property := filterCondition.Property; property {
        case "BundleId":
            bundleIds = append(bundleIds, &filterCondition.Value)
        default:
            log.DefaultLogger.Warn("Invalid property in filter", "property", property)
        }
    }
    if len(bundleIds) > 0 {
        log.DefaultLogger.Debug(fmt.Sprintf("Filter by BundleIds.", "ids_count", len(bundleIds)))
        samm.filter.SetBundleIds(bundleIds)
    }
    log.DefaultLogger.Debug("Input", "filter", samm.filter)
}

func (samm SammWorkspaceBundle) DefaultFieldList() []string {
    return samm.defaultFieldList
}

func (samm SammWorkspaceBundle) Elements() []interface{} {
    return samm.elements
}

func (samm SammWorkspaceBundle) Len() int {
    return len(samm.elements)
}

func (samm SammWorkspaceBundle) NextToken() *string {
    return samm.nextToken
}

func (samm *SammWorkspaceBundle) Query(elements []interface{}) ([]interface{}, *string, error) {
    var err error
    NextToken := samm.filter.NextToken
    for {
        var awsoutput *workspaces.DescribeWorkspaceBundlesOutput

        awsoutput, err = samm.svc.DescribeWorkspaceBundles(samm.filter)
        if (err != nil) {
            log.DefaultLogger.Error("Unable to collect objects.", "error", err.Error())
            return elements, NextToken, err
        }
        for _, e := range awsoutput.Bundles {
            elements = append(elements, e)
        }
        log.DefaultLogger.Debug("workspaces.DescribeWorkspaceBundles Elements.", "cache_length", len(elements))

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

func (samm *SammWorkspaceBundle) UpdateElements(cachedElements interface{}, nextToken *string, cacheIsValid bool) (error) {
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
