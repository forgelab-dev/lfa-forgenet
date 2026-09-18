// This is a generated file - do not edit.
//
// Generated from forgenet/v1/deploy.proto.

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

@$core.Deprecated('Use deploymentStateDescriptor instead')
const DeploymentState$json = {
  '1': 'DeploymentState',
  '2': [
    {'1': 'DEPLOYMENT_STATE_UNSPECIFIED', '2': 0},
    {'1': 'DEPLOYMENT_STATE_PENDING', '2': 1},
    {'1': 'DEPLOYMENT_STATE_BUILDING', '2': 2},
    {'1': 'DEPLOYMENT_STATE_DEPLOYING', '2': 3},
    {'1': 'DEPLOYMENT_STATE_HEALTH_CHECKING', '2': 4},
    {'1': 'DEPLOYMENT_STATE_SUCCESS', '2': 5},
    {'1': 'DEPLOYMENT_STATE_FAILED', '2': 6},
    {'1': 'DEPLOYMENT_STATE_CANCELLED', '2': 7},
    {'1': 'DEPLOYMENT_STATE_ROLLING_BACK', '2': 8},
    {'1': 'DEPLOYMENT_STATE_ROLLED_BACK', '2': 9},
  ],
};

/// Descriptor for `DeploymentState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deploymentStateDescriptor = $convert.base64Decode(
    'Cg9EZXBsb3ltZW50U3RhdGUSIAocREVQTE9ZTUVOVF9TVEFURV9VTlNQRUNJRklFRBAAEhwKGE'
    'RFUExPWU1FTlRfU1RBVEVfUEVORElORxABEh0KGURFUExPWU1FTlRfU1RBVEVfQlVJTERJTkcQ'
    'AhIeChpERVBMT1lNRU5UX1NUQVRFX0RFUExPWUlORxADEiQKIERFUExPWU1FTlRfU1RBVEVfSE'
    'VBTFRIX0NIRUNLSU5HEAQSHAoYREVQTE9ZTUVOVF9TVEFURV9TVUNDRVNTEAUSGwoXREVQTE9Z'
    'TUVOVF9TVEFURV9GQUlMRUQQBhIeChpERVBMT1lNRU5UX1NUQVRFX0NBTkNFTExFRBAHEiEKHU'
    'RFUExPWU1FTlRfU1RBVEVfUk9MTElOR19CQUNLEAgSIAocREVQTE9ZTUVOVF9TVEFURV9ST0xM'
    'RURfQkFDSxAJ');

@$core.Deprecated('Use recreateStrategyDescriptor instead')
const RecreateStrategy$json = {
  '1': 'RecreateStrategy',
};

/// Descriptor for `RecreateStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recreateStrategyDescriptor =
    $convert.base64Decode('ChBSZWNyZWF0ZVN0cmF0ZWd5');

@$core.Deprecated('Use rollingStrategyDescriptor instead')
const RollingStrategy$json = {
  '1': 'RollingStrategy',
  '2': [
    {'1': 'max_surge', '3': 1, '4': 1, '5': 13, '10': 'maxSurge'},
    {'1': 'max_unavailable', '3': 2, '4': 1, '5': 13, '10': 'maxUnavailable'},
  ],
};

/// Descriptor for `RollingStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rollingStrategyDescriptor = $convert.base64Decode(
    'Cg9Sb2xsaW5nU3RyYXRlZ3kSGwoJbWF4X3N1cmdlGAEgASgNUghtYXhTdXJnZRInCg9tYXhfdW'
    '5hdmFpbGFibGUYAiABKA1SDm1heFVuYXZhaWxhYmxl');

@$core.Deprecated('Use blueGreenStrategyDescriptor instead')
const BlueGreenStrategy$json = {
  '1': 'BlueGreenStrategy',
  '2': [
    {
      '1': 'auto_promote_seconds',
      '3': 1,
      '4': 1,
      '5': 13,
      '10': 'autoPromoteSeconds'
    },
  ],
};

/// Descriptor for `BlueGreenStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blueGreenStrategyDescriptor = $convert.base64Decode(
    'ChFCbHVlR3JlZW5TdHJhdGVneRIwChRhdXRvX3Byb21vdGVfc2Vjb25kcxgBIAEoDVISYXV0b1'
    'Byb21vdGVTZWNvbmRz');

