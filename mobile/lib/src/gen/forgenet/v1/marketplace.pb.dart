// This is a generated file - do not edit.
//
// Generated from forgenet/v1/marketplace.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// TemplateVariable est une valeur demandée à l'installation, telle que la
/// déclare le manifeste du template (templates/README.md).
class TemplateVariable extends $pb.GeneratedMessage {
  factory TemplateVariable({
    $core.String? name,
    $core.String? label,
    $core.String? help,
    $core.String? defaultValue,
    $core.String? pattern,
    $core.bool? optional,
    $core.bool? secret,
    $core.int? generate,
  }) {
    final result = TemplateVariable._();
    if (name != null) result.name = name;
    if (label != null) result.label = label;
    if (help != null) result.help = help;
    if (defaultValue != null) result.defaultValue = defaultValue;
    if (pattern != null) result.pattern = pattern;
    if (optional != null) result.optional = optional;
    if (secret != null) result.secret = secret;
    if (generate != null) result.generate = generate;
    return result;
  }

  TemplateVariable._();

  factory TemplateVariable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplateVariable()..mergeFromBuffer(data, registry);
  factory TemplateVariable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplateVariable()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TemplateVariable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: TemplateVariable.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'label')
    ..aOS(3, _omitFieldNames ? '' : 'help')
    ..aOS(4, _omitFieldNames ? '' : 'defaultValue')
    ..aOS(5, _omitFieldNames ? '' : 'pattern')
    ..aOB(6, _omitFieldNames ? '' : 'optional')
    ..aOB(7, _omitFieldNames ? '' : 'secret')
    ..aI(8, _omitFieldNames ? '' : 'generate', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplateVariable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplateVariable copyWith(void Function(TemplateVariable) updates) =>
      super.copyWith((message) => updates(message as TemplateVariable))
          as TemplateVariable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TemplateVariable() / TemplateVariable.new instead')
  static TemplateVariable create() => TemplateVariable._();
  static $pb.GeneratedMessage $_createMessage() => TemplateVariable._();
  @$core.override
  TemplateVariable createEmptyInstance() => TemplateVariable._();
  @$core.pragma('dart2js:noInline')
  static TemplateVariable getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TemplateVariable>(
          TemplateVariable.$_createMessage);
  static TemplateVariable? _defaultInstance;

  /// Nom de la variable dans le fichier Compose, ex: "POSTGRES_PASSWORD".
  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  /// Libellé du champ du formulaire d'installation.
  @$pb.TagNumber(2)
  $core.String get label => $_getSZ(1);
  @$pb.TagNumber(2)
  set label($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLabel() => $_has(1);
  @$pb.TagNumber(2)
  void clearLabel() => $_clearField(2);

  /// Phrase d'aide ; vide si le manifeste n'en donne pas.
  @$pb.TagNumber(3)
  $core.String get help => $_getSZ(2);
  @$pb.TagNumber(3)
  set help($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasHelp() => $_has(2);
  @$pb.TagNumber(3)
  void clearHelp() => $_clearField(3);

  /// Valeur proposée ; toujours vide pour un secret.
  @$pb.TagNumber(4)
  $core.String get defaultValue => $_getSZ(3);
  @$pb.TagNumber(4)
  set defaultValue($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDefaultValue() => $_has(3);
  @$pb.TagNumber(4)
  void clearDefaultValue() => $_clearField(4);

  /// Expression régulière que la valeur respecte ; vide si aucune.
  @$pb.TagNumber(5)
  $core.String get pattern => $_getSZ(4);
  @$pb.TagNumber(5)
  set pattern($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPattern() => $_has(4);
  @$pb.TagNumber(5)
  void clearPattern() => $_clearField(5);

  /// Vrai : une valeur vide est acceptée.
  @$pb.TagNumber(6)
  $core.bool get optional => $_getBF(5);
  @$pb.TagNumber(6)
  set optional($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasOptional() => $_has(5);
  @$pb.TagNumber(6)
  void clearOptional() => $_clearField(6);

  /// Vrai : la valeur est masquée, chiffrée sur le serveur et jamais
  /// réaffichée.
  @$pb.TagNumber(7)
  $core.bool get secret => $_getBF(6);
  @$pb.TagNumber(7)
  set secret($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasSecret() => $_has(6);
  @$pb.TagNumber(7)
  void clearSecret() => $_clearField(7);

  /// Pour un secret : longueur de la valeur tirée au hasard si
  /// l'administrateur ne saisit rien. 0 : la valeur doit être saisie.
  @$pb.TagNumber(8)
  $core.int get generate => $_getIZ(7);
  @$pb.TagNumber(8)
  set generate($core.int value) => $_setUnsignedInt32(7, value);
  @$pb.TagNumber(8)
  $core.bool hasGenerate() => $_has(7);
  @$pb.TagNumber(8)
  void clearGenerate() => $_clearField(8);
}

/// TemplatePort est un port où un service du template écoute, dans son
/// conteneur. Aucun n'est publié sur l'hôte (ADR-0018).
class TemplatePort extends $pb.GeneratedMessage {
  factory TemplatePort({
    $core.String? service,
    $core.int? port,
    $core.String? protocol,
    $core.String? description,
  }) {
    final result = TemplatePort._();
    if (service != null) result.service = service;
    if (port != null) result.port = port;
    if (protocol != null) result.protocol = protocol;
    if (description != null) result.description = description;
    return result;
  }

  TemplatePort._();

  factory TemplatePort.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplatePort()..mergeFromBuffer(data, registry);
  factory TemplatePort.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplatePort()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TemplatePort',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: TemplatePort.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'service')
    ..aI(2, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aOS(3, _omitFieldNames ? '' : 'protocol')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplatePort clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplatePort copyWith(void Function(TemplatePort) updates) =>
      super.copyWith((message) => updates(message as TemplatePort))
          as TemplatePort;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TemplatePort() / TemplatePort.new instead')
  static TemplatePort create() => TemplatePort._();
  static $pb.GeneratedMessage $_createMessage() => TemplatePort._();
  @$core.override
  TemplatePort createEmptyInstance() => TemplatePort._();
  @$core.pragma('dart2js:noInline')
  static TemplatePort getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TemplatePort>(
          TemplatePort.$_createMessage);
  static TemplatePort? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get service => $_getSZ(0);
  @$pb.TagNumber(1)
  set service($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasService() => $_has(0);
  @$pb.TagNumber(1)
  void clearService() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get port => $_getIZ(1);
  @$pb.TagNumber(2)
  set port($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPort() => $_has(1);
  @$pb.TagNumber(2)
  void clearPort() => $_clearField(2);

  /// "http" : peut être relié à un domaine servi en HTTPS. "tcp" : interne,
  /// pour d'autres applications du serveur.
  @$pb.TagNumber(3)
  $core.String get protocol => $_getSZ(2);
  @$pb.TagNumber(3)
  set protocol($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProtocol() => $_has(2);
  @$pb.TagNumber(3)
  void clearProtocol() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);
}

/// TemplateVolume est un volume nommé de la stack.
class TemplateVolume extends $pb.GeneratedMessage {
  factory TemplateVolume({
    $core.String? name,
    $core.String? description,
  }) {
    final result = TemplateVolume._();
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    return result;
  }

  TemplateVolume._();

  factory TemplateVolume.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplateVolume()..mergeFromBuffer(data, registry);
  factory TemplateVolume.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TemplateVolume()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TemplateVolume',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: TemplateVolume.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplateVolume clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TemplateVolume copyWith(void Function(TemplateVolume) updates) =>
      super.copyWith((message) => updates(message as TemplateVolume))
          as TemplateVolume;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TemplateVolume() / TemplateVolume.new instead')
  static TemplateVolume create() => TemplateVolume._();
  static $pb.GeneratedMessage $_createMessage() => TemplateVolume._();
  @$core.override
  TemplateVolume createEmptyInstance() => TemplateVolume._();
  @$core.pragma('dart2js:noInline')
  static TemplateVolume getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TemplateVolume>(
          TemplateVolume.$_createMessage);
  static TemplateVolume? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  /// Ce que le volume garde.
  @$pb.TagNumber(2)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => $_clearField(2);
}

/// Template est un template du catalogue (ADR-0018).
class Template extends $pb.GeneratedMessage {
  factory Template({
    $core.String? id,
    $core.String? name,
    $core.String? description,
    $core.String? category,
    $core.String? version,
    $core.String? license,
    $core.String? homepage,
    $core.String? iconSvg,
    $core.Iterable<TemplateVariable>? variables,
    $core.Iterable<TemplatePort>? ports,
    $core.Iterable<TemplateVolume>? volumes,
  }) {
    final result = Template._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (category != null) result.category = category;
    if (version != null) result.version = version;
    if (license != null) result.license = license;
    if (homepage != null) result.homepage = homepage;
    if (iconSvg != null) result.iconSvg = iconSvg;
    if (variables != null) result.variables.addAll(variables);
    if (ports != null) result.ports.addAll(ports);
    if (volumes != null) result.volumes.addAll(volumes);
    return result;
  }

  Template._();

  factory Template.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Template()..mergeFromBuffer(data, registry);
  factory Template.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Template()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Template',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: Template.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..aOS(4, _omitFieldNames ? '' : 'category')
    ..aOS(5, _omitFieldNames ? '' : 'version')
    ..aOS(6, _omitFieldNames ? '' : 'license')
    ..aOS(7, _omitFieldNames ? '' : 'homepage')
    ..aOS(8, _omitFieldNames ? '' : 'iconSvg')
    ..pPM<TemplateVariable>(9, _omitFieldNames ? '' : 'variables',
        subBuilder: TemplateVariable.$_createMessage)
    ..pPM<TemplatePort>(10, _omitFieldNames ? '' : 'ports',
        subBuilder: TemplatePort.$_createMessage)
    ..pPM<TemplateVolume>(11, _omitFieldNames ? '' : 'volumes',
        subBuilder: TemplateVolume.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Template clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Template copyWith(void Function(Template) updates) =>
      super.copyWith((message) => updates(message as Template)) as Template;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Template() / Template.new instead')
  static Template create() => Template._();
  static $pb.GeneratedMessage $_createMessage() => Template._();
  @$core.override
  Template createEmptyInstance() => Template._();
  @$core.pragma('dart2js:noInline')
  static Template getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Template>(Template.$_createMessage);
  static Template? _defaultInstance;

  /// Nom du dossier du template, ex: "postgres".
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  /// Nom affiché, ex: "PostgreSQL".
  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  /// Une phrase.
  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => $_clearField(3);

  /// Une catégorie de la liste fermée : database, cache, storage, cms,
  /// devtools, monitoring, communication ou other.
  @$pb.TagNumber(4)
  $core.String get category => $_getSZ(3);
  @$pb.TagNumber(4)
  set category($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCategory() => $_has(3);
  @$pb.TagNumber(4)
  void clearCategory() => $_clearField(4);

  /// Version de l'application installée, ex: "17.2".
  @$pb.TagNumber(5)
  $core.String get version => $_getSZ(4);
  @$pb.TagNumber(5)
  set version($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasVersion() => $_has(4);
  @$pb.TagNumber(5)
  void clearVersion() => $_clearField(5);

  /// Identifiant SPDX de la licence du logiciel, ex: "PostgreSQL".
  @$pb.TagNumber(6)
  $core.String get license => $_getSZ(5);
  @$pb.TagNumber(6)
  set license($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasLicense() => $_has(5);
  @$pb.TagNumber(6)
  void clearLicense() => $_clearField(6);

  /// Site du logiciel, en https://.
  @$pb.TagNumber(7)
  $core.String get homepage => $_getSZ(6);
  @$pb.TagNumber(7)
  set homepage($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasHomepage() => $_has(6);
  @$pb.TagNumber(7)
  void clearHomepage() => $_clearField(7);

  /// Contenu de icon.svg : carré, sans script ni lien hors du fichier
  /// (vérifié par le serveur au démarrage).
  @$pb.TagNumber(8)
  $core.String get iconSvg => $_getSZ(7);
  @$pb.TagNumber(8)
  set iconSvg($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasIconSvg() => $_has(7);
  @$pb.TagNumber(8)
  void clearIconSvg() => $_clearField(8);

  @$pb.TagNumber(9)
  $pb.PbList<TemplateVariable> get variables => $_getList(8);

  @$pb.TagNumber(10)
  $pb.PbList<TemplatePort> get ports => $_getList(9);

  @$pb.TagNumber(11)
  $pb.PbList<TemplateVolume> get volumes => $_getList(10);
}

class ListTemplatesRequest extends $pb.GeneratedMessage {
  factory ListTemplatesRequest() => ListTemplatesRequest._();

  ListTemplatesRequest._();

  factory ListTemplatesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListTemplatesRequest()..mergeFromBuffer(data, registry);
  factory ListTemplatesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListTemplatesRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTemplatesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListTemplatesRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTemplatesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTemplatesRequest copyWith(void Function(ListTemplatesRequest) updates) =>
      super.copyWith((message) => updates(message as ListTemplatesRequest))
          as ListTemplatesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListTemplatesRequest() / ListTemplatesRequest.new instead')
  static ListTemplatesRequest create() => ListTemplatesRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListTemplatesRequest._();
  @$core.override
  ListTemplatesRequest createEmptyInstance() => ListTemplatesRequest._();
  @$core.pragma('dart2js:noInline')
  static ListTemplatesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTemplatesRequest>(
          ListTemplatesRequest.$_createMessage);
  static ListTemplatesRequest? _defaultInstance;
}

class ListTemplatesResponse extends $pb.GeneratedMessage {
  factory ListTemplatesResponse({
    $core.Iterable<Template>? templates,
  }) {
    final result = ListTemplatesResponse._();
    if (templates != null) result.templates.addAll(templates);
    return result;
  }

  ListTemplatesResponse._();

  factory ListTemplatesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListTemplatesResponse()..mergeFromBuffer(data, registry);
  factory ListTemplatesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListTemplatesResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTemplatesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListTemplatesResponse.$_createMessage)
    ..pPM<Template>(1, _omitFieldNames ? '' : 'templates',
        subBuilder: Template.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTemplatesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTemplatesResponse copyWith(
          void Function(ListTemplatesResponse) updates) =>
      super.copyWith((message) => updates(message as ListTemplatesResponse))
          as ListTemplatesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListTemplatesResponse() / ListTemplatesResponse.new instead')
  static ListTemplatesResponse create() => ListTemplatesResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListTemplatesResponse._();
  @$core.override
  ListTemplatesResponse createEmptyInstance() => ListTemplatesResponse._();
  @$core.pragma('dart2js:noInline')
  static ListTemplatesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTemplatesResponse>(
          ListTemplatesResponse.$_createMessage);
  static ListTemplatesResponse? _defaultInstance;

  /// Triés par identifiant. Vide si le catalogue ne contient aucun template.
  @$pb.TagNumber(1)
  $pb.PbList<Template> get templates => $_getList(0);
}

class GetTemplateRequest extends $pb.GeneratedMessage {
  factory GetTemplateRequest({
    $core.String? id,
  }) {
    final result = GetTemplateRequest._();
    if (id != null) result.id = id;
    return result;
  }

  GetTemplateRequest._();

  factory GetTemplateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTemplateRequest()..mergeFromBuffer(data, registry);
  factory GetTemplateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTemplateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTemplateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetTemplateRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTemplateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTemplateRequest copyWith(void Function(GetTemplateRequest) updates) =>
      super.copyWith((message) => updates(message as GetTemplateRequest))
          as GetTemplateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetTemplateRequest() / GetTemplateRequest.new instead')
  static GetTemplateRequest create() => GetTemplateRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetTemplateRequest._();
  @$core.override
  GetTemplateRequest createEmptyInstance() => GetTemplateRequest._();
  @$core.pragma('dart2js:noInline')
  static GetTemplateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTemplateRequest>(
          GetTemplateRequest.$_createMessage);
  static GetTemplateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class GetTemplateResponse extends $pb.GeneratedMessage {
  factory GetTemplateResponse({
    Template? template,
    $core.String? readme,
  }) {
    final result = GetTemplateResponse._();
    if (template != null) result.template = template;
    if (readme != null) result.readme = readme;
    return result;
  }

  GetTemplateResponse._();

  factory GetTemplateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTemplateResponse()..mergeFromBuffer(data, registry);
  factory GetTemplateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTemplateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTemplateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetTemplateResponse.$_createMessage)
    ..aOM<Template>(1, _omitFieldNames ? '' : 'template',
        subBuilder: Template.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'readme')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTemplateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTemplateResponse copyWith(void Function(GetTemplateResponse) updates) =>
      super.copyWith((message) => updates(message as GetTemplateResponse))
          as GetTemplateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use GetTemplateResponse() / GetTemplateResponse.new instead')
  static GetTemplateResponse create() => GetTemplateResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetTemplateResponse._();
  @$core.override
  GetTemplateResponse createEmptyInstance() => GetTemplateResponse._();
  @$core.pragma('dart2js:noInline')
  static GetTemplateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTemplateResponse>(
          GetTemplateResponse.$_createMessage);
  static GetTemplateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Template get template => $_getN(0);
  @$pb.TagNumber(1)
  set template(Template value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasTemplate() => $_has(0);
  @$pb.TagNumber(1)
  void clearTemplate() => $_clearField(1);
  @$pb.TagNumber(1)
  Template ensureTemplate() => $_ensure(0);

  /// README.md du template, en Markdown : le service, ses variables, comment
  /// s'y connecter.
  @$pb.TagNumber(2)
  $core.String get readme => $_getSZ(1);
  @$pb.TagNumber(2)
  set readme($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReadme() => $_has(1);
  @$pb.TagNumber(2)
  void clearReadme() => $_clearField(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
