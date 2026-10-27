import 'package:flutter/material.dart';

import 'src/challenges/data/demo_challenge_repository.dart';
import 'src/challenges/domain/challenge_repository.dart';
import 'src/challenges/presentation/challenge_list_screen.dart';
import 'src/core/theme/app_theme.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatelessWidget {
  const MainApp({super.key, this.repository = const DemoChallengeRepository()});

  final ChallengeRepository repository;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: buildAppTheme(Brightness.light),
      darkTheme: buildAppTheme(Brightness.dark),
      home: ChallengeListScreen(repository: repository),
    );
  }
}