@$core.Deprecated('Use canaryStrategyDescriptor instead')
const CanaryStrategy$json = {
  '1': 'CanaryStrategy',
  '2': [
    {
      '1': 'initial_weight_percent',
      '3': 1,
      '4': 1,
      '5': 13,
      '10': 'initialWeightPercent'
    },
    {
      '1': 'step_weight_percent',
      '3': 2,
      '4': 1,
      '5': 13,
      '10': 'stepWeightPercent'
    },
    {
      '1': 'step_interval_seconds',
      '3': 3,
      '4': 1,
      '5': 13,
      '10': 'stepIntervalSeconds'
    },
  ],
};

/// Descriptor for `CanaryStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List canaryStrategyDescriptor = $convert.base64Decode(
    'Cg5DYW5hcnlTdHJhdGVneRI0ChZpbml0aWFsX3dlaWdodF9wZXJjZW50GAEgASgNUhRpbml0aW'
    'FsV2VpZ2h0UGVyY2VudBIuChNzdGVwX3dlaWdodF9wZXJjZW50GAIgASgNUhFzdGVwV2VpZ2h0'
    'UGVyY2VudBIyChVzdGVwX2ludGVydmFsX3NlY29uZHMYAyABKA1SE3N0ZXBJbnRlcnZhbFNlY2'
    '9uZHM=');

@$core.Deprecated('Use strategyConfigDescriptor instead')
const StrategyConfig$json = {
  '1': 'StrategyConfig',
  '2': [
    {
      '1': 'recreate',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.RecreateStrategy',
      '9': 0,
      '10': 'recreate'
    },
    {
      '1': 'rolling',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.RollingStrategy',
      '9': 0,
      '10': 'rolling'
    },
    {
      '1': 'blue_green',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.BlueGreenStrategy',
      '9': 0,
      '10': 'blueGreen'
    },
    {
      '1': 'canary',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.CanaryStrategy',
      '9': 0,
      '10': 'canary'
    },
  ],
  '8': [
    {'1': 'strategy'},
  ],
};

/// Descriptor for `StrategyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List strategyConfigDescriptor = $convert.base64Decode(
    'Cg5TdHJhdGVneUNvbmZpZxI7CghyZWNyZWF0ZRgBIAEoCzIdLmZvcmdlbmV0LnYxLlJlY3JlYX'
    'RlU3RyYXRlZ3lIAFIIcmVjcmVhdGUSOAoHcm9sbGluZxgCIAEoCzIcLmZvcmdlbmV0LnYxLlJv'
    'bGxpbmdTdHJhdGVneUgAUgdyb2xsaW5nEj8KCmJsdWVfZ3JlZW4YAyABKAsyHi5mb3JnZW5ldC'
    '52MS5CbHVlR3JlZW5TdHJhdGVneUgAUglibHVlR3JlZW4SNQoGY2FuYXJ5GAQgASgLMhsuZm9y'
    'Z2VuZXQudjEuQ2FuYXJ5U3RyYXRlZ3lIAFIGY2FuYXJ5QgoKCHN0cmF0ZWd5');

@$core.Deprecated('Use gitSourceDescriptor instead')
const GitSource$json = {
  '1': 'GitSource',
  '2': [
    {'1': 'repository_url', '3': 1, '4': 1, '5': 9, '10': 'repositoryUrl'},
    {'1': 'branch', '3': 2, '4': 1, '5': 9, '10': 'branch'},
    {'1': 'commit_sha', '3': 3, '4': 1, '5': 9, '10': 'commitSha'},
    {'1': 'dockerfile_path', '3': 4, '4': 1, '5': 9, '10': 'dockerfilePath'},
  ],
};

/// Descriptor for `GitSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List gitSourceDescriptor = $convert.base64Decode(
    'CglHaXRTb3VyY2USJQoOcmVwb3NpdG9yeV91cmwYASABKAlSDXJlcG9zaXRvcnlVcmwSFgoGYn'
    'JhbmNoGAIgASgJUgZicmFuY2gSHQoKY29tbWl0X3NoYRgDIAEoCVIJY29tbWl0U2hhEicKD2Rv'
    'Y2tlcmZpbGVfcGF0aBgEIAEoCVIOZG9ja2VyZmlsZVBhdGg=');

@$core.Deprecated('Use imageSourceDescriptor instead')
const ImageSource$json = {
  '1': 'ImageSource',
  '2': [
    {'1': 'image_ref', '3': 1, '4': 1, '5': 9, '10': 'imageRef'},
  ],
};

