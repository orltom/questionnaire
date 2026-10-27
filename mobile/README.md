# Questionnaire Mobile Application

[![Mobile](https://github.com/orltom/questionnaire/actions/workflows/mobile.yml/badge.svg?branch=main)](https://github.com/orltom/questionnaire/actions/workflows/mobile.yml)

Mobile application to list and start questionnaires.

## Requirements

- [Flutter SDK](https://docs.flutter.dev/get-started/install) with the Android SDK
- [Task](https://taskfile.dev)
- Docker (only to regenerate the API client)

## Getting Started

```bash
task check      # analyze and test
task generate   # regenerate the API client from api/openapi.yaml (needs Docker)
task build      # build android release APK
```
```bash
flutter emulators --launch pixel7 # use your android emulator
flutter run                       # run the app
```