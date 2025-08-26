import React, { useState } from 'react'; 
import type { 
    SammAwsService,
    SammAwsQuery,
    SammAwsServiceQuery,
    FilterCondition,
} from '../types'
import { 
    SammAwsDataSource,
    MetricFindValue,
 } from '../datasource'
import type { SelectableValue } from '@grafana/data';
import { 
    Select, 
    Input, 
    InlineFormLabel,
    Cascader,
    Button,
    CascaderOption,
} from '@grafana/ui';
import { 
    SERVICE_TYPES,
    SERVICE_QUERY_TYPES,
    EnableFilters,
    FILTERS_VALUE_QUERY,
    FILTER_PROPERTIES,
} from '../constants';


export interface AwsQueryProps {
    query: SammAwsQuery;
    onChange: (query: SammAwsQuery) => void;
    datasource: SammAwsDataSource;
}

export const AwsLimitSelector = (props: AwsQueryProps) => {
    return (
        <div className="gf-form">
            <InlineFormLabel>
                Limit
            </InlineFormLabel>
            <Input 
                width={18}
                value={props.query.Limit} 
                onChange={(e) => props.onChange({ ...props.query, Limit: Number((e.target as HTMLInputElement).value) })} />
        </div>
    );
}

export const AwsServiceSelector = (props: AwsQueryProps) => {

    const { query, onChange } = props;

    const onServiceChange = (e: SelectableValue<SammAwsService>) => {
        let newval = {...query, service: e.value, service_query: null}
        onChange(newval);
    }

    return (
        <>
        <div className="gf-form">
            <InlineFormLabel>
                Service
            </InlineFormLabel>
            <Select 
                width={18}
                options={SERVICE_TYPES}
                onChange={onServiceChange}
                value={query.service}
                key={query.service}
                menuShouldPortal={true}
            />
        </div>
        </>
    );
}

export const AwsServiceQuerySelector = (props: AwsQueryProps) => {
    const { query, onChange } = props;
    const onServiceQueryChange = (e: SelectableValue<SammAwsServiceQuery>) => {
        let newval = {...query, service_query: e.value}
        onChange(newval);
    }
    let serviceQueryTypes: Array<SelectableValue<SammAwsServiceQuery>> | undefined = [];
    serviceQueryTypes = SERVICE_QUERY_TYPES.find((e) => 
        e.service === query.service)?.service_queries;

    return (
        <>
            <div className="gf-form"> {/* ServiceQuery */}
            <InlineFormLabel>
                Service Query
            </InlineFormLabel>
            <Select 
                width={28}
                options={serviceQueryTypes}
                onChange={onServiceQueryChange}
                value={query.service_query}
                menuShouldPortal={true}
            />
            </div>
        </>
    )
}

export interface AwsFieldsProps {
    query: SammAwsQuery;
    onChange: (query: SammAwsQuery) => void;
    datasource: SammAwsDataSource;
    availableFields: MetricFindValue[];
    setFieldList: React.Dispatch<React.SetStateAction<Array<string | undefined>>>;
    allowAddRemove?: boolean;
}

export const AwsFieldsSelector = (props: AwsFieldsProps) => {
    props.allowAddRemove = props.allowAddRemove ?? true;
    const addField = (_: any): void => {
        props.query.fieldList = [...props.query.fieldList, ""];
        props.setFieldList(props.query.fieldList);
    }
    const removeField = (index: number): void => {
        props.setFieldList(prevItems => {
          props.query.fieldList = prevItems.filter((_, i) => i !== index);
          return props.query.fieldList;
        });
    }
      
    
    const fieldSelector = (fieldName: string | undefined, index: number, availableFields: Array<MetricFindValue>) => {
        return (
          <>
            <InlineFormLabel width={10} tooltip={'Add filter condition'}>
              Field
            </InlineFormLabel>
            <Cascader
              width={40}
              onSelect={(item) => {
                if (item === fieldName) {
                  return;
                }
                if (item) {
                  props.query.fieldList[index] = item
                }
                props.onChange(props.query);
                console.log("props");
                console.log(props);
              }}
              changeOnSelect={true}
              allowCustomValue={false}
              initialValue={fieldName}
              options={availableFields}
            />
          </>
        )
      }
    
    const fieldSelectors = props.query.fieldList.map((fieldName, index) => {
        if (props.availableFields.length == 0) {
            return;
        }
        return (
          <div className='gf-form-inline'>
            <div className='gf-form'>
              {fieldSelector(fieldName, index, props.availableFields)}
            </div>
            {props.allowAddRemove && (
                <Button variant={'secondary'} onClick={() => removeField(index)}>
                -
                </Button>
            )}
          </div>
        )
      });
    
    return (
        <>
          {fieldSelectors}
          {props.allowAddRemove && (
            <div className="gf-form-inline">
                <div className={'gf-form'}>
                    <Button variant={'secondary'} onClick={addField}>
                    + Field
                    </Button>
                </div>
            </div>

          )}
        </>
    )
}

