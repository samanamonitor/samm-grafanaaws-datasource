import { DataSourcePlugin } from '@grafana/data';
import { SammAwsDataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { SammAwsQueryEditor } from './components/QueryEditor';

import { SammAwsQuery, SammAwsDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<SammAwsDataSource, SammAwsQuery, SammAwsDataSourceOptions>(SammAwsDataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(SammAwsQueryEditor);
