import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';
import type {CascaderOption} from '@grafana/ui';

export type SammAwsService = 'workspaces' | 'appstream' | 'ec2';
export type SammAwsServiceQuery = (SammAwsWorkspacesServiceQuery | SammAwsAppstreamServiceQuery);
export type SammAwsWorkspacesServiceQuery = 'DescribeWorkspaces' | 'DescribeWorkspacesConnectionStatus' | 'DescribeWorkspaceDirectories' | 'DescribeWorkspaceBundles';
export type SammAwsAppstreamServiceQuery = 'DescribeStacks' |
    'DescribeFleets' | 
    'DescribeSessions' |
    'DescribeDirectoryConfigs' |
    'ListAssociatedStacks' |
    'ListAssociatedFleets';

export type SammAwsProps = (SammAwsWorkspacesProps | SammAwsAppstreamProps | SammAwsEc2Props);
export type SammAwsNoneProps = 'None';

export type SammAwsWorkspacesProps = SammAwsNoneProps | SammAwsDescWsProps | SammAwsConnWsProps | SammAwsDirWsProps;
export type SammAwsDescWsProps = 'WorkspaceId' | 'DirectoryId' | 'UserName' | 'IpAddress' | 'State' | 'BundleId' | 'SubnetId' | 'ComputerName';
export type SammAwsConnWsProps = 'WorkspaceId' | 'ConnectionState';
export type SammAwsDirWsProps = 'DirectoryId' | 'Alias' | 'DirectoryName' | 'RegistrationCode' | 'CustomerUserName' | 'IamRoleId' | 'DirectoryType' | 'WorkspaceSecurityGroupId' | 'State';

export type SammAwsAppstreamProps = SammAwsNoneProps | SammAwsDescAsFleetProps | SammAwsDescAsStackProps | SammAwsListAssociatedProps | SammAwsDescAsSessionProps;
export type SammAwsDescAsFleetProps = "Arn" | "Description" | "DisplayName" | "FleetType" | "IamRoleArn" | "ImageArn" | "ImageName" | "InstanceType" | "Name" | "Platform" | "State" | "StreamView";
export type SammAwsDescAsStackProps = "Arn" | "CreatedTime" | "Description" | "DisplayName" | "FeedbackURL" | "Name" | "RedirectURL";
export type SammAwsListAssociatedProps = "Names";
export type SammAwsDescAsSessionProps = "Id" | "UserId" | "StackName" | "FleetName";

export type SammAwsEc2Props = SammAwsNoneProps;

export type SammAwsWSDescribeFilter = 'none' | 'BundleId' | 'DirectoryId' | 'UserName' | 'WorkspaceId';
export type SammAwsWSConnectionFilter = 'none' | 'WorkspaceId';
export type SammAwsWSDescribeDirectoryFilter = 'none' | 'DirectoryId' | 'DirectoryName';

export type SammAwsASDescribeSessionFilter = 'none' | 'FleetName' | 'StackName';

export type SammAwsQueryWithDataSource<T extends SammAwsService> = {};

export type SammAwsQuery = (SammAwsWorkspacesQuery | SammAwsAppstreamQuery | SammAwsEc2Query);

export type SammAwsBaseQuery = (
  | DataQuery & {
    service: SammAwsService | null | undefined;
    service_query: SammAwsServiceQuery | null | undefined;
    text_prop?: string;
    value_prop?: string;
    query?: string;
    filter_prop?: string;
    filter_value?: string;
    Limit?: number;
    filterConditions: Array<FilterCondition>;
    fieldList: Array<string | undefined>;
  }
);

export type SammAwsWorkspacesQuery = (
  | SammAwsBaseQuery & { 
    filter_field?: string;
    BundleId?: string;
    DirectoryId?: string;
    UserName?: string;
    WorkspaceId?: string;
  }
);
export interface SammAwsAppstreamQuery extends SammAwsBaseQuery {
    filter_field?: string;
    filter_value?: string;
    StackName?: string;
    FleetName?: string;
    authenticationType?: string | undefined;
    userId?: string | undefined;
};

export type SammAwsEc2Query = (
  | SammAwsBaseQuery & {
    filter_field?: string;
  }
);

export interface MyQuery extends DataQuery {
  queryText?: string;
  service: string;
  service_query: string;
  constant: number;
}

export const DEFAULT_QUERY: Partial<SammAwsQuery> = {
  service: "workspaces",
  service_query: "DescribeWorkspaces",
};

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface SammAwsDataSourceOptions extends DataSourceJsonData {
  region: string;
  accessKey: string;
  maxRetries: number;
  minRetryDelay: number;
  minThrottleDelay: number;
  maxRetryDelay: number;
  maxThrottleDelay: number;
  cacheSeconds: number;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface SammAwsSecureJsonData {
  accessSecret?: string;
  accessToken?: string;
}

export interface FilterCondition {
  property: string;
  value: string;
  outProperty: string;
}

export type FilterQueryDefinition = ({
  key: string,
  query: SammAwsQuery,
  static_options?: Array<CascaderOption>,
});

export type FilterMapping = {
  value_prop: string;
  text_prop: string;
  serviceQuery: string;
}

export interface SammAwsFilterCascader extends CascaderOption {
  filterField?: string;
};
