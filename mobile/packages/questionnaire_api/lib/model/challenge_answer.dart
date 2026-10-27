//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ChallengeAnswer {
  /// Returns a new [ChallengeAnswer] instance.
  ChallengeAnswer({
    required this.id,
    required this.description,
    required this.position,
  });

  String id;

  String description;

  /// Minimum value: 0
  int position;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ChallengeAnswer &&
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
  String toString() => 'ChallengeAnswer[id=$id, description=$description, position=$position]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'description'] = this.description;
      json[r'position'] = this.position;
    return json;
  }

  /// Returns a new [ChallengeAnswer] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ChallengeAnswer? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ChallengeAnswer[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ChallengeAnswer[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ChallengeAnswer(
        id: mapValueOfType<String>(json, r'id')!,
        description: mapValueOfType<String>(json, r'description')!,
        position: mapValueOfType<int>(json, r'position')!,
      );
    }
    return null;
  }

  static List<ChallengeAnswer> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ChallengeAnswer>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ChallengeAnswer.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ChallengeAnswer> mapFromJson(dynamic json) {
    final map = <String, ChallengeAnswer>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ChallengeAnswer.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ChallengeAnswer-objects as value to a dart map
  static Map<String, List<ChallengeAnswer>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ChallengeAnswer>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ChallengeAnswer.listFromJson(entry.value, growable: growable,);
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

