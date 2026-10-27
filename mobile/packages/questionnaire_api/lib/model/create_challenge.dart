//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class CreateChallenge {
  /// Returns a new [CreateChallenge] instance.
  CreateChallenge({
    required this.quizId,
    required this.startsAt,
    this.endsAt,
    required this.access,
    this.invitees = const [],
  });

  String quizId;

  DateTime startsAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? endsAt;

  Access access;

  List<String> invitees;

  @override
  bool operator ==(Object other) => identical(this, other) || other is CreateChallenge &&
    other.quizId == quizId &&
    other.startsAt == startsAt &&
    other.endsAt == endsAt &&
    other.access == access &&
    _deepEquality.equals(other.invitees, invitees);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (quizId.hashCode) +
    (startsAt.hashCode) +
    (endsAt == null ? 0 : endsAt!.hashCode) +
    (access.hashCode) +
    (invitees.hashCode);

  @override
  String toString() => 'CreateChallenge[quizId=$quizId, startsAt=$startsAt, endsAt=$endsAt, access=$access, invitees=$invitees]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'quizId'] = this.quizId;
      json[r'startsAt'] = this.startsAt.toUtc().toIso8601String();
    if (this.endsAt != null) {
      json[r'endsAt'] = this.endsAt!.toUtc().toIso8601String();
    } else {
      json[r'endsAt'] = null;
    }
      json[r'access'] = this.access;
      json[r'invitees'] = this.invitees;
    return json;
  }

  /// Returns a new [CreateChallenge] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static CreateChallenge? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "CreateChallenge[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "CreateChallenge[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return CreateChallenge(
        quizId: mapValueOfType<String>(json, r'quizId')!,
        startsAt: mapDateTime(json, r'startsAt', r'')!,
        endsAt: mapDateTime(json, r'endsAt', r''),
        access: Access.fromJson(json[r'access'])!,
        invitees: json[r'invitees'] is Iterable
            ? (json[r'invitees'] as Iterable).cast<String>().toList(growable: false)
            : const [],
      );
    }
    return null;
  }

  static List<CreateChallenge> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <CreateChallenge>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CreateChallenge.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, CreateChallenge> mapFromJson(dynamic json) {
    final map = <String, CreateChallenge>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = CreateChallenge.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of CreateChallenge-objects as value to a dart map
  static Map<String, List<CreateChallenge>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<CreateChallenge>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = CreateChallenge.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'quizId',
    'startsAt',
    'access',
  };
}

