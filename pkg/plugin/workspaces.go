package plugin

import (
    "fmt"
    "errors"
    "encoding/json"
    "github.com/aws/aws-sdk-go/service/workspaces"

    "github.com/grafana/grafana-plugin-sdk-go/backend"

    "github.com/samana-group/sammaws/pkg/models"
    "github.com/samana-group/sammaws/pkg/samm"
)

type WorkspacesQuery struct {
    svc        *workspaces.WorkSpaces
    queryData  models.QueryModel
    actionData models.ActionModel
    role       string
    dataSource *Datasource
    refID      string
}

func NewWorkspacesQuery(queryData models.QueryModel, actionData models.ActionModel, dataSource *Datasource, role string, refID string) WorkspacesQuery {
    return WorkspacesQuery{
        svc: workspaces.New(dataSource.AwsSession),
        queryData: queryData,
        actionData: actionData,
        role: role,
        dataSource: dataSource,
        refID: refID,
    }
}

func (w WorkspacesQuery) QueryData() backend.DataResponse {
    switch w.queryData.ServiceQuery {
    case "DescribeWorkspaces":
        return w.workspacesToResponse()
    case "DescribeWorkspacesFields":
        return w.workspaceFieldsToResponse()
    
    case "DescribeWorkspacesConnectionStatus":
        return w.workspacesConnectionStatusToResponse()
    case "DescribeWorkspacesConnectionStatusFields":
        return w.workspacesConnectionStatusFieldsToResponse()
    
    case "DescribeWorkspaceDirectories":
        return w.workspaceDirectoriesToResponse()
    case "DescribeWorkspaceDirectoriesFields":
        return w.workspaceDirectoriesFieldsToResponse()
    
    case "DescribeWorkspaceBundles":
        return w.workspaceBundlesToResponse()
    case "DescribeWorkspaceBundlesFields":
        return w.workspaceBundlesFieldsToResponse()
    
    case "Echo":
        return w.echoToResponse()
    }
    return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Not Implemented service_query %v", w.queryData.ServiceQuery))
}

func (w WorkspacesQuery) QueryVariable() ([]byte, error) {

    return w.toVariables(w.QueryData())

}

func (w WorkspacesQuery) CallAction() ([]byte, error) {
    if w.role != "Admin" {
        return []byte{}, errors.New("Not Authorized")
    }

    switch w.actionData.Action {
    case "start-workspaces":
        return w.startWorkspaces()
    
    case "stop-workspaces":
        return w.stopWorkspaces()
    
    case "reboot-workspaces":
        return w.rebootWorkspaces()
    
    case "restore-workspace":
        return w.restoreWorkspace()
    
    case "echo":
        return w.echo()
    
    case "list-actions":
        return w.listActions()
    }

    return []byte{}, errors.New(fmt.Sprintf("Not Implemented action: %s", w.actionData.Action))
}

func (w WorkspacesQuery) ListActions() ([]byte, error) {
    disabled := w.role != "Admin"
    actions := []SammAwsAction {
        {
            Action: "start-workspaces",
            DisplayName: "Start",
            Disabled: disabled,
            Confirm: true,
        },
        {
            Action: "stop-workspaces",
            DisplayName: "Stop",
            Disabled: disabled,
            Confirm: true,
        },
        {
            Action: "reboot-workspaces",
            DisplayName: "Reboot",
            Disabled: disabled,
            Confirm: true,
        },
        {
            Action: "restore-workspace",
            DisplayName: "Restore",
            Disabled: disabled,
            Confirm: true,
        },
        {
            Action: "echo",
            DisplayName: "Echo",
            Disabled: disabled,
            Confirm: true,
        },
    }
    return json.Marshal(actions)
}


/*    Variables    */
func (w WorkspacesQuery) toVariables(fl backend.DataResponse) ([]byte, error) {
    if fl.Error != nil {
        return []byte{}, fl.Error
    }
    return json.Marshal(fl.Frames)
}

