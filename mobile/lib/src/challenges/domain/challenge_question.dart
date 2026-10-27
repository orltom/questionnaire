class ChallengeAnswer {
  const ChallengeAnswer({required this.id, required this.description});

  final String id;
  final String description;
}

class ChallengeQuestion {
  const ChallengeQuestion({
    required this.id,
    required this.description,
    required this.answers,
  });

  final String id;
  final String description;
  final List<ChallengeAnswer> answers;
}
