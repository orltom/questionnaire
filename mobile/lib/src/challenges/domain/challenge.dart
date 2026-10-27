enum ChallengeState { scheduled, running, finished }

class Challenge {
  const Challenge({
    required this.id,
    required this.quizTitle,
    required this.state,
    required this.startsAt,
    this.endsAt,
  });

  final String id;
  final String quizTitle;
  final ChallengeState state;
  final DateTime startsAt;
  final DateTime? endsAt;

  bool get canStart => state == ChallengeState.running;
}
