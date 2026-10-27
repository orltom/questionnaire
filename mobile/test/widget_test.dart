import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/main.dart';
import 'package:mobile/src/challenges/domain/challenge.dart';
import 'package:mobile/src/challenges/domain/challenge_question.dart';
import 'package:mobile/src/challenges/domain/challenge_repository.dart';
import 'package:mobile/src/challenges/presentation/challenge_list_screen.dart';

class FakeChallengeRepository implements ChallengeRepository {
  FakeChallengeRepository({
    this.challenges = const [],
    this.questionItems = const [],
  });

  final List<Challenge> challenges;
  final List<ChallengeQuestion> questionItems;
  final List<String> submitted = [];

  @override
  Future<List<Challenge>> startable() async => challenges;

  @override
  Future<List<ChallengeQuestion>> questions(String challengeId) async =>
      questionItems;

  @override
  Future<void> submitAnswer({
    required String challengeId,
    required String questionId,
    required String answerId,
  }) async {
    submitted.add('$questionId:$answerId');
  }
}

final _goBasics = Challenge(
  id: '1',
  quizTitle: 'Go basics',
  state: ChallengeState.running,
  startsAt: DateTime(2026, 10, 26, 14),
);

const _questions = [
  ChallengeQuestion(
    id: 'q1',
    description: 'Which keyword declares a function?',
    answers: [
      ChallengeAnswer(id: 'a1', description: 'func'),
      ChallengeAnswer(id: 'a2', description: 'def'),
    ],
  ),
  ChallengeQuestion(
    id: 'q2',
    description: 'Which operator assigns a new variable?',
    answers: [
      ChallengeAnswer(id: 'a3', description: ':='),
      ChallengeAnswer(id: 'a4', description: '='),
    ],
  ),
];

void main() {
  testWidgets('lists the questionnaires that can be started', (tester) async {
    await tester.pumpWidget(
      MainApp(repository: FakeChallengeRepository(challenges: [_goBasics])),
    );
    await tester.pumpAndSettle();

    expect(find.text('Questionnaires'), findsOneWidget);
    expect(find.text('Go basics'), findsOneWidget);
  });

  testWidgets('shows a message when nothing can be started', (tester) async {
    await tester.pumpWidget(MainApp(repository: FakeChallengeRepository()));
    await tester.pumpAndSettle();

    expect(find.text('No questionnaires to start'), findsOneWidget);
  });

  testWidgets('answering moves through the questions and submits each answer', (
    tester,
  ) async {
    final repository = FakeChallengeRepository(
      challenges: [_goBasics],
      questionItems: _questions,
    );

    await tester.pumpWidget(
      MaterialApp(home: ChallengeListScreen(repository: repository)),
    );
    await tester.pumpAndSettle();
    await tester.tap(find.text('Go basics'));
    await tester.pumpAndSettle();

    expect(find.text('Which keyword declares a function?'), findsOneWidget);
    await tester.tap(find.text('func'));
    await tester.pumpAndSettle();

    expect(find.text('Which operator assigns a new variable?'), findsOneWidget);
    await tester.tap(find.text(':='));
    await tester.pumpAndSettle();

    expect(find.text('Questionnaire finished'), findsOneWidget);
    expect(repository.submitted, ['q1:a1', 'q2:a3']);
  });
}
