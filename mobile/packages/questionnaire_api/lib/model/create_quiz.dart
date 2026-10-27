//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class CreateQuiz {
  /// Returns a new [CreateQuiz] instance.
  CreateQuiz({
    required this.title,
    this.description = '',
    required this.visibility,
    this.questions = const [],
  });

  String title;

  String description;

  Visibility visibility;

  List<QuizQuestion> questions;

  @override
  bool operator ==(Object other) => identical(this, other) || other is CreateQuiz &&
    other.title == title &&
    other.description == description &&
    other.visibility == visibility &&
    _deepEquality.equals(other.questions, questions);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (title.hashCode) +
    (description.hashCode) +
    (visibility.hashCode) +
    (questions.hashCode);

  @override
  String toString() => 'CreateQuiz[title=$title, description=$description, visibility=$visibility, questions=$questions]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'title'] = this.title;
      json[r'description'] = this.description;
      json[r'visibility'] = this.visibility;
      json[r'questions'] = this.questions;
    return json;
  }

  /// Returns a new [CreateQuiz] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static CreateQuiz? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "CreateQuiz[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "CreateQuiz[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return CreateQuiz(
        title: mapValueOfType<String>(json, r'title')!,
        description: mapValueOfType<String>(json, r'description')!,
        visibility: Visibility.fromJson(json[r'visibility'])!,
        questions: QuizQuestion.listFromJson(json[r'questions']),
      );
    }
    return null;
  }

  static List<CreateQuiz> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <CreateQuiz>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CreateQuiz.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, CreateQuiz> mapFromJson(dynamic json) {
    final map = <String, CreateQuiz>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = CreateQuiz.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of CreateQuiz-objects as value to a dart map
  static Map<String, List<CreateQuiz>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<CreateQuiz>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = CreateQuiz.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'title',
    'description',
    'visibility',
  };
}