/// Descriptor for `ImageSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List imageSourceDescriptor = $convert
    .base64Decode('CgtJbWFnZVNvdXJjZRIbCglpbWFnZV9yZWYYASABKAlSCGltYWdlUmVm');

@$core.Deprecated('Use composeSourceDescriptor instead')
const ComposeSource$json = {
  '1': 'ComposeSource',
  '2': [
    {'1': 'compose_content', '3': 1, '4': 1, '5': 9, '10': 'composeContent'},
  ],
};

/// Descriptor for `ComposeSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List composeSourceDescriptor = $convert.base64Decode(
    'Cg1Db21wb3NlU291cmNlEicKD2NvbXBvc2VfY29udGVudBgBIAEoCVIOY29tcG9zZUNvbnRlbn'
    'Q=');

@$core.Deprecated('Use deploymentSourceDescriptor instead')
const DeploymentSource$json = {
  '1': 'DeploymentSource',
  '2': [
    {
      '1': 'git',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.GitSource',
      '9': 0,
      '10': 'git'
    },
    {
      '1': 'image',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.ImageSource',
      '9': 0,
      '10': 'image'
    },
    {
      '1': 'compose',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.ComposeSource',
      '9': 0,
      '10': 'compose'
    },
  ],
  '8': [
    {'1': 'source'},
  ],
};

/// Descriptor for `DeploymentSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentSourceDescriptor = $convert.base64Decode(
    'ChBEZXBsb3ltZW50U291cmNlEioKA2dpdBgBIAEoCzIWLmZvcmdlbmV0LnYxLkdpdFNvdXJjZU'
    'gAUgNnaXQSMAoFaW1hZ2UYAiABKAsyGC5mb3JnZW5ldC52MS5JbWFnZVNvdXJjZUgAUgVpbWFn'
    'ZRI2Cgdjb21wb3NlGAMgASgLMhouZm9yZ2VuZXQudjEuQ29tcG9zZVNvdXJjZUgAUgdjb21wb3'
    'NlQggKBnNvdXJjZQ==');

@$core.Deprecated('Use deploymentDescriptor instead')
const Deployment$json = {
  '1': 'Deployment',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'project_id', '3': 2, '4': 1, '5': 9, '10': 'projectId'},
    {
      '1': 'source',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.DeploymentSource',
      '10': 'source'
    },
    {
      '1': 'state',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.DeploymentState',
      '10': 'state'
    },
    {
      '1': 'strategy',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.StrategyConfig',
      '10': 'strategy'
    },
    {'1': 'error_message', '3': 6, '4': 1, '5': 9, '10': 'errorMessage'},
    {'1': 'target_env', '3': 7, '4': 1, '5': 9, '10': 'targetEnv'},
    {
      '1': 'previous_deployment_id',
      '3': 8,
      '4': 1,
      '5': 9,
      '10': 'previousDeploymentId'
    },
    {'1': 'triggered_by', '3': 9, '4': 1, '5': 9, '10': 'triggeredBy'},
    {
      '1': 'created_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    {
      '1': 'finished_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'finishedAt'
    },
  ],
};

/// Descriptor for `Deployment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentDescriptor = $convert.base64Decode(
    'CgpEZXBsb3ltZW50Eg4KAmlkGAEgASgJUgJpZBIdCgpwcm9qZWN0X2lkGAIgASgJUglwcm9qZW'
    'N0SWQSNQoGc291cmNlGAMgASgLMh0uZm9yZ2VuZXQudjEuRGVwbG95bWVudFNvdXJjZVIGc291'
    'cmNlEjIKBXN0YXRlGAQgASgOMhwuZm9yZ2VuZXQudjEuRGVwbG95bWVudFN0YXRlUgVzdGF0ZR'
    'I3CghzdHJhdGVneRgFIAEoCzIbLmZvcmdlbmV0LnYxLlN0cmF0ZWd5Q29uZmlnUghzdHJhdGVn'
    'eRIjCg1lcnJvcl9tZXNzYWdlGAYgASgJUgxlcnJvck1lc3NhZ2USHQoKdGFyZ2V0X2VudhgHIA'
    'EoCVIJdGFyZ2V0RW52EjQKFnByZXZpb3VzX2RlcGxveW1lbnRfaWQYCCABKAlSFHByZXZpb3Vz'
    'RGVwbG95bWVudElkEiEKDHRyaWdnZXJlZF9ieRgJIAEoCVILdHJpZ2dlcmVkQnkSOQoKY3JlYX'
    'RlZF9hdBgKIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdBI5Cgp1'
    'cGRhdGVkX2F0GAsgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdXBkYXRlZEF0Ej'
    'sKC2ZpbmlzaGVkX2F0GAwgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIKZmluaXNo'
    'ZWRBdA==');

