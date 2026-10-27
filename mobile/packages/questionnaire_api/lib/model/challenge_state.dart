//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class ChallengeState {
  /// Instantiate a new enum with the provided [value].
  const ChallengeState._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const scheduled = ChallengeState._(r'scheduled');
  static const running = ChallengeState._(r'running');
  static const finished = ChallengeState._(r'finished');

  /// List of all possible values in this [enum][ChallengeState].
  static const values = <ChallengeState>[
    scheduled,
    running,
    finished,
  ];

  static ChallengeState? fromJson(dynamic value) => ChallengeStateTypeTransformer().decode(value);

  static List<ChallengeState> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ChallengeState>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ChallengeState.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [ChallengeState] to String,
/// and [decode] dynamic data back to [ChallengeState].
class ChallengeStateTypeTransformer {
  factory ChallengeStateTypeTransformer() => _instance ??= const ChallengeStateTypeTransformer._();

  const ChallengeStateTypeTransformer._();

  String encode(ChallengeState data) => data.value;

  /// Decodes a [dynamic value][data] to a ChallengeState.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  ChallengeState? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'scheduled': return ChallengeState.scheduled;
        case r'running': return ChallengeState.running;
        case r'finished': return ChallengeState.finished;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [ChallengeStateTypeTransformer] instance.
  static ChallengeStateTypeTransformer? _instance;
}

