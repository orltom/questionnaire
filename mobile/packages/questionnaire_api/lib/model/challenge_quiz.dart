//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ChallengeQuiz {
  /// Returns a new [ChallengeQuiz] instance.
  ChallengeQuiz({
    required this.quizId,
    required this.title,
    this.questions = const [],
  });

  String quizId;

  String title;

  List<ChallengeQuestion> questions;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ChallengeQuiz &&
    other.quizId == quizId &&
    other.title == title &&
    _deepEquality.equals(other.questions, questions);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (quizId.hashCode) +
    (title.hashCode) +
    (questions.hashCode);

  @override
  String toString() => 'ChallengeQuiz[quizId=$quizId, title=$title, questions=$questions]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'quizId'] = this.quizId;
      json[r'title'] = this.title;
      json[r'questions'] = this.questions;
    return json;
  }

  /// Returns a new [ChallengeQuiz] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ChallengeQuiz? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ChallengeQuiz[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ChallengeQuiz[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ChallengeQuiz(
        quizId: mapValueOfType<String>(json, r'quizId')!,
        title: mapValueOfType<String>(json, r'title')!,
        questions: ChallengeQuestion.listFromJson(json[r'questions']),
      );
    }
    return null;
  }

  static List<ChallengeQuiz> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ChallengeQuiz>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ChallengeQuiz.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ChallengeQuiz> mapFromJson(dynamic json) {
    final map = <String, ChallengeQuiz>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ChallengeQuiz.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ChallengeQuiz-objects as value to a dart map
  static Map<String, List<ChallengeQuiz>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ChallengeQuiz>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ChallengeQuiz.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'quizId',
    'title',
    'questions',
  };
}

