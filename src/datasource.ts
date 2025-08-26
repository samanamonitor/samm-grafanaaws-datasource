import { Observable, of, from } from 'rxjs';
import { map } from 'rxjs/operators';
import { 
    DataQueryRequest, 
    DataQueryResponse,
    DataSourceInstanceSettings, 
    ScopedVars,
    CustomVariableSupport,
} from '@grafana/data';
import { 
    DataSourceWithBackend, 
    getTemplateSrv,
    TemplateSrv, 
} from '@grafana/runtime';

import {
    SammAwsQuery, 
    SammAwsDataSourceOptions, 
} from './types';
import {
    VariableQueryEditor,
} from './components/VariableQueryEditor'
import { CascaderOption } from '@grafana/ui';

export interface MetricFindValue extends CascaderOption {
    text: string;
    data?: any;
}

export class SammAwsVariableSupport extends CustomVariableSupport<SammAwsDataSource> {
    constructor(
        private readonly datasource: SammAwsDataSource,
        private readonly templateSrv: TemplateSrv = getTemplateSrv()
    ) {
        super();
    }

    editor = VariableQueryEditor;

    query(request: DataQueryRequest): Observable<DataQueryResponse> {
        const query = request.targets[0] as SammAwsQuery;
        const nf = query.filterConditions.map((item, index) => {
            return {
                ...item, 
                value: this.templateSrv.replace(item.value, request.scopedVars)
            }});
        const newquery = {...query, filterConditions: nf};
        if (query.service && query.service_query && 
            query.fieldList && query.fieldList.length == 2 
            && query.fieldList[0] && query.fieldList[1]) {
            return from(this.datasource.metricFindQuery(newquery)).pipe(map((results) => ({data: results})));
        }
        return of({ data: []});
    }
}

export class SammAwsDataSource extends DataSourceWithBackend<SammAwsQuery, SammAwsDataSourceOptions> {

    constructor(instanceSettings: DataSourceInstanceSettings<SammAwsDataSourceOptions>,
                  private readonly templateSrv: TemplateSrv = getTemplateSrv(),
                ) {
        super(instanceSettings);
        this.variables = new SammAwsVariableSupport(this, this.templateSrv);
    }

    applyTemplateVariables(query: SammAwsQuery, scopedVars: ScopedVars) {
        if (query.filterConditions) {
            return {
                ...query,
                filterConditions: query.filterConditions.map((item, index) => {
                    return {
                        ...item, 
                        value: this.templateSrv.replace(item.value, scopedVars)
                    };
                }),
            };
        }
        return query;
    }

    async metricFindQuery(query: SammAwsQuery, options?: any): Promise<MetricFindValue[]> {
        const newquery = this.applyTemplateVariables(query, {});
        return this.postResource('query', newquery).then((result) => {
            return this.queryToLabel(result as MetricFindValue[])
        });
    }

    async getLabels(query: SammAwsQuery, options?: any): Promise<MetricFindValue[]> {
        const newquery = this.applyTemplateVariables(query, {});
        return this.postResource('query', newquery).then((result) => result as MetricFindValue[]);
    }

    queryToLabel(result: MetricFindValue[]): MetricFindValue[] {
        let res: MetricFindValue[] = []
        if (result.length === 0 || result[0].data.values.length === 0) {
            console.log(result);
            return res;
        }
        for (let i = 0; i < result[0].data.values[0].length; i++) {
            res.push({
                label: result[0].data.values[0][i],
                text: result[0].data.values[0][i],
                value: result[0].data.values[1][i],
            });
        }
        return res;
    }

    getVariables(): string[] {
        return this.templateSrv.getVariables().map((v) => `$${v.name}`);
    }
}
