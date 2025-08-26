import React, { useState } from 'react';
import { 
  SammAwsQuery,
  SammAwsAppstreamServiceQuery,
} from './../types';
import {
  DefaultSammAwsQuery,
} from './../constants'
import { 
  SammAwsDataSource,
  MetricFindValue,
} from '../datasource'
import { 
  AwsServiceSelector, 
  AwsServiceQuerySelector,
  AwsFieldsSelector,
  AwsFilterSelector,
} from './AwsSelector';

interface VariableQueryProps {
  query: SammAwsQuery;
  onChange: (query: SammAwsQuery) => void;
  datasource: SammAwsDataSource;
}

export const VariableQueryEditor = (props: VariableQueryProps) => {
  if (!props.query.service) {
    props.query = {...DefaultSammAwsQuery, ...props.query};
    props.query.fieldList = ["", ""];
  }

  const [ _, setFieldList ] = useState<Array<string | undefined>>(props.query.fieldList);
  const [ availableFields, setAvailableFields ] = useState<MetricFindValue[]>([])

  const loadFields = (query: SammAwsQuery) => {
    if (!query.service || !query.service_query) return;
    if (availableFields.length > 0) return;
    const q: SammAwsQuery = {
      ...DefaultSammAwsQuery,
      service: query.service,
      service_query: query.service_query + "Fields" as SammAwsAppstreamServiceQuery,
      refId: 'GetFields',
    }
    props.datasource.metricFindQuery(q).then(result => setAvailableFields(result))
  }
  loadFields(props.query);

  const saveQuery = (query: SammAwsQuery) => {
    setAvailableFields([]);
    query.fieldList = [undefined, undefined];
    props.onChange({...query, 
      query: `${query.service}.${query.service_query}(fieldList[0]=${query.fieldList[0]},fieldList[1]=${query.fieldList[1]})`});
  };

  const fieldOnChange = (query: SammAwsQuery) => {
    props.onChange(query);
  }

  return (
    <>
        <AwsServiceSelector {...props} onChange={saveQuery} />
        <AwsServiceQuerySelector {...props} onChange={saveQuery} />
        <AwsFieldsSelector 
            {...props} 
            onChange={fieldOnChange} 
            availableFields={availableFields}
            setFieldList={setFieldList}
            allowAddRemove={false}
            />
        <AwsFilterSelector {...props} onChange={props.onChange} />

    </>
  );
};

