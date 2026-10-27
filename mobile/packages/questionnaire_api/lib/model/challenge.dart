//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Challenge {
  /// Returns a new [Challenge] instance.
  Challenge({
    required this.id,
    required this.quiz,
    required this.startsAt,
    this.endsAt,
    this.closedAt,
    required this.access,
    this.invitees = const [],
    required this.state,
  });

  String id;

  ChallengeQuiz quiz;

  DateTime startsAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? endsAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? closedAt;

  Access access;

  List<String> invitees;

  ChallengeState state;

  @override
  bool operator ==(Object other) => identical(this, other) || other is Challenge &&
    other.id == id &&
    other.quiz == quiz &&
    other.startsAt == startsAt &&
    other.endsAt == endsAt &&
    other.closedAt == closedAt &&
    other.access == access &&
    _deepEquality.equals(other.invitees, invitees) &&
    other.state == state;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (quiz.hashCode) +
    (startsAt.hashCode) +
    (endsAt == null ? 0 : endsAt!.hashCode) +
    (closedAt == null ? 0 : closedAt!.hashCode) +
    (access.hashCode) +
    (invitees.hashCode) +
    (state.hashCode);

  @override
  String toString() => 'Challenge[id=$id, quiz=$quiz, startsAt=$startsAt, endsAt=$endsAt, closedAt=$closedAt, access=$access, invitees=$invitees, state=$state]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'quiz'] = this.quiz;
      json[r'startsAt'] = this.startsAt.toUtc().toIso8601String();
    if (this.endsAt != null) {
      json[r'endsAt'] = this.endsAt!.toUtc().toIso8601String();
    } else {
      json[r'endsAt'] = null;
    }
    if (this.closedAt != null) {
      json[r'closedAt'] = this.closedAt!.toUtc().toIso8601String();
    } else {
      json[r'closedAt'] = null;
    }
      json[r'access'] = this.access;
      json[r'invitees'] = this.invitees;
      json[r'state'] = this.state;
    return json;
  }

  /// Returns a new [Challenge] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Challenge? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "Challenge[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "Challenge[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return Challenge(
        id: mapValueOfType<String>(json, r'id')!,
        quiz: ChallengeQuiz.fromJson(json[r'quiz'])!,
        startsAt: mapDateTime(json, r'startsAt', r'')!,
        endsAt: mapDateTime(json, r'endsAt', r''),
        closedAt: mapDateTime(json, r'closedAt', r''),
        access: Access.fromJson(json[r'access'])!,
        invitees: json[r'invitees'] is Iterable
            ? (json[r'invitees'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        state: ChallengeState.fromJson(json[r'state'])!,
      );
    }
    return null;
  }

  static List<Challenge> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <Challenge>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Challenge.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Challenge> mapFromJson(dynamic json) {
    final map = <String, Challenge>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Challenge.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Challenge-objects as value to a dart map
  static Map<String, List<Challenge>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<Challenge>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Challenge.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'quiz',
    'startsAt',
    'access',
    'invitees',
    'state',
  };
}