/*    Queries    */
func (w WorkspacesQuery) workspaceFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "BundleId",
        "ComputerName",
        "DataReplicationSettings",
        "DirectoryId",
        "ErrorCode",
        "ErrorMessage",
        "IpAddress",
        "ModificationStates",
        "RelatedWorkspaces",
        "RootVolumeEncryptionEnabled",
        "StandbyWorkspacesProperties",
        "State",
        "SubnetId",
        "UserName",
        "UserVolumeEncryptionEnabled",
        "VolumeEncryptionKey",
        "WorkspaceId",
        "WorkspaceName",
        "WorkspaceProperties",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (w WorkspacesQuery) workspacesToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammWorkspace(w.svc, w.queryData.FilterConditions, w.queryData.Limit)

    /* Process Cache */
    serviceKey := "workspaces.Workspace"
    if len(w.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := w.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, w.queryData.FieldList, w.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */


func (w WorkspacesQuery) workspacesConnectionStatusFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "WorkspaceId",
        "ConnectionState",
        "ConnectionStateCheckTimestamp",
        "LastKnownUserConnectionTimestamp",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (w WorkspacesQuery) workspacesConnectionStatusToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammWorkspacesConnectionStatus(w.svc, w.queryData.FilterConditions, w.queryData.Limit)

    /* Process Cache */
    serviceKey := "workspaces.WorkspaceConnectionStatus"
    if len(w.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := w.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, w.queryData.FieldList, w.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */


func (w WorkspacesQuery) workspaceDirectoriesFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "ActiveDirectoryConfig",
        "Alias",
        "CertificateBasedAuthProperties",
        "CustomerUserName",
        "DirectoryId",
        "DirectoryName",
        "DirectoryType",
        "DnsIpAddresses",
        "ErrorMessage",
        "IamRoleId",
        "IpGroupIds",
        "RegistrationCode",
        "SamlProperties",
        "SelfservicePermissions",
        "State",
        "StreamingProperties",
        "SubnetIds",
        "Tenancy",
        "UserIdentityType",
        "WorkspaceAccessProperties",
        "WorkspaceCreationProperties",
        "WorkspaceDirectoryDescription",
        "WorkspaceDirectoryName",
        "WorkspaceSecurityGroupId",
        "WorkspaceType",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (w WorkspacesQuery) workspaceDirectoriesToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammWorkspacesDirectory(w.svc, w.queryData.FilterConditions, w.queryData.Limit)

    /* Process Cache */
    serviceKey := "workspaces.WorkspacesDirectory"
    if len(w.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := w.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */

    frame, err := CreateFrame(sw, w.queryData.FieldList, w.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */

func (w WorkspacesQuery) workspaceBundlesFieldsToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "BundleId",
        "BundleType",
        "ComputeType",
        "CreationTime",
        "Description",
        "ImageId",
        "LastUpdatedTime",
        "Name",
        "Owner",
        "RootStorage",
        "State",
        "UserStorage",
    }
    return fieldsToResponse(fields, fieldlist)
}

func (w WorkspacesQuery) workspaceBundlesToResponse() backend.DataResponse {
    var response backend.DataResponse
    sw := samm.NewSammWorkspaceBundle(w.svc, w.queryData.FilterConditions, w.queryData.Limit)

    /* Process Cache */
    serviceKey := "workspaces.WorkspaceBundle"
    if len(w.queryData.FilterConditions) > 0 {
        err := sw.UpdateElements([]interface{}{}, nil, false)
        if err != nil {
            response.Error = err
        }
    } else {
        cacheItem := w.dataSource.Cache.Get(serviceKey)
        err := sw.UpdateElements(cacheItem.Objects.([]interface{}), cacheItem.NextToken, cacheItem.IsValid())
        if err != nil {
            response.Error = err
        }
        cacheItem.Update(sw, err)
    }
    /* End Process Cache */


    frame, err := CreateFrame(sw, w.queryData.FieldList, w.refID)
    if err != nil {
        response.Error = err
        return response
    }

    response.Frames = append(response.Frames, frame)

    return response
}
/* ************************************************************* */
func (w WorkspacesQuery) echoToResponse() backend.DataResponse {
    fieldlist := []string{ "Label", "Value" }
    fields := []string{
        "echo1",
        "echo2",
        "echo3",
        "echo4",
        "echo5",
        "echo6",
        "echo7",
        "echo8",
        "echo9",
        "echo0",
    }
    return fieldsToResponse(fields, fieldlist)
}

/*    Actions    */
func (w WorkspacesQuery) startWorkspaces() ([]byte, error) {
    sr := workspaces.StartRequest {
        WorkspaceId: &w.actionData.Id,
    }
    input := workspaces.StartWorkspacesInput{
        StartWorkspaceRequests: []*workspaces.StartRequest{&sr},
    }
    _, err := w.svc.StartWorkspaces(&input)
    if err != nil {
        return []byte{}, err
    }
    return []byte("{ \"message\": \"Workspace is being Started\" }"), nil
}

func (w WorkspacesQuery) stopWorkspaces() ([]byte, error) {
    sr := workspaces.StopRequest {
        WorkspaceId: &w.actionData.Id,
    }
    input := workspaces.StopWorkspacesInput{
        StopWorkspaceRequests: []*workspaces.StopRequest{&sr},
    }
    _, err := w.svc.StopWorkspaces(&input)
    if err != nil {
        return []byte{}, err
    }
    return []byte("{ \"message\": \"Workspace is being Stopped\" }"), nil
}

func (w WorkspacesQuery) rebootWorkspaces() ([]byte, error) {
    sr := workspaces.RebootRequest {
        WorkspaceId: &w.actionData.Id,
    }
    input := workspaces.RebootWorkspacesInput{
        RebootWorkspaceRequests: []*workspaces.RebootRequest{&sr},
    }
    _, err := w.svc.RebootWorkspaces(&input)
    if err != nil {
        return []byte{}, err
    }
    return []byte("{ \"message\": \"Workspace is being Rebooted\" }"), nil
}

func (w WorkspacesQuery) restoreWorkspace() ([]byte, error) {
    return []byte{}, errors.New("Not Implemented restore")
}

func (w WorkspacesQuery) echo() ([]byte, error) {
    return []byte(fmt.Sprintf("{ \"message\": \"You requested an echo from: %s\" }", w.actionData.Id)), nil
}

func (w WorkspacesQuery) listActions() ([]byte, error) {
    workspaceId := w.actionData.Id
    sw := samm.NewSammWorkspace(w.svc, []models.FilterCondition{{"WorkspaceId", workspaceId}} , 1)
    err := sw.UpdateElements([]interface{}{}, nil, false)
    if err != nil || sw.Len() != 1 {
        return []byte{}, fmt.Errorf("Unable to get information for workspaceId=\"%s\".", workspaceId)
    }
    ws := sw.At(0).(*workspaces.Workspace)
    isAdmin := w.role == "Admin"
    actions := []SammAwsAction {
        {
            Action: "start-workspaces",
            DisplayName: "Start",
            Disabled: !(isAdmin && ( 
                *ws.State == "STOPPED" ||
                *ws.State == "SUSPENDED")),
            Confirm: true,
        },
        {
            Action: "stop-workspaces",
            DisplayName: "Stop",
            Disabled: !(isAdmin && (
                *ws.State == "AVAILABLE" ||
                *ws.State == "IMPAIRED" ||
                *ws.State == "UNHEALTHY" ||
                *ws.State == "REBOOTING" ||
                *ws.State == "STARTING" ||
                *ws.State == "SUSPENDED")),
            Confirm: true,
        },
        {
            Action: "reboot-workspaces",
            DisplayName: "Reboot",
            Disabled: !(isAdmin && (
                *ws.State == "AVAILABLE" ||
                *ws.State == "IMPAIRED" ||
                *ws.State == "UNHEALTHY")),
            Confirm: true,
        },
        {
            Action: "restore-workspace",
            DisplayName: "Restore",
            Disabled: !isAdmin,
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
