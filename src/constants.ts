import type {SelectableValue } from '@grafana/data';
import type { 
  SammAwsQuery, 
  SammAwsService,
  SammAwsWorkspacesServiceQuery,
  SammAwsAppstreamServiceQuery,
  SammAwsWSDescribeFilter,
  SammAwsWSConnectionFilter,
  SammAwsWSDescribeDirectoryFilter,
  SammAwsASDescribeSessionFilter,
  SammAwsNoneProps,
  SammAwsDescWsProps,
  SammAwsConnWsProps,
  SammAwsDirWsProps,
  SammAwsDescAsFleetProps,
  SammAwsDescAsStackProps,
  SammAwsListAssociatedProps,
  SammAwsDescAsSessionProps,
  FilterQueryDefinition,
  SammAwsFilterCascader,
} from './types';

export const EnableLimit = false;
export const EnableFilters = true;

export const DefaultSammAwsQuery: SammAwsQuery = {
  refId: '',
  service: null,
  service_query: null,
  filter_field: 'none',
  Limit: -1,
  filterConditions: [],
  fieldList: [],
};

export const SERVICE_TYPES: Array<SelectableValue<SammAwsService>> = [
  { label: 'Workspaces', value: 'workspaces' },
  { label: 'Appstream', value: 'appstream' },
  { label: 'Ec2', value: 'ec2' },
];

export const WS_SERVICE_QUERY_TYPES: Array<SelectableValue<SammAwsWorkspacesServiceQuery>> | undefined = [
  { label: 'Describe', value: 'DescribeWorkspaces' },
  { label: 'Connections', value: 'DescribeWorkspacesConnectionStatus' },
  { label: 'Directories', value: 'DescribeWorkspaceDirectories' },
  { label: 'Bundles', value: 'DescribeWorkspaceBundles' },
];

export const AS_SERVICE_QUERY_TYPES: Array<SelectableValue<SammAwsAppstreamServiceQuery>> | undefined = [
  { label: 'Stacks', value: 'DescribeStacks' },
  { label: 'Fleets', value: 'DescribeFleets' },
  { label: 'Sessions', value: 'DescribeSessions' },
  { label: 'Directory Configs', value: 'DescribeDirectoryConfigs' },
  { label: 'Associated Stacks', value: 'ListAssociatedStacks' },
  { label: 'Associated Fleets', value: 'ListAssociatedFleets' },
];

export const WS_DESCRIBE_FILTERS: Array<SelectableValue<SammAwsWSDescribeFilter>> = [
  { label: 'Bundle', value: 'BundleId' },
  { label: 'Directory', value: 'DirectoryId' },
  { label: 'User', value: 'UserName' },
  { label: 'Workspace', value: 'WorkspaceId' },
];

export const WS_CONNECTION_FILTERS: Array<SelectableValue<SammAwsWSConnectionFilter>> = [
  { label: 'Workspace Id', value: 'WorkspaceId' },
];

export const WS_DIRECTORY_FILTERS: Array<SelectableValue<SammAwsWSDescribeDirectoryFilter>> = [
  { label: 'Directory Name', value: 'DirectoryName' },
  { label: 'Directory Id', value: 'DirectoryId' },
];

export const SERVICE_QUERY_TYPES = [
  { service: 'workspaces', service_queries: WS_SERVICE_QUERY_TYPES},
  { service: 'appstream', service_queries: AS_SERVICE_QUERY_TYPES},
];

export const AS_SESSION_FILTERS: Array<SelectableValue<SammAwsASDescribeSessionFilter>> = [
  { label: 'Fleet Name', value: 'FleetName' },
  { label: 'Stack Name', value: 'StackName' },
];

export const AS_DESC_FLEET_PROPS: Array<SelectableValue<SammAwsDescAsFleetProps>> = [
  { label: 'Arn', value: 'Arn' },
  { label: 'Name', value: 'Name' },
  { label: 'Description', value: 'Description' },
  { label: 'DisplayName', value: 'DisplayName' },
  { label: 'FleetType', value: 'FleetType' },
  { label: 'IamRoleArn', value: 'IamRoleArn' },
  { label: 'ImageArn', value: 'ImageArn' },
  { label: 'ImageName', value: 'ImageName' },
  { label: 'InstanceType', value: 'InstanceType' },
  { label: 'Platform', value: 'Platform' },
  { label: 'State', value: 'State' },
  { label: 'StreamView', value: 'StreamView' },
];

