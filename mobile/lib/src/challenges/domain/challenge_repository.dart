import 'challenge.dart';
import 'challenge_question.dart';

abstract interface class ChallengeRepository {
  Future<List<Challenge>> startable();

  Future<List<ChallengeQuestion>> questions(String challengeId);

  Future<void> submitAnswer({
    required String challengeId,
    required String questionId,
    required String answerId,
  });
}
