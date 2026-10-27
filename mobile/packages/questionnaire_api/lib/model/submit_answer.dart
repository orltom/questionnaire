//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubmitAnswer {
  /// Returns a new [SubmitAnswer] instance.
  SubmitAnswer({
    required this.questionId,
    required this.answerId,
  });

  String questionId;

  String answerId;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubmitAnswer &&
    other.questionId == questionId &&
    other.answerId == answerId;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (questionId.hashCode) +
    (answerId.hashCode);

  @override
  String toString() => 'SubmitAnswer[questionId=$questionId, answerId=$answerId]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'questionId'] = this.questionId;
      json[r'answerId'] = this.answerId;
    return json;
  }

  /// Returns a new [SubmitAnswer] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubmitAnswer? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubmitAnswer[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubmitAnswer[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubmitAnswer(
        questionId: mapValueOfType<String>(json, r'questionId')!,
        answerId: mapValueOfType<String>(json, r'answerId')!,
      );
    }
    return null;
  }

  static List<SubmitAnswer> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubmitAnswer>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubmitAnswer.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubmitAnswer> mapFromJson(dynamic json) {
    final map = <String, SubmitAnswer>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubmitAnswer.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubmitAnswer-objects as value to a dart map
  static Map<String, List<SubmitAnswer>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubmitAnswer>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubmitAnswer.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'questionId',
    'answerId',
  };
}

