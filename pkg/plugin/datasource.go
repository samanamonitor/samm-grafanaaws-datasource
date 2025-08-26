package plugin

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "fmt"

    "github.com/grafana/grafana-plugin-sdk-go/backend"
    "github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
    //"github.com/grafana/grafana-plugin-sdk-go/backend/log"

    "github.com/samana-group/sammaws/pkg/models"
    "github.com/samana-group/sammaws/pkg/cache"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/client"
    "github.com/aws/aws-sdk-go/aws/request"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
    _ backend.QueryDataHandler      = (*Datasource)(nil)
    _ backend.CheckHealthHandler    = (*Datasource)(nil)
    _ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func NewSammAwsResponse(message string, status int, sender backend.CallResourceResponseSender) error {
    body, _ := json.Marshal(SammAwsResponse{ Message: message })
    return sender.Send(&backend.CallResourceResponse{
        Status: status,
        Body: body,
    })
}

type SammAwsAction struct {
    Action      string `json:"action"`
    DisplayName string `json:"displayname"`
    Disabled    bool   `json:"disabled"`
    Confirm     bool   `json:"confirm"`
}


type Datasource struct{
    AwsSession *session.Session
    Cache cache.CacheMap
    CacheDuration time.Duration
}

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
    var err error
    config, err := models.LoadPluginSettings(settings)
    if err != nil {
        return nil, err
    }

    retryer := client.DefaultRetryer {
        NumMaxRetries: config.NumMaxRetries,
        MinRetryDelay: time.Duration(config.MinRetryDelay) * time.Millisecond,
        MinThrottleDelay: time.Duration(config.MinThrottleDelay) * time.Millisecond,
        MaxRetryDelay: time.Duration(config.MaxRetryDelay) * time.Second,
        MaxThrottleDelay: time.Duration(config.MaxThrottleDelay) * time.Second,
    }
    awsconfig := aws.NewConfig().WithRegion(config.Region)
    if config.AccessKey != "" {
        awsconfig.WithCredentials(credentials.NewStaticCredentials(
            config.AccessKey, 
            config.Secrets.AccessSecret, 
            config.Secrets.AccessToken))
    }
    awsconfig = request.WithRetryer(awsconfig, retryer)
    sess, err := session.NewSession(awsconfig)
    if err != nil {
        return nil, err
    }
    d := Datasource{
        AwsSession: sess,
        Cache: cache.NewCacheMap(time.Duration(config.CacheSeconds) * time.Second),
    }
    return &d, nil
}

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
    var role string

    if req.PluginContext.User != nil {
        role = req.PluginContext.User.Role
    } else {
        role = ""
    }

    if req.Path == "query" {
        return d.queryVariable(role, ctx, req, sender)

    } else if req.Path == "list-actions" {
        return d.listActions(role, ctx, req, sender)

    } else if req.Path == "action" {
        return d.callAction(role, ctx, req, sender)

    } else {
        return NewSammAwsResponse("Resource not Found", http.StatusNotFound, sender)
    }
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
    // Clean up datasource instance resources.
}


// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
    // create response struct
    response := backend.NewQueryDataResponse()
    var query SammAwsQuery

    for _, q := range req.Queries {

        queryData, err := models.NewQueryModelFromJSON(q.JSON)
        if err != nil {
            response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, 
                err.Error())
            continue
        }

        if (queryData.Service == "workspaces") {
            query = NewWorkspacesQuery(queryData, models.ActionModel{}, d, "", q.RefID)

        } else if queryData.Service == "appstream" {
            query = NewAppstreamQuery(queryData, models.ActionModel{}, d, "", q.RefID)

        } else {
            query = NotImplemented{}
        }

        // save the response in a hashmap
        // based on with RefID as identifier
        response.Responses[q.RefID] = query.QueryData()
    }

    return response, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
    var err error
    settings := *req.PluginContext.DataSourceInstanceSettings
    config, err := models.LoadPluginSettings(settings)
    if err != nil {
        return nil, err
    }

    awsconfig := aws.Config{
        Region: &config.Region,
    }
    if config.AccessKey != "" {
        awsconfig.WithCredentials(credentials.NewStaticCredentials(
            config.AccessKey, 
            config.Secrets.AccessSecret, 
            config.Secrets.AccessToken))
    }

    w := NewWorkspacesQuery(models.QueryModel{ 
            ServiceQuery: "DescribeWorkspaces",
            Limit: 1,
        }, models.ActionModel{}, d, "", "test")
    w.QueryData()

    /*
    res := &backend.CheckHealthResult{}
    config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

    if err != nil {
        res.Status = backend.HealthStatusError
        res.Message = "Unable to load settings"
        return res, nil
    }

    
    if config.Secrets.AccessSecret == "" {
        res.Status = backend.HealthStatusError
        res.Message = "Access Secret is missing"
        return res, nil
    }
    */
    return &backend.CheckHealthResult{
        Status:  backend.HealthStatusOk,
        Message: "Data source is working",
    }, nil
}

func (d *Datasource) callAction(role string, ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {

    var query SammAwsQuery
    var actionData models.ActionModel

    err := json.Unmarshal(req.Body, &actionData)
    if err != nil {
        return NewSammAwsResponse(err.Error(), http.StatusBadRequest, sender)
    }

    if actionData.Service == "workspaces" {
        query = NewWorkspacesQuery(models.QueryModel{}, actionData, d, role, "action")

    } else if actionData.Service == "appstream" {
        query = NewAppstreamQuery(models.QueryModel{}, actionData, d, role, "action")

    } else {
        query = NotImplemented{}
    }

    body, err := query.CallAction()
    if err != nil {
        return NewSammAwsResponse(err.Error(), http.StatusBadRequest, sender)
    }
    return sender.Send(&backend.CallResourceResponse{
        Status: http.StatusOK,
        Body: body,
        })
}

func (d *Datasource) listActions(role string, ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
    var query SammAwsQuery

    queryData, err := models.NewQueryModelFromJSON(req.Body)
    if err != nil {
        return NewSammAwsResponse(
            err.Error(), 
            http.StatusBadRequest, 
            sender)
    }

    if queryData.Service == "workspaces" {
        query = NewWorkspacesQuery(queryData, models.ActionModel{}, d, role, "listactions")

    } else if queryData.Service == "appstream" {
        query = NewAppstreamQuery(queryData, models.ActionModel{}, d, role, "listactions")

    } else {
        query = NotImplemented{}
    }

    body, err := query.ListActions()
    if err != nil {
        return NewSammAwsResponse(err.Error(), http.StatusBadRequest, sender)
    }
    return sender.Send(&backend.CallResourceResponse{
        Status: http.StatusOK,
        Body: body,
        })
}

func (d *Datasource) queryVariable(role string, ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
    var query SammAwsQuery

    queryData, err := models.NewQueryModelFromJSON(req.Body)
    if err != nil {
        return NewSammAwsResponse(
            err.Error(), 
            http.StatusBadRequest, 
            sender)
    }

    if queryData.Service == "workspaces" {
        query = NewWorkspacesQuery(queryData, models.ActionModel{}, d, role, "variable")

    } else if queryData.Service == "appstream" {
        query = NewAppstreamQuery(queryData, models.ActionModel{}, d, role, "variable")
    } else {
        return NewSammAwsResponse(
            fmt.Sprintf("Unable to process query. query=%+v", queryData), 
            http.StatusBadRequest, 
            sender)
    }

    var body []byte
    body, err = query.QueryVariable()
    if err != nil {
        return NewSammAwsResponse(err.Error(), http.StatusBadRequest, sender)
    }
    return sender.Send(&backend.CallResourceResponse{
        Status: http.StatusOK,
        Body: body,
        })
}