export const AS_DESC_STACK_PROPS: Array<SelectableValue<SammAwsDescAsStackProps>> = [
  { label: 'Arn', value: 'Arn' },
  { label: 'CreatedTime', value: 'CreatedTime' },
  { label: 'Description', value: 'Description' },
  { label: 'DisplayName', value: 'DisplayName' },
  { label: 'FeedbackURL', value: 'FeedbackURL' },
  { label: 'Name', value: 'Name' },
  { label: 'RedirectURL', value: 'RedirectURL' },
];

export const AS_LIST_ASSOCIATEDSTACKS_PROPS: Array<SelectableValue<SammAwsListAssociatedProps>> = [
  { label: 'Name', value: 'Names' },
];

export const AS_LIST_ASSOCIATEDFLEETS_PROPS: Array<SelectableValue<SammAwsListAssociatedProps>> = [
  { label: 'Name', value: 'Names' },
];

export const AS_DESC_SESSION_PROPS: Array<SelectableValue<SammAwsDescAsSessionProps>> = [
  { label: 'Session Id', value: 'Id' },
  { label: 'User Id', value: 'UserId' },
  { label: 'Stack Name', value: 'StackName' },
  { label: 'Fleet Name', value: 'FleetName' },
];

export const WS_DESC_PROPS: Array<SelectableValue<SammAwsDescWsProps>> = [
  { label: 'Workspace Id', value: 'WorkspaceId' },
  { label: 'Directory Id', value: 'DirectoryId' },
  { label: 'User Name', value: 'UserName' },
  { label: 'IP Address', value: 'IpAddress' },
  { label: 'State', value: 'State' },
  { label: 'Bundle Id', value: 'BundleId' },
  { label: 'Computer Name', value: 'ComputerName' },
];

export const WS_CONN_PROPS: Array<SelectableValue<SammAwsConnWsProps>> = [
  { label: 'Workspace Id', value: 'WorkspaceId' },
  { label: 'State', value: 'ConnectionState' },
];

export const WS_DIR_PROPS: Array<SelectableValue<SammAwsDirWsProps>> = [
  { label: 'Directory Id', value: 'DirectoryId' },
  { label: 'Alias', value: 'Alias' },
  { label: 'Directory Name', value: 'DirectoryName' },
  { label: 'Registration Code', value: 'RegistrationCode' },
  { label: 'Customer User Name', value: 'CustomerUserName' },
  { label: 'IAM Role Id', value: 'IamRoleId' },
  { label: 'Directory Type', value: 'DirectoryType' },
  { label: 'Security Group Id', value: 'WorkspaceSecurityGroupId' },
  { label: 'State', value: 'State' },
];

export const NONE_PROPS: Array<SelectableValue<SammAwsNoneProps>> = [
  { label: 'None', value: 'None' },
];

export const FILTERS_VALUE_QUERY: Array<FilterQueryDefinition> = [
  {
    key: 'BundleId', 
    query: { 
      ...DefaultSammAwsQuery,
      service: 'workspaces', 
      service_query: 'DescribeWorkspaceBundles', 
      fieldList: [ 'BundleId', 'BundleId' ],
    }
  },
  {
    key: 'BundleName',
    query: { 
      ...DefaultSammAwsQuery,
      service: 'workspaces', 
      service_query: 'DescribeWorkspaceBundles', 
      fieldList: [ 'Name', 'BundleId' ],
    } 
  },
  {
    key: 'DirectoryId',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaceDirectories',
      fieldList: [ 'DirectoryId', 'DirectoryId' ],
    }
  },
  {
    key: 'DirectoryName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaceDirectories',
      fieldList: [ 'DirectoryName', 'DirectoryId' ],
    }
  },
  {
    key: 'WorkspaceId',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaces',
      fieldList: [ 'WorkspaceId', 'WorkspaceId' ],
    }
  },
  {
    key: 'WorkspaceName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaces',
      fieldList: [ 'WorkspaceName', 'WorkspaceId' ],
    }
  },
  {
    key: 'ComputerName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaces',
      fieldList: [ 'ComputerName', 'WorkspaceId' ],
    }
  },
  {
    key: 'UserName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'workspaces',
      service_query: 'DescribeWorkspaces',
      fieldList: [ 'UserName', 'WorkspaceId' ],
    }
  },
  {
    key: 'StackName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'appstream',
      service_query: 'DescribeStacks',
      fieldList: [ 'Name', 'Name' ],
    },
  },
  {
    key: 'FleetName',
    query: {
      ...DefaultSammAwsQuery,
      service: 'appstream',
      service_query: 'DescribeFleets',
      fieldList: [ 'Name', 'Name' ],
    },
  },
  {
    key: 'UserId',
    query: {
      ...DefaultSammAwsQuery,
      service: 'appstream',
      service_query: 'DescribeSessions',
      fieldList: [ 'Name', 'Name' ],
    },
    static_options: [],
  },
  {
    key: 'AuthenticationType',
    query: {
      ...DefaultSammAwsQuery
    },
    static_options: [
      { label: 'SAML', value: 'SAML' },
      { label: 'API', value: 'API' },
      { label: 'USERPOOL', value: 'USERPOOL' },
      { label: 'AWS_AD', value: 'AWS_AD' },
    ]
  }
];

