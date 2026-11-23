// This is a generated file - do not edit.
//
// Generated from forgenet/v1/marketplace.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use templateVariableDescriptor instead')
const TemplateVariable$json = {
  '1': 'TemplateVariable',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'label', '3': 2, '4': 1, '5': 9, '10': 'label'},
    {'1': 'help', '3': 3, '4': 1, '5': 9, '10': 'help'},
    {'1': 'default_value', '3': 4, '4': 1, '5': 9, '10': 'defaultValue'},
    {'1': 'pattern', '3': 5, '4': 1, '5': 9, '10': 'pattern'},
    {'1': 'optional', '3': 6, '4': 1, '5': 8, '10': 'optional'},
    {'1': 'secret', '3': 7, '4': 1, '5': 8, '10': 'secret'},
    {'1': 'generate', '3': 8, '4': 1, '5': 13, '10': 'generate'},
  ],
};

/// Descriptor for `TemplateVariable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templateVariableDescriptor = $convert.base64Decode(
    'ChBUZW1wbGF0ZVZhcmlhYmxlEhIKBG5hbWUYASABKAlSBG5hbWUSFAoFbGFiZWwYAiABKAlSBW'
    'xhYmVsEhIKBGhlbHAYAyABKAlSBGhlbHASIwoNZGVmYXVsdF92YWx1ZRgEIAEoCVIMZGVmYXVs'
    'dFZhbHVlEhgKB3BhdHRlcm4YBSABKAlSB3BhdHRlcm4SGgoIb3B0aW9uYWwYBiABKAhSCG9wdG'
    'lvbmFsEhYKBnNlY3JldBgHIAEoCFIGc2VjcmV0EhoKCGdlbmVyYXRlGAggASgNUghnZW5lcmF0'
    'ZQ==');

@$core.Deprecated('Use templatePortDescriptor instead')
const TemplatePort$json = {
  '1': 'TemplatePort',
  '2': [
    {'1': 'service', '3': 1, '4': 1, '5': 9, '10': 'service'},
    {'1': 'port', '3': 2, '4': 1, '5': 13, '10': 'port'},
    {'1': 'protocol', '3': 3, '4': 1, '5': 9, '10': 'protocol'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `TemplatePort`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templatePortDescriptor = $convert.base64Decode(
    'CgxUZW1wbGF0ZVBvcnQSGAoHc2VydmljZRgBIAEoCVIHc2VydmljZRISCgRwb3J0GAIgASgNUg'
    'Rwb3J0EhoKCHByb3RvY29sGAMgASgJUghwcm90b2NvbBIgCgtkZXNjcmlwdGlvbhgEIAEoCVIL'
    'ZGVzY3JpcHRpb24=');

@$core.Deprecated('Use templateVolumeDescriptor instead')
const TemplateVolume$json = {
  '1': 'TemplateVolume',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `TemplateVolume`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templateVolumeDescriptor = $convert.base64Decode(
    'Cg5UZW1wbGF0ZVZvbHVtZRISCgRuYW1lGAEgASgJUgRuYW1lEiAKC2Rlc2NyaXB0aW9uGAIgAS'
    'gJUgtkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use templateDescriptor instead')
const Template$json = {
  '1': 'Template',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'category', '3': 4, '4': 1, '5': 9, '10': 'category'},
    {'1': 'version', '3': 5, '4': 1, '5': 9, '10': 'version'},
    {'1': 'license', '3': 6, '4': 1, '5': 9, '10': 'license'},
    {'1': 'homepage', '3': 7, '4': 1, '5': 9, '10': 'homepage'},
    {'1': 'icon_svg', '3': 8, '4': 1, '5': 9, '10': 'iconSvg'},
    {
      '1': 'variables',
      '3': 9,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.TemplateVariable',
      '10': 'variables'
    },
    {
      '1': 'ports',
      '3': 10,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.TemplatePort',
      '10': 'ports'
    },
    {
      '1': 'volumes',
      '3': 11,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.TemplateVolume',
      '10': 'volumes'
    },
  ],
};

/// Descriptor for `Template`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templateDescriptor = $convert.base64Decode(
    'CghUZW1wbGF0ZRIOCgJpZBgBIAEoCVICaWQSEgoEbmFtZRgCIAEoCVIEbmFtZRIgCgtkZXNjcm'
    'lwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24SGgoIY2F0ZWdvcnkYBCABKAlSCGNhdGVnb3J5EhgK'
    'B3ZlcnNpb24YBSABKAlSB3ZlcnNpb24SGAoHbGljZW5zZRgGIAEoCVIHbGljZW5zZRIaCghob2'
    '1lcGFnZRgHIAEoCVIIaG9tZXBhZ2USGQoIaWNvbl9zdmcYCCABKAlSB2ljb25TdmcSOwoJdmFy'
    'aWFibGVzGAkgAygLMh0uZm9yZ2VuZXQudjEuVGVtcGxhdGVWYXJpYWJsZVIJdmFyaWFibGVzEi'
    '8KBXBvcnRzGAogAygLMhkuZm9yZ2VuZXQudjEuVGVtcGxhdGVQb3J0UgVwb3J0cxI1Cgd2b2x1'
    'bWVzGAsgAygLMhsuZm9yZ2VuZXQudjEuVGVtcGxhdGVWb2x1bWVSB3ZvbHVtZXM=');

@$core.Deprecated('Use listTemplatesRequestDescriptor instead')
const ListTemplatesRequest$json = {
  '1': 'ListTemplatesRequest',
};

/// Descriptor for `ListTemplatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTemplatesRequestDescriptor =
    $convert.base64Decode('ChRMaXN0VGVtcGxhdGVzUmVxdWVzdA==');

@$core.Deprecated('Use listTemplatesResponseDescriptor instead')
const ListTemplatesResponse$json = {
  '1': 'ListTemplatesResponse',
  '2': [
    {
      '1': 'templates',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.Template',
      '10': 'templates'
    },
  ],
};

/// Descriptor for `ListTemplatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTemplatesResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0VGVtcGxhdGVzUmVzcG9uc2USMwoJdGVtcGxhdGVzGAEgAygLMhUuZm9yZ2VuZXQudj'
    'EuVGVtcGxhdGVSCXRlbXBsYXRlcw==');

