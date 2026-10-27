import 'package:flutter/foundation.dart';

import '../domain/challenge.dart';
import '../domain/challenge_question.dart';
import '../domain/challenge_repository.dart';

class DemoChallengeRepository implements ChallengeRepository {
  const DemoChallengeRepository();

  static const _questions = {
    '1': [
      ChallengeQuestion(
        id: 'q1',
        description: 'What is the capital of France?',
        answers: [
          ChallengeAnswer(id: 'a1', description: 'Paris'),
          ChallengeAnswer(id: 'a2', description: 'Lyon'),
          ChallengeAnswer(id: 'a3', description: 'Marseille'),
        ],
      ),
      ChallengeQuestion(
        id: 'q2',
        description: 'How many continents are there?',
        answers: [
          ChallengeAnswer(id: 'a4', description: '5'),
          ChallengeAnswer(id: 'a5', description: '6'),
          ChallengeAnswer(id: 'a6', description: '7'),
        ],
      ),
    ],
    '2': [
      ChallengeQuestion(
        id: 'q3',
        description: 'Which keyword declares a function in Go?',
        answers: [
          ChallengeAnswer(id: 'a7', description: 'func'),
          ChallengeAnswer(id: 'a8', description: 'def'),
          ChallengeAnswer(id: 'a9', description: 'fn'),
        ],
      ),
      ChallengeQuestion(
        id: 'q4',
        description: 'Which operator declares and assigns a variable?',
        answers: [
          ChallengeAnswer(id: 'a10', description: ':='),
          ChallengeAnswer(id: 'a11', description: '='),
          ChallengeAnswer(id: 'a12', description: 'var'),
        ],
      ),
    ],
  };

  @override
  Future<List<Challenge>> startable() async {
    return [
      Challenge(
        id: '1',
        quizTitle: 'General knowledge',
        state: ChallengeState.running,
        startsAt: DateTime(2026, 10, 27, 9),
        endsAt: DateTime(2026, 10, 31, 18),
      ),
      Challenge(
        id: '2',
        quizTitle: 'Go basics',
        state: ChallengeState.running,
        startsAt: DateTime(2026, 10, 26, 14),
      ),
    ];
  }

  @override
  Future<List<ChallengeQuestion>> questions(String challengeId) async {
    return _questions[challengeId] ?? const [];
  }

  @override
  Future<void> submitAnswer({
    required String challengeId,
    required String questionId,
    required String answerId,
  }) async {
    debugPrint(
      'challenge $challengeId: question $questionId answered with $answerId',
    );
  }
}
