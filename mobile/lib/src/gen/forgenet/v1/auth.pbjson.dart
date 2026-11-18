// This is a generated file - do not edit.
//
// Generated from forgenet/v1/auth.proto.

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

@$core.Deprecated('Use userRoleDescriptor instead')
const UserRole$json = {
  '1': 'UserRole',
  '2': [
    {'1': 'USER_ROLE_UNSPECIFIED', '2': 0},
    {'1': 'USER_ROLE_ADMIN', '2': 1},
    {'1': 'USER_ROLE_MEMBER', '2': 2},
  ],
};

/// Descriptor for `UserRole`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userRoleDescriptor = $convert.base64Decode(
    'CghVc2VyUm9sZRIZChVVU0VSX1JPTEVfVU5TUEVDSUZJRUQQABITCg9VU0VSX1JPTEVfQURNSU'
    '4QARIUChBVU0VSX1JPTEVfTUVNQkVSEAI=');

@$core.Deprecated('Use accessLinkPurposeDescriptor instead')
const AccessLinkPurpose$json = {
  '1': 'AccessLinkPurpose',
  '2': [
    {'1': 'ACCESS_LINK_PURPOSE_UNSPECIFIED', '2': 0},
    {'1': 'ACCESS_LINK_PURPOSE_FIRST_LOGIN', '2': 1},
    {'1': 'ACCESS_LINK_PURPOSE_PASSWORD_RESET', '2': 2},
    {'1': 'ACCESS_LINK_PURPOSE_INVITATION', '2': 3},
  ],
};

/// Descriptor for `AccessLinkPurpose`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accessLinkPurposeDescriptor = $convert.base64Decode(
    'ChFBY2Nlc3NMaW5rUHVycG9zZRIjCh9BQ0NFU1NfTElOS19QVVJQT1NFX1VOU1BFQ0lGSUVEEA'
    'ASIwofQUNDRVNTX0xJTktfUFVSUE9TRV9GSVJTVF9MT0dJThABEiYKIkFDQ0VTU19MSU5LX1BV'
    'UlBPU0VfUEFTU1dPUkRfUkVTRVQQAhIiCh5BQ0NFU1NfTElOS19QVVJQT1NFX0lOVklUQVRJT0'
    '4QAw==');

@$core.Deprecated('Use userDescriptor instead')
const User$json = {
  '1': 'User',
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
  ],
};

/// Descriptor for `User`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userDescriptor = $convert.base64Decode(
    'CgRVc2VyEg4KAmlkGAEgASgJUgJpZBIaCgh1c2VybmFtZRgCIAEoCVIIdXNlcm5hbWUSFAoFZW'
    '1haWwYAyABKAlSBWVtYWlsEikKBHJvbGUYBCABKA4yFS5mb3JnZW5ldC52MS5Vc2VyUm9sZVIE'
    'cm9sZQ==');

