// This is a generated file - do not edit.
//
// Generated from forgenet/v1/team.proto.

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

@$core.Deprecated('Use memberStatusDescriptor instead')
const MemberStatus$json = {
  '1': 'MemberStatus',
  '2': [
    {'1': 'MEMBER_STATUS_UNSPECIFIED', '2': 0},
    {'1': 'MEMBER_STATUS_ACTIVE', '2': 1},
    {'1': 'MEMBER_STATUS_INVITED', '2': 2},
    {'1': 'MEMBER_STATUS_INVITATION_EXPIRED', '2': 3},
  ],
};

/// Descriptor for `MemberStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List memberStatusDescriptor = $convert.base64Decode(
    'CgxNZW1iZXJTdGF0dXMSHQoZTUVNQkVSX1NUQVRVU19VTlNQRUNJRklFRBAAEhgKFE1FTUJFUl'
    '9TVEFUVVNfQUNUSVZFEAESGQoVTUVNQkVSX1NUQVRVU19JTlZJVEVEEAISJAogTUVNQkVSX1NU'
    'QVRVU19JTlZJVEFUSU9OX0VYUElSRUQQAw==');

@$core.Deprecated('Use memberDescriptor instead')
const Member$json = {
  '1': 'Member',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'username', '3': 2, '4': 1, '5': 9, '10': 'username'},
    {'1': 'email', '3': 3, '4': 1, '5': 9, '10': 'email'},
    {
      '1': 'role',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.UserRole',
      '10': 'role'
    },
    {
      '1': 'status',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MemberStatus',
      '10': 'status'
    },
    {
      '1': 'invitation_expires_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'invitationExpiresAt'
    },
    {
      '1': 'created_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Member`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List memberDescriptor = $convert.base64Decode(
    'CgZNZW1iZXISDgoCaWQYASABKAlSAmlkEhoKCHVzZXJuYW1lGAIgASgJUgh1c2VybmFtZRIUCg'
    'VlbWFpbBgDIAEoCVIFZW1haWwSKQoEcm9sZRgEIAEoDjIVLmZvcmdlbmV0LnYxLlVzZXJSb2xl'
    'UgRyb2xlEjEKBnN0YXR1cxgFIAEoDjIZLmZvcmdlbmV0LnYxLk1lbWJlclN0YXR1c1IGc3RhdH'
    'VzEk4KFWludml0YXRpb25fZXhwaXJlc19hdBgGIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1l'
    'c3RhbXBSE2ludml0YXRpb25FeHBpcmVzQXQSOQoKY3JlYXRlZF9hdBgHIAEoCzIaLmdvb2dsZS'
    '5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use listMembersRequestDescriptor instead')
const ListMembersRequest$json = {
  '1': 'ListMembersRequest',
};

/// Descriptor for `ListMembersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMembersRequestDescriptor =
    $convert.base64Decode('ChJMaXN0TWVtYmVyc1JlcXVlc3Q=');

@$core.Deprecated('Use listMembersResponseDescriptor instead')
const ListMembersResponse$json = {
  '1': 'ListMembersResponse',
  '2': [
    {
      '1': 'members',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.Member',
      '10': 'members'
    },
    {'1': 'mail_enabled', '3': 2, '4': 1, '5': 8, '10': 'mailEnabled'},
  ],
};

/// Descriptor for `ListMembersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMembersResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0TWVtYmVyc1Jlc3BvbnNlEi0KB21lbWJlcnMYASADKAsyEy5mb3JnZW5ldC52MS5NZW'
    '1iZXJSB21lbWJlcnMSIQoMbWFpbF9lbmFibGVkGAIgASgIUgttYWlsRW5hYmxlZA==');

@$core.Deprecated('Use inviteMemberRequestDescriptor instead')
const InviteMemberRequest$json = {
  '1': 'InviteMemberRequest',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {
      '1': 'role',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.UserRole',
      '10': 'role'
    },
  ],
};

/// Descriptor for `InviteMemberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inviteMemberRequestDescriptor = $convert.base64Decode(
    'ChNJbnZpdGVNZW1iZXJSZXF1ZXN0EhoKCHVzZXJuYW1lGAEgASgJUgh1c2VybmFtZRIUCgVlbW'
    'FpbBgCIAEoCVIFZW1haWwSKQoEcm9sZRgDIAEoDjIVLmZvcmdlbmV0LnYxLlVzZXJSb2xlUgRy'
    'b2xl');

@$core.Deprecated('Use invitationResultDescriptor instead')
const InvitationResult$json = {
  '1': 'InvitationResult',
  '2': [
    {
      '1': 'member',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Member',
      '10': 'member'
    },
    {'1': 'emailed', '3': 2, '4': 1, '5': 8, '10': 'emailed'},
    {'1': 'invitation_url', '3': 3, '4': 1, '5': 9, '10': 'invitationUrl'},
    {'1': 'local_dashboard', '3': 4, '4': 1, '5': 8, '10': 'localDashboard'},
  ],
};

/// Descriptor for `InvitationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invitationResultDescriptor = $convert.base64Decode(
    'ChBJbnZpdGF0aW9uUmVzdWx0EisKBm1lbWJlchgBIAEoCzITLmZvcmdlbmV0LnYxLk1lbWJlcl'
    'IGbWVtYmVyEhgKB2VtYWlsZWQYAiABKAhSB2VtYWlsZWQSJQoOaW52aXRhdGlvbl91cmwYAyAB'
    'KAlSDWludml0YXRpb25VcmwSJwoPbG9jYWxfZGFzaGJvYXJkGAQgASgIUg5sb2NhbERhc2hib2'
    'FyZA==');

@$core.Deprecated('Use inviteMemberResponseDescriptor instead')
const InviteMemberResponse$json = {
  '1': 'InviteMemberResponse',
  '2': [
    {
      '1': 'result',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.InvitationResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `InviteMemberResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inviteMemberResponseDescriptor = $convert.base64Decode(
    'ChRJbnZpdGVNZW1iZXJSZXNwb25zZRI1CgZyZXN1bHQYASABKAsyHS5mb3JnZW5ldC52MS5Jbn'
    'ZpdGF0aW9uUmVzdWx0UgZyZXN1bHQ=');

@$core.Deprecated('Use renewInvitationRequestDescriptor instead')
const RenewInvitationRequest$json = {
  '1': 'RenewInvitationRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `RenewInvitationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renewInvitationRequestDescriptor =
    $convert.base64Decode(
        'ChZSZW5ld0ludml0YXRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use renewInvitationResponseDescriptor instead')
const RenewInvitationResponse$json = {
  '1': 'RenewInvitationResponse',
  '2': [
    {
      '1': 'result',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.InvitationResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `RenewInvitationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renewInvitationResponseDescriptor =
    $convert.base64Decode(
        'ChdSZW5ld0ludml0YXRpb25SZXNwb25zZRI1CgZyZXN1bHQYASABKAsyHS5mb3JnZW5ldC52MS'
        '5JbnZpdGF0aW9uUmVzdWx0UgZyZXN1bHQ=');
