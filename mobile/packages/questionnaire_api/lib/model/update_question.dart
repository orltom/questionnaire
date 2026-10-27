//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class UpdateQuestion {
  /// Returns a new [UpdateQuestion] instance.
  UpdateQuestion({
    required this.description,
    required this.visibility,
    this.labels = const [],
    this.answers = const [],
  });

  String description;

  Visibility visibility;

  List<String> labels;

  List<Answer> answers;

  @override
  bool operator ==(Object other) => identical(this, other) || other is UpdateQuestion &&
    other.description == description &&
    other.visibility == visibility &&
    _deepEquality.equals(other.labels, labels) &&
    _deepEquality.equals(other.answers, answers);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (description.hashCode) +
    (visibility.hashCode) +
    (labels.hashCode) +
    (answers.hashCode);

  @override
  String toString() => 'UpdateQuestion[description=$description, visibility=$visibility, labels=$labels, answers=$answers]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'description'] = this.description;
      json[r'visibility'] = this.visibility;
      json[r'labels'] = this.labels;
      json[r'answers'] = this.answers;
    return json;
  }

  /// Returns a new [UpdateQuestion] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static UpdateQuestion? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "UpdateQuestion[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "UpdateQuestion[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return UpdateQuestion(
        description: mapValueOfType<String>(json, r'description')!,
        visibility: Visibility.fromJson(json[r'visibility'])!,
        labels: json[r'labels'] is Iterable
            ? (json[r'labels'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        answers: Answer.listFromJson(json[r'answers']),
      );
    }
    return null;
  }

  static List<UpdateQuestion> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <UpdateQuestion>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = UpdateQuestion.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, UpdateQuestion> mapFromJson(dynamic json) {
    final map = <String, UpdateQuestion>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = UpdateQuestion.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of UpdateQuestion-objects as value to a dart map
  static Map<String, List<UpdateQuestion>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<UpdateQuestion>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = UpdateQuestion.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'description',
    'visibility',
    'answers',
  };
}