@$core.Deprecated('Use loginRequestDescriptor instead')
const LoginRequest$json = {
  '1': 'LoginRequest',
  '2': [
    {'1': 'login', '3': 1, '4': 1, '5': 9, '10': 'login'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `LoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginRequestDescriptor = $convert.base64Decode(
    'CgxMb2dpblJlcXVlc3QSFAoFbG9naW4YASABKAlSBWxvZ2luEhoKCHBhc3N3b3JkGAIgASgJUg'
    'hwYXNzd29yZA==');

@$core.Deprecated('Use loginResponseDescriptor instead')
const LoginResponse$json = {
  '1': 'LoginResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'expires_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {
      '1': 'user',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `LoginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginResponseDescriptor = $convert.base64Decode(
    'Cg1Mb2dpblJlc3BvbnNlEhQKBXRva2VuGAEgASgJUgV0b2tlbhI5CgpleHBpcmVzX2F0GAIgAS'
    'gLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJZXhwaXJlc0F0EiUKBHVzZXIYAyABKAsy'
    'ES5mb3JnZW5ldC52MS5Vc2VyUgR1c2Vy');

@$core.Deprecated('Use logoutRequestDescriptor instead')
const LogoutRequest$json = {
  '1': 'LogoutRequest',
};

/// Descriptor for `LogoutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutRequestDescriptor =
    $convert.base64Decode('Cg1Mb2dvdXRSZXF1ZXN0');

@$core.Deprecated('Use logoutResponseDescriptor instead')
const LogoutResponse$json = {
  '1': 'LogoutResponse',
};

/// Descriptor for `LogoutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutResponseDescriptor =
    $convert.base64Decode('Cg5Mb2dvdXRSZXNwb25zZQ==');

@$core.Deprecated('Use getCurrentUserRequestDescriptor instead')
const GetCurrentUserRequest$json = {
  '1': 'GetCurrentUserRequest',
};

/// Descriptor for `GetCurrentUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCurrentUserRequestDescriptor =
    $convert.base64Decode('ChVHZXRDdXJyZW50VXNlclJlcXVlc3Q=');

@$core.Deprecated('Use getCurrentUserResponseDescriptor instead')
const GetCurrentUserResponse$json = {
  '1': 'GetCurrentUserResponse',
  '2': [
    {
      '1': 'user',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `GetCurrentUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCurrentUserResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRDdXJyZW50VXNlclJlc3BvbnNlEiUKBHVzZXIYASABKAsyES5mb3JnZW5ldC52MS5Vc2'
        'VyUgR1c2Vy');

@$core.Deprecated('Use changePasswordRequestDescriptor instead')
const ChangePasswordRequest$json = {
  '1': 'ChangePasswordRequest',
  '2': [
    {'1': 'current_password', '3': 1, '4': 1, '5': 9, '10': 'currentPassword'},
    {'1': 'new_password', '3': 2, '4': 1, '5': 9, '10': 'newPassword'},
  ],
};

/// Descriptor for `ChangePasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changePasswordRequestDescriptor = $convert.base64Decode(
    'ChVDaGFuZ2VQYXNzd29yZFJlcXVlc3QSKQoQY3VycmVudF9wYXNzd29yZBgBIAEoCVIPY3Vycm'
    'VudFBhc3N3b3JkEiEKDG5ld19wYXNzd29yZBgCIAEoCVILbmV3UGFzc3dvcmQ=');

@$core.Deprecated('Use changePasswordResponseDescriptor instead')
const ChangePasswordResponse$json = {
  '1': 'ChangePasswordResponse',
  '2': [
    {'1': 'revoked_sessions', '3': 1, '4': 1, '5': 5, '10': 'revokedSessions'},
  ],
};

/// Descriptor for `ChangePasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changePasswordResponseDescriptor =
    $convert.base64Decode(
        'ChZDaGFuZ2VQYXNzd29yZFJlc3BvbnNlEikKEHJldm9rZWRfc2Vzc2lvbnMYASABKAVSD3Jldm'
        '9rZWRTZXNzaW9ucw==');

@$core.Deprecated('Use inspectAccessLinkRequestDescriptor instead')
const InspectAccessLinkRequest$json = {
  '1': 'InspectAccessLinkRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `InspectAccessLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectAccessLinkRequestDescriptor =
    $convert.base64Decode(
        'ChhJbnNwZWN0QWNjZXNzTGlua1JlcXVlc3QSFAoFdG9rZW4YASABKAlSBXRva2Vu');

@$core.Deprecated('Use inspectAccessLinkResponseDescriptor instead')
const InspectAccessLinkResponse$json = {
  '1': 'InspectAccessLinkResponse',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {
      '1': 'purpose',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.AccessLinkPurpose',
      '10': 'purpose'
    },
    {
      '1': 'expires_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
  ],
};

/// Descriptor for `InspectAccessLinkResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectAccessLinkResponseDescriptor = $convert.base64Decode(
    'ChlJbnNwZWN0QWNjZXNzTGlua1Jlc3BvbnNlEhoKCHVzZXJuYW1lGAEgASgJUgh1c2VybmFtZR'
    'I4CgdwdXJwb3NlGAIgASgOMh4uZm9yZ2VuZXQudjEuQWNjZXNzTGlua1B1cnBvc2VSB3B1cnBv'
    'c2USOQoKZXhwaXJlc19hdBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWV4cG'
    'lyZXNBdA==');

@$core.Deprecated('Use redeemAccessLinkRequestDescriptor instead')
const RedeemAccessLinkRequest$json = {
  '1': 'RedeemAccessLinkRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {'1': 'new_password', '3': 2, '4': 1, '5': 9, '10': 'newPassword'},
  ],
};

/// Descriptor for `RedeemAccessLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemAccessLinkRequestDescriptor =
    $convert.base64Decode(
        'ChdSZWRlZW1BY2Nlc3NMaW5rUmVxdWVzdBIUCgV0b2tlbhgBIAEoCVIFdG9rZW4SIQoMbmV3X3'
        'Bhc3N3b3JkGAIgASgJUgtuZXdQYXNzd29yZA==');

@$core.Deprecated('Use redeemAccessLinkResponseDescriptor instead')
const RedeemAccessLinkResponse$json = {
  '1': 'RedeemAccessLinkResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'expires_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {
      '1': 'user',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `RedeemAccessLinkResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redeemAccessLinkResponseDescriptor = $convert.base64Decode(
    'ChhSZWRlZW1BY2Nlc3NMaW5rUmVzcG9uc2USFAoFdG9rZW4YASABKAlSBXRva2VuEjkKCmV4cG'
    'lyZXNfYXQYAiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUglleHBpcmVzQXQSJQoE'
    'dXNlchgDIAEoCzIRLmZvcmdlbmV0LnYxLlVzZXJSBHVzZXI=');
