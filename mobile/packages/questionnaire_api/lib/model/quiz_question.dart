//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class QuizQuestion {
  /// Returns a new [QuizQuestion] instance.
  QuizQuestion({
    required this.questionId,
    required this.position,
  });

  String questionId;

  /// Minimum value: 0
  int position;

  @override
  bool operator ==(Object other) => identical(this, other) || other is QuizQuestion &&
    other.questionId == questionId &&
    other.position == position;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (questionId.hashCode) +
    (position.hashCode);

  @override
  String toString() => 'QuizQuestion[questionId=$questionId, position=$position]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'questionId'] = this.questionId;
      json[r'position'] = this.position;
    return json;
  }

  /// Returns a new [QuizQuestion] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static QuizQuestion? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "QuizQuestion[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "QuizQuestion[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return QuizQuestion(
        questionId: mapValueOfType<String>(json, r'questionId')!,
        position: mapValueOfType<int>(json, r'position')!,
      );
    }
    return null;
  }

  static List<QuizQuestion> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <QuizQuestion>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = QuizQuestion.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, QuizQuestion> mapFromJson(dynamic json) {
    final map = <String, QuizQuestion>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = QuizQuestion.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of QuizQuestion-objects as value to a dart map
  static Map<String, List<QuizQuestion>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<QuizQuestion>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = QuizQuestion.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'questionId',
    'position',
  };
}

