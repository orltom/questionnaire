//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ChallengeQuestionSummary {
  /// Returns a new [ChallengeQuestionSummary] instance.
  ChallengeQuestionSummary({
    required this.id,
    required this.description,
    required this.position,
  });

  String id;

  String description;

  /// Minimum value: 0
  int position;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ChallengeQuestionSummary &&
    other.id == id &&
    other.description == description &&
    other.position == position;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (description.hashCode) +
    (position.hashCode);

  @override
  String toString() => 'ChallengeQuestionSummary[id=$id, description=$description, position=$position]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'description'] = this.description;
      json[r'position'] = this.position;
    return json;
  }

  /// Returns a new [ChallengeQuestionSummary] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ChallengeQuestionSummary? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ChallengeQuestionSummary[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ChallengeQuestionSummary[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ChallengeQuestionSummary(
        id: mapValueOfType<String>(json, r'id')!,
        description: mapValueOfType<String>(json, r'description')!,
        position: mapValueOfType<int>(json, r'position')!,
      );
    }
    return null;
  }

  static List<ChallengeQuestionSummary> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ChallengeQuestionSummary>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ChallengeQuestionSummary.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ChallengeQuestionSummary> mapFromJson(dynamic json) {
    final map = <String, ChallengeQuestionSummary>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ChallengeQuestionSummary.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ChallengeQuestionSummary-objects as value to a dart map
  static Map<String, List<ChallengeQuestionSummary>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ChallengeQuestionSummary>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ChallengeQuestionSummary.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'description',
    'position',
  };
}

