package plugin

import (
    "fmt"
    "errors"
    "encoding/json"

    "github.com/aws/aws-sdk-go/service/appstream"

    "github.com/grafana/grafana-plugin-sdk-go/backend"
    //"github.com/grafana/grafana-plugin-sdk-go/data"
    //"github.com/grafana/grafana-plugin-sdk-go/backend/log"

    "github.com/samana-group/sammaws/pkg/models"
    "github.com/samana-group/sammaws/pkg/samm"
)

type AppstreamQuery struct {
    svc        *appstream.AppStream
    queryData  models.QueryModel
    actionData models.ActionModel
    role       string
    dataSource *Datasource
    refID      string
}

func NewAppstreamQuery(queryData models.QueryModel, actionData models.ActionModel, dataSource *Datasource, role string, refID string) AppstreamQuery {
    return AppstreamQuery{
        svc: appstream.New(dataSource.AwsSession),
        queryData: queryData,
        actionData: actionData,
        role: role,
        dataSource: dataSource,
        refID: refID,
    }
}

func (a AppstreamQuery) QueryData() backend.DataResponse {
    switch a.queryData.ServiceQuery {
    case "DescribeFleets":
        return a.fleetsToResponse()
    case "DescribeFleetsFields":
        return a.fleetFieldsToResponse()

    case "DescribeStacks":
        return a.stacksToResponse()
    case "DescribeStacksFields":
        return a.stacksFieldsToResponse()

    case "DescribeSessions":
        return a.sessionsToResponse()
    case "DescribeSessionsFields":
        return a.sessionsFieldsToResponse()

    case "DescribeDirectoryConfigs":
        return a.directoryConfigsToResponse()
    case "DescribeDirectoryConfigsFields":
        return a.directoryConfigsFieldsToResponse()

    case "ListAssociatedStacks":
        return a.associatedStacksToResponse()
    case "ListAssociatedStacksFields":
        return a.associatedStacksFieldsToResponse()

    case "ListAssociatedFleets":
        return a.associatedFleetsToResponse()
    case "ListAssociatedFleetsFields":
        return a.associatedFleetsFieldsToResponse()

    }
    return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Not Implemented service_query %v", a.queryData.ServiceQuery))
}

func (a AppstreamQuery) QueryVariable() ([]byte, error) {

    return a.toVariables(a.QueryData())

}

func (a AppstreamQuery) CallAction() ([]byte, error) {
    if a.role != "Admin" {
        return []byte{}, errors.New("Not Authorized")
    }
    switch a.actionData.Action {
    case "expire-session":
        return a.expireSession()
    
    case "echo":
        return a.echo()
    
    case "list-actions":
        return a.listActions()

    }
    return []byte{}, errors.New(fmt.Sprintf("Not Implemented action: %s", a.actionData.Action))
}

func (a AppstreamQuery) ListActions() ([]byte, error) {
    disabled := a.role != "Admin"
    actions := []SammAwsAction{
        {
            Action: "expiresession",
            DisplayName: "Logoff User",
            Disabled: disabled,
            Confirm: true,
        },
        {
            Action: "echo",
            DisplayName: "Echo" + a.role,
            Disabled: disabled,
            Confirm: true,
        },
    }
    return json.Marshal(actions)
}

/*    Variables    */

func (a AppstreamQuery) toVariables(fl backend.DataResponse) ([]byte, error) {
    if fl.Error != nil {
        return []byte{}, fl.Error
    }
    return json.Marshal(fl.Frames)
}

/*    Queries    */
func (w AppstreamQuery) fleetFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "Arn",
        "ComputeCapacityStatus",
        "CreatedTime",
        "Description",
        "DisconnectTimeoutInSeconds",
        "DisplayName",
        "DomainJoinInfo",
        "EnableDefaultInternetAccess",
        "FleetErrors",
        "FleetType",
        "IamRoleArn",
        "IdleDisconnectTimeoutInSeconds",
        "ImageArn",
        "ImageName",
        "InstanceType",
        "MaxConcurrentSessions",
        "MaxSessionsPerInstance",
        "MaxUserDurationInSeconds",
        "Name",
        "Platform",
        "SessionScriptS3Location",
        "State",
        "StreamView",
        "UsbDeviceFilterStrings",
        "VpcConfig",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) fleetsToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammFleet(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    /* Process Cache */
    serviceKey := "appstream.Fleet"
    if len(a.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := a.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w AppstreamQuery) stacksFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
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
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) stacksToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammStack(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    /* Process Cache */
    serviceKey := "appstream.Stack"
    if len(a.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := a.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w AppstreamQuery) sessionsFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
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
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) sessionsToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammSession(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    /* Process Cache */
    serviceKey := "appstream.Session"
    if len(a.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := a.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w AppstreamQuery) directoryConfigsFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "CertificateBasedAuthProperties",
        "CreatedTime",
        "DirectoryName",
        "OrganizationalUnitDistinguishedNames",
        "ServiceAccountCredentials",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) directoryConfigsToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammDirectoryConfigs(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    /* Process Cache */
    serviceKey := "appstream.DirectoryConfig"
    if len(a.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := a.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w AppstreamQuery) associatedStacksFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "Names",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) associatedStacksToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammAssociatedStacks(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    err := sw.UpdateElements([]interface{}{}, nil, false)
    if err != nil {
        response.Error = err
    }

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w AppstreamQuery) associatedFleetsFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "Names",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (a AppstreamQuery) associatedFleetsToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammAssociatedFleets(a.svc, a.queryData.FilterConditions, a.queryData.Limit)

    /* Process Cache */
    err := sw.UpdateElements([]interface{}{}, nil, false)
    if err != nil {
        response.Error = err
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, a.queryData.FieldList, a.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}

func (a AppstreamQuery) expireSession() ([]byte, error) {
    filter := appstream.ExpireSessionInput{ 
        SessionId: &a.actionData.Id,
    }
    _, err := a.svc.ExpireSession(&filter)
    if err != nil {
        return []byte{}, err
    }
    return []byte("{ \"message\": \"Session is being Expired\" }"), nil
}

func (a AppstreamQuery) echo() ([]byte, error) {
    return []byte(fmt.Sprintf("{ \"message\": \"You requested an echo from: %s\" }", a.actionData.Id)), nil
}

func (a AppstreamQuery) listActions() ([]byte, error) {
    isAdmin := a.role == "Admin"
    actions := []SammAwsAction {
        {
            Action: "expire-session",
            DisplayName: "Expire Session",
            Disabled: !(isAdmin),
            Confirm: true,
        },
        /*
        {
            Action: "echo",
            DisplayName: "Echo",
            Disabled: false,
            Confirm: true,
        },
        */
    }
    return json.Marshal(actions)
}
