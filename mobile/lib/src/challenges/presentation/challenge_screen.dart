import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';
import '../domain/challenge.dart';
import '../domain/challenge_question.dart';
import '../domain/challenge_repository.dart';

class ChallengeScreen extends StatefulWidget {
  const ChallengeScreen({
    super.key,
    required this.challenge,
    required this.repository,
  });

  final Challenge challenge;
  final ChallengeRepository repository;

  @override
  State<ChallengeScreen> createState() => _ChallengeScreenState();
}

class _ChallengeScreenState extends State<ChallengeScreen> {
  late final Future<List<ChallengeQuestion>> _questions = widget.repository
      .questions(widget.challenge.id);

  int _current = 0;
  bool _submitting = false;

  Future<void> _answer(
    ChallengeQuestion question,
    ChallengeAnswer answer,
  ) async {
    setState(() => _submitting = true);
    await widget.repository.submitAnswer(
      challengeId: widget.challenge.id,
      questionId: question.id,
      answerId: answer.id,
    );
    if (!mounted) return;
    setState(() {
      _current++;
      _submitting = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return Scaffold(
      appBar: AppBar(title: Text(widget.challenge.quizTitle)),
      body: FutureBuilder<List<ChallengeQuestion>>(
        future: _questions,
        builder: (context, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return const Center(child: Text('Could not load questions'));
          }

          final questions = snapshot.requireData;
          if (_current >= questions.length) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.celebration_rounded,
                      size: 72,
                      color: colorScheme.primary,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      'Questionnaire finished',
                      style: textTheme.headlineSmall,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Thanks for taking part.',
                      style: textTheme.bodyMedium?.copyWith(
                        color: colorScheme.onSurfaceVariant,
                      ),
                    ),
                    const SizedBox(height: 24),
                    FilledButton(
                      onPressed: () => Navigator.of(context).pop(),
                      child: const Text('Back to questionnaires'),
                    ),
                  ],
                ),
              ),
            );
          }

          final question = questions[_current];
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              LinearProgressIndicator(value: _current / questions.length),
              const SizedBox(height: 8),
              Text(
                'Question ${_current + 1} of ${questions.length}',
                style: textTheme.labelLarge?.copyWith(
                  color: colorScheme.onSurfaceVariant,
                ),
              ),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: colorScheme.primaryContainer,
                  borderRadius: appCardRadius,
                ),
                child: Text(
                  question.description,
                  style: textTheme.titleLarge?.copyWith(
                    color: colorScheme.onPrimaryContainer,
                  ),
                ),
              ),
              const SizedBox(height: 24),
              for (final answer in question.answers)
                Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: FilledButton.tonal(
                    onPressed: _submitting
                        ? null
                        : () => _answer(question, answer),
                    style: FilledButton.styleFrom(
                      minimumSize: const Size.fromHeight(56),
                      alignment: Alignment.centerLeft,
                      padding: const EdgeInsets.symmetric(horizontal: 20),
                    ),
                    child: Text(answer.description),
                  ),
                ),
            ],
          );
        },
      ),
    );
  }
}