@$core.Deprecated('Use deploymentLogEntryDescriptor instead')
const DeploymentLogEntry$json = {
  '1': 'DeploymentLogEntry',
  '2': [
    {
      '1': 'emitted_at',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'emittedAt'
    },
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {'1': 'step', '3': 3, '4': 1, '5': 9, '10': 'step'},
    {
      '1': 'stream',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.LogStream',
      '10': 'stream'
    },
  ],
};

/// Descriptor for `DeploymentLogEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentLogEntryDescriptor = $convert.base64Decode(
    'ChJEZXBsb3ltZW50TG9nRW50cnkSOQoKZW1pdHRlZF9hdBgBIAEoCzIaLmdvb2dsZS5wcm90b2'
    'J1Zi5UaW1lc3RhbXBSCWVtaXR0ZWRBdBIYCgdtZXNzYWdlGAIgASgJUgdtZXNzYWdlEhIKBHN0'
    'ZXAYAyABKAlSBHN0ZXASLgoGc3RyZWFtGAQgASgOMhYuZm9yZ2VuZXQudjEuTG9nU3RyZWFtUg'
    'ZzdHJlYW0=');

@$core.Deprecated('Use createDeploymentRequestDescriptor instead')
const CreateDeploymentRequest$json = {
  '1': 'CreateDeploymentRequest',
  '2': [
    {'1': 'project_id', '3': 1, '4': 1, '5': 9, '10': 'projectId'},
    {
      '1': 'source',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.DeploymentSource',
      '10': 'source'
    },
    {
      '1': 'strategy',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.StrategyConfig',
      '10': 'strategy'
    },
    {'1': 'target_env', '3': 4, '4': 1, '5': 9, '10': 'targetEnv'},
    {'1': 'request_id', '3': 5, '4': 1, '5': 9, '10': 'requestId'},
  ],
};

/// Descriptor for `CreateDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeploymentRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVEZXBsb3ltZW50UmVxdWVzdBIdCgpwcm9qZWN0X2lkGAEgASgJUglwcm9qZWN0SW'
    'QSNQoGc291cmNlGAIgASgLMh0uZm9yZ2VuZXQudjEuRGVwbG95bWVudFNvdXJjZVIGc291cmNl'
    'EjcKCHN0cmF0ZWd5GAMgASgLMhsuZm9yZ2VuZXQudjEuU3RyYXRlZ3lDb25maWdSCHN0cmF0ZW'
    'd5Eh0KCnRhcmdldF9lbnYYBCABKAlSCXRhcmdldEVudhIdCgpyZXF1ZXN0X2lkGAUgASgJUgly'
    'ZXF1ZXN0SWQ=');

@$core.Deprecated('Use createDeploymentResponseDescriptor instead')
const CreateDeploymentResponse$json = {
  '1': 'CreateDeploymentResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
  ],
};

/// Descriptor for `CreateDeploymentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeploymentResponseDescriptor =
    $convert.base64Decode(
        'ChhDcmVhdGVEZXBsb3ltZW50UmVzcG9uc2USNwoKZGVwbG95bWVudBgBIAEoCzIXLmZvcmdlbm'
        'V0LnYxLkRlcGxveW1lbnRSCmRlcGxveW1lbnQ=');

@$core.Deprecated('Use getDeploymentRequestDescriptor instead')
const GetDeploymentRequest$json = {
  '1': 'GetDeploymentRequest',
  '2': [
    {'1': 'deployment_id', '3': 1, '4': 1, '5': 9, '10': 'deploymentId'},
  ],
};

/// Descriptor for `GetDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeploymentRequestDescriptor = $convert.base64Decode(
    'ChRHZXREZXBsb3ltZW50UmVxdWVzdBIjCg1kZXBsb3ltZW50X2lkGAEgASgJUgxkZXBsb3ltZW'
    '50SWQ=');