export interface AwsFilterProps {
    query: SammAwsQuery;
    onChange: (query: SammAwsQuery) => void;
    datasource: SammAwsDataSource;

}

export const AwsFilterSelector = (props: AwsFilterProps) => {
    const [ filterConditions, setFilterConditions ] = useState<Array<FilterCondition>>(props.query.filterConditions);
    const [ filterValues, setFilterValues ] = useState<Array<CascaderOption>>([]);

    const filterConditionProperties = FILTER_PROPERTIES.find((e) => 
        (e.service === props.query.service && e.service_query === props.query.service_query));
    
    const addFilterCondition = (_: any): void => {
        props.query.filterConditions = [...filterConditions, { property: '', value: '', outProperty: '' }];
        setFilterConditions(props.query.filterConditions);
    }

    const removeFilterCondition = (index: number): void => {
        setFilterConditions(prevItems => {
          props.query.filterConditions = prevItems.filter((_, i) => i !== index);
          return props.query.filterConditions;
        });
    }

    const loadFilterValues = (item: string, index: number) => {
        if (!item) {
          return;
        }
        const m = FILTERS_VALUE_QUERY.find((e) => e.key === item);
        if ( m === undefined) {
          console.log("No definition for " + item);
          //clearIsLoadingIndex(index);
          return;
        }
        //setIsLoadingIndex(index);
        const q: SammAwsQuery = {
          ...props.query,
          ...m.query,
          refId: 'GetFilterProperty'+item,
        }
        if (m.static_options) {
          setFilterValues( [
            ...props.datasource.getVariables().map((item) => {
              return { label: item, value: item };
            }),
            ...m.static_options,
          ]);
          return true;
        }
        props.datasource.metricFindQuery(q).then(result => setFilterValues([
            ...props.datasource.getVariables().map((item, index) => {
                return { label: item, value: item, text: item };
            }),
            ...result,
        ]))
        return true;
    }

    const updateFilterConditionValue = (newValue: string, index: number, oldValue: string): void => {
        if (newValue === oldValue) {
          return;
        }
        if (! newValue) {
          return;
        }
        const fc = props.query.filterConditions[index];
        if (fc !== undefined && fc !== null) {
          fc.value = newValue;
          if (!fc.outProperty) {
            fc.outProperty = fc.property
          }
        }
        setFilterConditions(props.query.filterConditions);
        props.onChange(props.query);
    }
    
    const filterPropertyName = (filterCondition: FilterCondition, index: number) => {
        return (
          <>
            <InlineFormLabel width={10} tooltip={'Add filter condition'}>
              {index === 0 ? 'Filter' : 'AND'}
            </InlineFormLabel>
            <Cascader
              width={40}
              disabled={filterCondition.property?.length > 0}
              onSelect={(item) => {
                if (item === filterCondition.property) {
                  return;
                }
                if (item) {
                  const filterField = filterConditionProperties?.filter.find((propitem) => propitem.value === item);
                  if (filterField) {
                    filterCondition.property = filterField.value;
                    ////filterCondition.property = filterquery.query.fieldList[1] ?? "";
                  } else {
                    filterCondition.property = item;
                  }
                }
                loadFilterValues(item, index);
                props.onChange(props.query);
    
              }}
              changeOnSelect={true}
              initialValue={filterCondition.property}
              options={filterConditionProperties?.filter ?? []}
            />
          </>
        )
    }

    const filterPropertyValue = (filterCondition: FilterCondition, index: number) => {
        if (filterCondition.property) {
          return (<>
            <Cascader
                        placeholder={"Choose"}
                        options={filterValues}
                        onSelect={item => updateFilterConditionValue(item, index, filterCondition.value)}
                        changeOnSelect={true}
                        initialValue={filterCondition.value}
                        allowCustomValue={true}
                      />
            </>);
        } 
        return (<></>)
    }
        
    const listFilters = filterConditions.map((filterCondition, index) => {
        const filter = (
            <div className="gf-form-inline">
                <div className={'gf-form'}>
                    { filterPropertyName(filterCondition, index) }
                    =
                    { filterPropertyValue(filterCondition, index) }
                    <Button variant={'secondary'} onClick={() => removeFilterCondition(index)}>
                    -
                    </Button>
                </div>
            </div>)
        return filter;
    });
    
    return (
        <>
            {listFilters}
            { EnableFilters && (
                <div className="gf-form-inline">
                <div className={'gf-form'}>
                    <Button variant={'secondary'} onClick={addFilterCondition}>
                    + Filter condition
                    </Button>
                </div>
                </div>
            )}
        </>
    )
}