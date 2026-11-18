// This is a generated file - do not edit.
//
// Generated from forgenet/v1/mail.proto.

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

@$core.Deprecated('Use mailTLSDescriptor instead')
const MailTLS$json = {
  '1': 'MailTLS',
  '2': [
    {'1': 'MAIL_TLS_UNSPECIFIED', '2': 0},
    {'1': 'MAIL_TLS_STARTTLS', '2': 1},
    {'1': 'MAIL_TLS_IMPLICIT', '2': 2},
    {'1': 'MAIL_TLS_NONE', '2': 3},
  ],
};

/// Descriptor for `MailTLS`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mailTLSDescriptor = $convert.base64Decode(
    'CgdNYWlsVExTEhgKFE1BSUxfVExTX1VOU1BFQ0lGSUVEEAASFQoRTUFJTF9UTFNfU1RBUlRUTF'
    'MQARIVChFNQUlMX1RMU19JTVBMSUNJVBACEhEKDU1BSUxfVExTX05PTkUQAw==');

@$core.Deprecated('Use getMailStatusRequestDescriptor instead')
const GetMailStatusRequest$json = {
  '1': 'GetMailStatusRequest',
};

/// Descriptor for `GetMailStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMailStatusRequestDescriptor =
    $convert.base64Decode('ChRHZXRNYWlsU3RhdHVzUmVxdWVzdA==');

@$core.Deprecated('Use getMailStatusResponseDescriptor instead')
const GetMailStatusResponse$json = {
  '1': 'GetMailStatusResponse',
  '2': [
    {'1': 'enabled', '3': 1, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'host', '3': 2, '4': 1, '5': 9, '10': 'host'},
    {'1': 'port', '3': 3, '4': 1, '5': 13, '10': 'port'},
    {
      '1': 'tls',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MailTLS',
      '10': 'tls'
    },
    {'1': 'from', '3': 5, '4': 1, '5': 9, '10': 'from'},
    {'1': 'authenticated', '3': 6, '4': 1, '5': 8, '10': 'authenticated'},
  ],
};

/// Descriptor for `GetMailStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMailStatusResponseDescriptor = $convert.base64Decode(
    'ChVHZXRNYWlsU3RhdHVzUmVzcG9uc2USGAoHZW5hYmxlZBgBIAEoCFIHZW5hYmxlZBISCgRob3'
    'N0GAIgASgJUgRob3N0EhIKBHBvcnQYAyABKA1SBHBvcnQSJgoDdGxzGAQgASgOMhQuZm9yZ2Vu'
    'ZXQudjEuTWFpbFRMU1IDdGxzEhIKBGZyb20YBSABKAlSBGZyb20SJAoNYXV0aGVudGljYXRlZB'
    'gGIAEoCFINYXV0aGVudGljYXRlZA==');

@$core.Deprecated('Use sendTestMailRequestDescriptor instead')
const SendTestMailRequest$json = {
  '1': 'SendTestMailRequest',
};

/// Descriptor for `SendTestMailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTestMailRequestDescriptor =
    $convert.base64Decode('ChNTZW5kVGVzdE1haWxSZXF1ZXN0');

@$core.Deprecated('Use sendTestMailResponseDescriptor instead')
const SendTestMailResponse$json = {
  '1': 'SendTestMailResponse',
  '2': [
    {'1': 'to', '3': 1, '4': 1, '5': 9, '10': 'to'},
  ],
};

/// Descriptor for `SendTestMailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTestMailResponseDescriptor = $convert
    .base64Decode('ChRTZW5kVGVzdE1haWxSZXNwb25zZRIOCgJ0bxgBIAEoCVICdG8=');