@$core.Deprecated('Use getDeploymentResponseDescriptor instead')
const GetDeploymentResponse$json = {
  '1': 'GetDeploymentResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
  ],
};

/// Descriptor for `GetDeploymentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeploymentResponseDescriptor = $convert.base64Decode(
    'ChVHZXREZXBsb3ltZW50UmVzcG9uc2USNwoKZGVwbG95bWVudBgBIAEoCzIXLmZvcmdlbmV0Ln'
    'YxLkRlcGxveW1lbnRSCmRlcGxveW1lbnQ=');

@$core.Deprecated('Use listDeploymentsRequestDescriptor instead')
const ListDeploymentsRequest$json = {
  '1': 'ListDeploymentsRequest',
  '2': [
    {'1': 'project_id', '3': 1, '4': 1, '5': 9, '10': 'projectId'},
    {'1': 'page_size', '3': 2, '4': 1, '5': 5, '10': 'pageSize'},
    {'1': 'page_token', '3': 3, '4': 1, '5': 9, '10': 'pageToken'},
  ],
};

/// Descriptor for `ListDeploymentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeploymentsRequestDescriptor = $convert.base64Decode(
    'ChZMaXN0RGVwbG95bWVudHNSZXF1ZXN0Eh0KCnByb2plY3RfaWQYASABKAlSCXByb2plY3RJZB'
    'IbCglwYWdlX3NpemUYAiABKAVSCHBhZ2VTaXplEh0KCnBhZ2VfdG9rZW4YAyABKAlSCXBhZ2VU'
    'b2tlbg==');

@$core.Deprecated('Use listDeploymentsResponseDescriptor instead')
const ListDeploymentsResponse$json = {
  '1': 'ListDeploymentsResponse',
  '2': [
    {
      '1': 'deployments',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployments'
    },
    {'1': 'next_page_token', '3': 2, '4': 1, '5': 9, '10': 'nextPageToken'},
  ],
};

/// Descriptor for `ListDeploymentsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeploymentsResponseDescriptor = $convert.base64Decode(
    'ChdMaXN0RGVwbG95bWVudHNSZXNwb25zZRI5CgtkZXBsb3ltZW50cxgBIAMoCzIXLmZvcmdlbm'
    'V0LnYxLkRlcGxveW1lbnRSC2RlcGxveW1lbnRzEiYKD25leHRfcGFnZV90b2tlbhgCIAEoCVIN'
    'bmV4dFBhZ2VUb2tlbg==');

