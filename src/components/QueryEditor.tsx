import React, { useState } from 'react';
import { 
  SammAwsDataSource,
  MetricFindValue,
} from '../datasource'
import type { 
  SammAwsQuery,
  SammAwsBaseQuery,
  SammAwsAppstreamServiceQuery,
//  SammAwsService,
} from '../types'
import { 
  DefaultSammAwsQuery,
  EnableLimit,
} from '../constants';
import {
  AwsLimitSelector,
  AwsServiceSelector,
  AwsServiceQuerySelector,
  AwsFieldsSelector,
  AwsFilterSelector,
} from './AwsSelector'

export type SammAwsQueryEditorProps = {
    query: SammAwsQuery;
    onChange: (query: SammAwsQuery) => void;
    onRunQuery: () => void;
    datasource: SammAwsDataSource;
};

export type FilterConditionCache = {
  property: string;
  values: MetricFindValue[]
}

export const SammAwsQueryEditor = (props: SammAwsQueryEditorProps ) => {
  if (!props.query.service) {
    props.query = {...DefaultSammAwsQuery, ...props.query};
  }

  props.query.filterConditions = props.query.filterConditions ?? [];
  const [ _, setFieldList ] = useState<Array<string | undefined>>(props.query.fieldList);
  const [ availableFields, setAvailableFields ] = useState<MetricFindValue[]>([])


  const loadFields = (query: SammAwsQuery) => {
    if (!query.service || !query.service_query || availableFields.length > 0) return;
    const q: SammAwsBaseQuery = {
      ...DefaultSammAwsQuery,
      service: query.service,
      service_query: query.service_query + "Fields" as SammAwsAppstreamServiceQuery,
      refId: 'GetFields',
    }
    props.datasource.metricFindQuery(q).then(result => setAvailableFields(result))
  }
  loadFields(props.query);

  const serviceOnChange = (query: SammAwsQuery) => {
    props.query = {
      ...query,
      filterConditions: [],
      fieldList: []
    }
    setFieldList(props.query.fieldList);
    setAvailableFields([]);
    loadFields(props.query);
    props.onChange(props.query)
  }

  return (
    <>
      <div className="gf-form-group">
        { EnableLimit && (<AwsLimitSelector {...props} />)}
        <AwsServiceSelector {...props} onChange={serviceOnChange} />
        <AwsServiceQuerySelector {...props} onChange={serviceOnChange} />
        <AwsFieldsSelector 
          {...props} 
          onChange={props.onChange} 
          availableFields={availableFields}
          setFieldList={setFieldList}
        />
        <AwsFilterSelector {...props} onChange={props.onChange} />
      </div>
    </>
  );
};