export const WS_DESCRIBE_FILTERS_CC: SammAwsFilterCascader[] = [
  { label: 'Bundle Id',      value: 'BundleId' },
  { label: 'Bundle Name',    value: 'BundleName' },
  { label: 'Directory Id',   value: 'DirectoryId' },
  { label: 'Directory Name', value: 'DirectoryName' },
  { label: 'User',           value: 'UserName' },
  { label: 'Workspace Id',   value: 'WorkspaceId', },
  { label: 'Computer Name',  value: 'ComputerName' },
];

export const WS_CONNECTION_FILTERS_CC: SammAwsFilterCascader[] = [
  { label: 'Workspace Id', value: 'WorkspaceId' },
  { label: 'Workspace Name', value: 'WorkspaceName' },
  { label: 'Computer Name', value: 'ComputerName' },
];

export const WS_DIRECTORY_FILTERS_CC: SammAwsFilterCascader[] = [
  { label: 'Directory Name', value: 'DirectoryName' },
  { label: 'Directory Id', value: 'DirectoryId' },
];

export const AS_STACKS_FILTER_CC: SammAwsFilterCascader[] = [
  { label: 'Stack Name', value: 'StackName', filterField: 'Name' },
]

export const AS_FLEETS_FILTER_CC: SammAwsFilterCascader[] = [
  { label: 'Fleet Name', value: 'FleetName', filterField: 'Name' },
]

export const AS_SESSIONS_FILTER_CC: SammAwsFilterCascader[] = [
  { label: 'Fleet Name', value: 'FleetName' },
  { label: 'Stack Name', value: 'StackName' },
  { label: 'User', value: 'UserId' },
  { label: 'Authentication Type', value: 'AuthenticationType'}
]

export const AS_ASSOCIATED_STACKS_FILTER_CC: SammAwsFilterCascader[] = [
  { label: 'Fleet Name', value: 'FleetName' },
]

export const AS_ASSOCIATED_FLEETS_FILTER_CC: SammAwsFilterCascader[] = [
  { label: 'Stack Name', value: 'StackName' },
]

export const FILTER_PROPERTIES = [
  {service: 'workspaces', service_query: 'DescribeWorkspaces', filter: WS_DESCRIBE_FILTERS_CC},
  {service: 'workspaces', service_query: 'DescribeWorkspacesConnectionStatus', filter: WS_CONNECTION_FILTERS_CC},
  {service: 'workspaces', service_query: 'DescribeWorkspaceDirectories', filter: WS_DIRECTORY_FILTERS_CC},
  {service: 'appstream',  service_query: 'DescribeStacks', filter: AS_STACKS_FILTER_CC},
  {service: 'appstream',  service_query: 'DescribeFleets', filter: AS_FLEETS_FILTER_CC},
  {service: 'appstream',  service_query: 'DescribeSessions', filter: AS_SESSIONS_FILTER_CC },
  {service: 'appstream',  service_query: 'ListAssociatedStacks', filter: AS_ASSOCIATED_STACKS_FILTER_CC },
  {service: 'appstream',  service_query: 'ListAssociatedFleets', filter: AS_ASSOCIATED_FLEETS_FILTER_CC },
];