@$core.Deprecated('Use getTemplateRequestDescriptor instead')
const GetTemplateRequest$json = {
  '1': 'GetTemplateRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTemplateRequestDescriptor =
    $convert.base64Decode('ChJHZXRUZW1wbGF0ZVJlcXVlc3QSDgoCaWQYASABKAlSAmlk');

@$core.Deprecated('Use getTemplateResponseDescriptor instead')
const GetTemplateResponse$json = {
  '1': 'GetTemplateResponse',
  '2': [
    {
      '1': 'template',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Template',
      '10': 'template'
    },
    {'1': 'readme', '3': 2, '4': 1, '5': 9, '10': 'readme'},
  ],
};

/// Descriptor for `GetTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTemplateResponseDescriptor = $convert.base64Decode(
    'ChNHZXRUZW1wbGF0ZVJlc3BvbnNlEjEKCHRlbXBsYXRlGAEgASgLMhUuZm9yZ2VuZXQudjEuVG'
    'VtcGxhdGVSCHRlbXBsYXRlEhYKBnJlYWRtZRgCIAEoCVIGcmVhZG1l');

@$core.Deprecated('Use installTemplateRequestDescriptor instead')
const InstallTemplateRequest$json = {
  '1': 'InstallTemplateRequest',
  '2': [
    {'1': 'template_id', '3': 1, '4': 1, '5': 9, '10': 'templateId'},
    {'1': 'project_id', '3': 2, '4': 1, '5': 9, '10': 'projectId'},
    {
      '1': 'values',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.InstallTemplateRequest.ValuesEntry',
      '10': 'values'
    },
    {'1': 'request_id', '3': 4, '4': 1, '5': 9, '10': 'requestId'},
  ],
  '3': [InstallTemplateRequest_ValuesEntry$json],
};

@$core.Deprecated('Use installTemplateRequestDescriptor instead')
const InstallTemplateRequest_ValuesEntry$json = {
  '1': 'ValuesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InstallTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List installTemplateRequestDescriptor = $convert.base64Decode(
    'ChZJbnN0YWxsVGVtcGxhdGVSZXF1ZXN0Eh8KC3RlbXBsYXRlX2lkGAEgASgJUgp0ZW1wbGF0ZU'
    'lkEh0KCnByb2plY3RfaWQYAiABKAlSCXByb2plY3RJZBJHCgZ2YWx1ZXMYAyADKAsyLy5mb3Jn'
    'ZW5ldC52MS5JbnN0YWxsVGVtcGxhdGVSZXF1ZXN0LlZhbHVlc0VudHJ5UgZ2YWx1ZXMSHQoKcm'
    'VxdWVzdF9pZBgEIAEoCVIJcmVxdWVzdElkGjkKC1ZhbHVlc0VudHJ5EhAKA2tleRgBIAEoCVID'
    'a2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use generatedSecretDescriptor instead')
const GeneratedSecret$json = {
  '1': 'GeneratedSecret',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'label', '3': 2, '4': 1, '5': 9, '10': 'label'},
    {'1': 'value', '3': 3, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `GeneratedSecret`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generatedSecretDescriptor = $convert.base64Decode(
    'Cg9HZW5lcmF0ZWRTZWNyZXQSEgoEbmFtZRgBIAEoCVIEbmFtZRIUCgVsYWJlbBgCIAEoCVIFbG'
    'FiZWwSFAoFdmFsdWUYAyABKAlSBXZhbHVl');

@$core.Deprecated('Use installTemplateResponseDescriptor instead')
const InstallTemplateResponse$json = {
  '1': 'InstallTemplateResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
    {
      '1': 'generated_secrets',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.GeneratedSecret',
      '10': 'generatedSecrets'
    },
  ],
};

/// Descriptor for `InstallTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List installTemplateResponseDescriptor = $convert.base64Decode(
    'ChdJbnN0YWxsVGVtcGxhdGVSZXNwb25zZRI3CgpkZXBsb3ltZW50GAEgASgLMhcuZm9yZ2VuZX'
    'QudjEuRGVwbG95bWVudFIKZGVwbG95bWVudBJJChFnZW5lcmF0ZWRfc2VjcmV0cxgCIAMoCzIc'
    'LmZvcmdlbmV0LnYxLkdlbmVyYXRlZFNlY3JldFIQZ2VuZXJhdGVkU2VjcmV0cw==');