@$core.Deprecated('Use cancelDeploymentRequestDescriptor instead')
const CancelDeploymentRequest$json = {
  '1': 'CancelDeploymentRequest',
  '2': [
    {'1': 'deployment_id', '3': 1, '4': 1, '5': 9, '10': 'deploymentId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `CancelDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelDeploymentRequestDescriptor =
    $convert.base64Decode(
        'ChdDYW5jZWxEZXBsb3ltZW50UmVxdWVzdBIjCg1kZXBsb3ltZW50X2lkGAEgASgJUgxkZXBsb3'
        'ltZW50SWQSFgoGcmVhc29uGAIgASgJUgZyZWFzb24=');

@$core.Deprecated('Use cancelDeploymentResponseDescriptor instead')
const CancelDeploymentResponse$json = {
  '1': 'CancelDeploymentResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
  ],
};

/// Descriptor for `CancelDeploymentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelDeploymentResponseDescriptor =
    $convert.base64Decode(
        'ChhDYW5jZWxEZXBsb3ltZW50UmVzcG9uc2USNwoKZGVwbG95bWVudBgBIAEoCzIXLmZvcmdlbm'
        'V0LnYxLkRlcGxveW1lbnRSCmRlcGxveW1lbnQ=');

@$core.Deprecated('Use rollbackDeploymentRequestDescriptor instead')
const RollbackDeploymentRequest$json = {
  '1': 'RollbackDeploymentRequest',
  '2': [
    {
      '1': 'target_deployment_id',
      '3': 1,
      '4': 1,
      '5': 9,
      '10': 'targetDeploymentId'
    },
    {'1': 'request_id', '3': 2, '4': 1, '5': 9, '10': 'requestId'},
  ],
};

/// Descriptor for `RollbackDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rollbackDeploymentRequestDescriptor =
    $convert.base64Decode(
        'ChlSb2xsYmFja0RlcGxveW1lbnRSZXF1ZXN0EjAKFHRhcmdldF9kZXBsb3ltZW50X2lkGAEgAS'
        'gJUhJ0YXJnZXREZXBsb3ltZW50SWQSHQoKcmVxdWVzdF9pZBgCIAEoCVIJcmVxdWVzdElk');

@$core.Deprecated('Use rollbackDeploymentResponseDescriptor instead')
const RollbackDeploymentResponse$json = {
  '1': 'RollbackDeploymentResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
  ],
};

/// Descriptor for `RollbackDeploymentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rollbackDeploymentResponseDescriptor =
    $convert.base64Decode(
        'ChpSb2xsYmFja0RlcGxveW1lbnRSZXNwb25zZRI3CgpkZXBsb3ltZW50GAEgASgLMhcuZm9yZ2'
        'VuZXQudjEuRGVwbG95bWVudFIKZGVwbG95bWVudA==');

@$core.Deprecated('Use streamDeploymentLogsRequestDescriptor instead')
const StreamDeploymentLogsRequest$json = {
  '1': 'StreamDeploymentLogsRequest',
  '2': [
    {'1': 'deployment_id', '3': 1, '4': 1, '5': 9, '10': 'deploymentId'},
    {'1': 'follow', '3': 2, '4': 1, '5': 8, '10': 'follow'},
    {
      '1': 'since',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'since'
    },
    {
      '1': 'stream',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.LogStream',
      '10': 'stream'
    },
  ],
};

/// Descriptor for `StreamDeploymentLogsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDeploymentLogsRequestDescriptor = $convert.base64Decode(
    'ChtTdHJlYW1EZXBsb3ltZW50TG9nc1JlcXVlc3QSIwoNZGVwbG95bWVudF9pZBgBIAEoCVIMZG'
    'VwbG95bWVudElkEhYKBmZvbGxvdxgCIAEoCFIGZm9sbG93EjAKBXNpbmNlGAMgASgLMhouZ29v'
    'Z2xlLnByb3RvYnVmLlRpbWVzdGFtcFIFc2luY2USLgoGc3RyZWFtGAQgASgOMhYuZm9yZ2VuZX'
    'QudjEuTG9nU3RyZWFtUgZzdHJlYW0=');

@$core.Deprecated('Use streamDeploymentLogsResponseDescriptor instead')
const StreamDeploymentLogsResponse$json = {
  '1': 'StreamDeploymentLogsResponse',
  '2': [
    {
      '1': 'log_entry',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.DeploymentLogEntry',
      '10': 'logEntry'
    },
  ],
};

/// Descriptor for `StreamDeploymentLogsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDeploymentLogsResponseDescriptor =
    $convert.base64Decode(
        'ChxTdHJlYW1EZXBsb3ltZW50TG9nc1Jlc3BvbnNlEjwKCWxvZ19lbnRyeRgBIAEoCzIfLmZvcm'
        'dlbmV0LnYxLkRlcGxveW1lbnRMb2dFbnRyeVIIbG9nRW50cnk=');

@$core.Deprecated('Use streamDeploymentRequestDescriptor instead')
const StreamDeploymentRequest$json = {
  '1': 'StreamDeploymentRequest',
  '2': [
    {'1': 'deployment_id', '3': 1, '4': 1, '5': 9, '10': 'deploymentId'},
  ],
};

/// Descriptor for `StreamDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDeploymentRequestDescriptor =
    $convert.base64Decode(
        'ChdTdHJlYW1EZXBsb3ltZW50UmVxdWVzdBIjCg1kZXBsb3ltZW50X2lkGAEgASgJUgxkZXBsb3'
        'ltZW50SWQ=');

@$core.Deprecated('Use streamDeploymentResponseDescriptor instead')
const StreamDeploymentResponse$json = {
  '1': 'StreamDeploymentResponse',
  '2': [
    {
      '1': 'deployment',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Deployment',
      '10': 'deployment'
    },
  ],
};

/// Descriptor for `StreamDeploymentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDeploymentResponseDescriptor =
    $convert.base64Decode(
        'ChhTdHJlYW1EZXBsb3ltZW50UmVzcG9uc2USNwoKZGVwbG95bWVudBgBIAEoCzIXLmZvcmdlbm'
        'V0LnYxLkRlcGxveW1lbnRSCmRlcGxveW1lbnQ=');
