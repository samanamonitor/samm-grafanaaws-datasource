import React, { ChangeEvent } from 'react';
import { FieldSet, InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SammAwsDataSourceOptions, SammAwsSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<SammAwsDataSourceOptions, SammAwsSecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;

  jsonData.maxRetries = jsonData.maxRetries ?? 5;
  jsonData.minRetryDelay = jsonData.minRetryDelay ?? 100;
  jsonData.maxRetryDelay = jsonData.maxRetryDelay ?? 1000;
  jsonData.minThrottleDelay = jsonData.minThrottleDelay ?? 500;
  jsonData.maxThrottleDelay = jsonData.maxThrottleDelay ?? 30000;
  jsonData.cacheSeconds = jsonData.cacheSeconds ?? 3600;

  const onRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        region: event.target.value,
      },
    });
  };

  // AWS Access Key
  const onAccessKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        accessKey: event.target.value,
      },
    });
  };

  // AWS Access Secret
  const onSecureChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
		  ...options.secureJsonData,
        [event.target.name]: event.target.value,
      },
    });
  };

  const onResetAccessSecret = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        accessSecret: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        accessSecret: '',
      },
    });
  };

  const onResetAccessToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        accessToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        accessToken: '',
      },
    });
  };

  return (
    <>
      <FieldSet>
        <InlineField label="Region" labelWidth={20} interactive tooltip={'AWS Region'}>
          <Input
            required
            id="config-editor-region"
            onChange={onRegionChange}
            value={jsonData.region}
            placeholder="Enter AWS Region, e.g. us-east-1"
            width={40}
          />
        </InlineField>
        <InlineField label="Access Key" labelWidth={20} interactive tooltip={'AWS Access Key ID'}>
          <Input
            required
            id="config-editor-access-key"
            onChange={onAccessKeyChange}
            value={jsonData.accessKey}
            placeholder="Enter your Access Key"
            width={40}
          />
        </InlineField>
        <InlineField label="Access Secret" labelWidth={20} interactive tooltip={'Secure AWS Access Secret'}>
          <SecretInput
          name="accessSecret"
            required
            id="config-editor-access-secret"
            isConfigured={secureJsonFields.accessSecret}
            value={secureJsonData?.accessSecret}
            placeholder="Enter your Access secret"
            width={40}
            onReset={onResetAccessSecret}
            onChange={onSecureChange}
          />
        </InlineField>
        <InlineField label="Access Token" labelWidth={20} interactive tooltip={'Secure AWS Access Token (optional)'}>
          <SecretInput
          name="accessToken"
            id="config-editor-access-token"
            isConfigured={secureJsonFields.accessToken}
            value={secureJsonData?.accessToken}
            placeholder="Enter your Access Token"
            width={40}
            onReset={onResetAccessToken}
            onChange={onSecureChange}
          />
        </InlineField>
      </FieldSet>
      <div className='gf-form-group'>
        <h3 className='page-heading'>Additional settings</h3>
        <FieldSet>
          <InlineField label="Max Retries" labelWidth={25} interactive tooltip={'Max number of query retries to AWS'}>
            <Input
              required
              id="config-editor-max-retries"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  maxRetries: Number(event.target.value),
                },
                });}}
              value={jsonData.maxRetries}
              placeholder="5"
              width={15}
            />
          </InlineField>
          <InlineField label="Min Retry Delay (ms)" labelWidth={25} interactive tooltip={'Min time in milliseconds to delay between retries'}>
            <Input
              required
              id="config-editor-min-retry-delay"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  minRetryDelay: Number(event.target.value),
                },
                });}}
              value={jsonData.minRetryDelay}
              placeholder="500"
              width={15}
            />
          </InlineField>
          <InlineField label="Max Retry Delay (ms)" labelWidth={25} interactive tooltip={'Max time in milliseconds to delay between retries'}>
            <Input
              required
              id="config-editor-max-retry-delay"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  maxRetryDelay: Number(event.target.value),
                },
                });}}
              value={jsonData.maxRetryDelay}
              placeholder="1000"
              width={15}
            />
          </InlineField>
          <InlineField label="Min Throttle Delay (ms)" labelWidth={25} interactive tooltip={'Min time in milliseconds to delay between retries when throttled'}>
            <Input
              required
              id="config-editor-min-throttle-delay"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  minThrottleDelay: Number(event.target.value),
                },
                });}}
              value={jsonData.minThrottleDelay}
              placeholder="500"
              width={15}
            />
          </InlineField>
          <InlineField label="Max Throttle Delay (ms)" labelWidth={25} interactive tooltip={'Max time in milliseconds to delay between retries when throttled'}>
            <Input
              required
              id="config-editor-max-throttle-delay"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  maxThrottleDelay: Number(event.target.value),
                },
                });}}
              value={jsonData.maxThrottleDelay}
              placeholder="30000"
              width={15}
            />
          </InlineField>
          <InlineField label="Cache Expiration (s)" labelWidth={25} interactive tooltip={'Time in seconds that the plugin will keep the results in cache.'}>
            <Input
              required
              id="config-editor-cache-seconds"
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({
                ...options,
                jsonData: {
                  ...jsonData,
                  cacheSeconds: Number(event.target.value),
                },
                });}}
              value={jsonData.cacheSeconds}
              placeholder="3600"
              width={15}
            />
          </InlineField>
        </FieldSet>
      </div>
    </>
  );
}
