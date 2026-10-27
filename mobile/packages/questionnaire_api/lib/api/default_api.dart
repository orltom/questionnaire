//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class DefaultApi {
  DefaultApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Add question to quiz
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  ///
  /// * [AddQuestionToQuiz] addQuestionToQuiz (required):
  Future<Response> addQuestionToQuizWithHttpInfo(String quizId, AddQuestionToQuiz addQuestionToQuiz,) async {
    // ignore: prefer_const_declarations
    final path = r'/quizzes/{quizId}/questions'
      .replaceAll('{quizId}', quizId);

    // ignore: prefer_final_locals
    Object? postBody = addQuestionToQuiz;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Add question to quiz
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  ///
  /// * [AddQuestionToQuiz] addQuestionToQuiz (required):
  Future<void> addQuestionToQuiz(String quizId, AddQuestionToQuiz addQuestionToQuiz,) async {
    final response = await addQuestionToQuizWithHttpInfo(quizId, addQuestionToQuiz,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Close challenge
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<Response> closeChallengeWithHttpInfo(String challengeId,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}/close'
      .replaceAll('{challengeId}', challengeId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Close challenge
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<void> closeChallenge(String challengeId,) async {
    final response = await closeChallengeWithHttpInfo(challengeId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Create challenge
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [CreateChallenge] createChallenge (required):
  Future<Response> createChallengeWithHttpInfo(CreateChallenge createChallenge,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges';

    // ignore: prefer_final_locals
    Object? postBody = createChallenge;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Create challenge
  ///
  /// Parameters:
  ///
  /// * [CreateChallenge] createChallenge (required):
  Future<Challenge?> createChallenge(CreateChallenge createChallenge,) async {
    final response = await createChallengeWithHttpInfo(createChallenge,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'Challenge',) as Challenge;
    
    }
    return null;
  }

  /// Create question
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [CreateQuestion] createQuestion (required):
  Future<Response> createQuestionWithHttpInfo(CreateQuestion createQuestion,) async {
    // ignore: prefer_const_declarations
    final path = r'/questions';

    // ignore: prefer_final_locals
    Object? postBody = createQuestion;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Create question
  ///
  /// Parameters:
  ///
  /// * [CreateQuestion] createQuestion (required):
  Future<Question?> createQuestion(CreateQuestion createQuestion,) async {
    final response = await createQuestionWithHttpInfo(createQuestion,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'Question',) as Question;
    
    }
    return null;
  }

  /// Create quiz
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [CreateQuiz] createQuiz (required):
  Future<Response> createQuizWithHttpInfo(CreateQuiz createQuiz,) async {
    // ignore: prefer_const_declarations
    final path = r'/quizzes';

    // ignore: prefer_final_locals
    Object? postBody = createQuiz;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Create quiz
  ///
  /// Parameters:
  ///
  /// * [CreateQuiz] createQuiz (required):
  Future<Quiz?> createQuiz(CreateQuiz createQuiz,) async {
    final response = await createQuizWithHttpInfo(createQuiz,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'Quiz',) as Quiz;
    
    }
    return null;
  }

  /// Delete challenge
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<Response> deleteChallengeWithHttpInfo(String challengeId,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}'
      .replaceAll('{challengeId}', challengeId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete challenge
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<void> deleteChallenge(String challengeId,) async {
    final response = await deleteChallengeWithHttpInfo(challengeId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Delete question
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  Future<Response> deleteQuestionWithHttpInfo(String questionId,) async {
    // ignore: prefer_const_declarations
    final path = r'/questions/{questionId}'
      .replaceAll('{questionId}', questionId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete question
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  Future<void> deleteQuestion(String questionId,) async {
    final response = await deleteQuestionWithHttpInfo(questionId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Delete quiz
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  Future<Response> deleteQuizWithHttpInfo(String quizId,) async {
    // ignore: prefer_const_declarations
    final path = r'/quizzes/{quizId}'
      .replaceAll('{quizId}', quizId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete quiz
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  Future<void> deleteQuiz(String quizId,) async {
    final response = await deleteQuizWithHttpInfo(quizId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get challenge overview
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<Response> getChallengeWithHttpInfo(String challengeId,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}'
      .replaceAll('{challengeId}', challengeId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get challenge overview
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<ChallengeOverview?> getChallenge(String challengeId,) async {
    final response = await getChallengeWithHttpInfo(challengeId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ChallengeOverview',) as ChallengeOverview;
    
    }
    return null;
  }

  /// Get question from challenge
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  ///
  /// * [String] questionId (required):
  Future<Response> getChallengeQuestionWithHttpInfo(String challengeId, String questionId,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}/question/{questionId}'
      .replaceAll('{challengeId}', challengeId)
      .replaceAll('{questionId}', questionId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get question from challenge
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  ///
  /// * [String] questionId (required):
  Future<ChallengeQuestion?> getChallengeQuestion(String challengeId, String questionId,) async {
    final response = await getChallengeQuestionWithHttpInfo(challengeId, questionId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ChallengeQuestion',) as ChallengeQuestion;
    
    }
    return null;
  }

  /// Get challenge results
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<Response> getChallengeResultsWithHttpInfo(String challengeId,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}/results'
      .replaceAll('{challengeId}', challengeId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get challenge results
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  Future<List<Result>?> getChallengeResults(String challengeId,) async {
    final response = await getChallengeResultsWithHttpInfo(challengeId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      final responseBody = await _decodeBodyBytes(response);
      return (await apiClient.deserializeAsync(responseBody, 'List<Result>') as List)
        .cast<Result>()
        .toList(growable: false);

    }
    return null;
  }

  /// Get question
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  Future<Response> getQuestionWithHttpInfo(String questionId,) async {
    // ignore: prefer_const_declarations
    final path = r'/questions/{questionId}'
      .replaceAll('{questionId}', questionId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get question
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  Future<Question?> getQuestion(String questionId,) async {
    final response = await getQuestionWithHttpInfo(questionId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'Question',) as Question;
    
    }
    return null;
  }

  /// Get quiz
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  Future<Response> getQuizWithHttpInfo(String quizId,) async {
    // ignore: prefer_const_declarations
    final path = r'/quizzes/{quizId}'
      .replaceAll('{quizId}', quizId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get quiz
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  Future<Quiz?> getQuiz(String quizId,) async {
    final response = await getQuizWithHttpInfo(quizId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'Quiz',) as Quiz;
    
    }
    return null;
  }

  /// Submit an answer
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  ///
  /// * [SubmitAnswer] submitAnswer (required):
  Future<Response> submitAnswerWithHttpInfo(String challengeId, SubmitAnswer submitAnswer,) async {
    // ignore: prefer_const_declarations
    final path = r'/challenges/{challengeId}/answers'
      .replaceAll('{challengeId}', challengeId);

    // ignore: prefer_final_locals
    Object? postBody = submitAnswer;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Submit an answer
  ///
  /// Parameters:
  ///
  /// * [String] challengeId (required):
  ///
  /// * [SubmitAnswer] submitAnswer (required):
  Future<void> submitAnswer(String challengeId, SubmitAnswer submitAnswer,) async {
    final response = await submitAnswerWithHttpInfo(challengeId, submitAnswer,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Update question
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  ///
  /// * [UpdateQuestion] updateQuestion (required):
  Future<Response> updateQuestionWithHttpInfo(String questionId, UpdateQuestion updateQuestion,) async {
    // ignore: prefer_const_declarations
    final path = r'/questions/{questionId}'
      .replaceAll('{questionId}', questionId);

    // ignore: prefer_final_locals
    Object? postBody = updateQuestion;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Update question
  ///
  /// Parameters:
  ///
  /// * [String] questionId (required):
  ///
  /// * [UpdateQuestion] updateQuestion (required):
  Future<void> updateQuestion(String questionId, UpdateQuestion updateQuestion,) async {
    final response = await updateQuestionWithHttpInfo(questionId, updateQuestion,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Update quiz
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  ///
  /// * [UpdateQuiz] updateQuiz (required):
  Future<Response> updateQuizWithHttpInfo(String quizId, UpdateQuiz updateQuiz,) async {
    // ignore: prefer_const_declarations
    final path = r'/quizzes/{quizId}'
      .replaceAll('{quizId}', quizId);

    // ignore: prefer_final_locals
    Object? postBody = updateQuiz;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Update quiz
  ///
  /// Parameters:
  ///
  /// * [String] quizId (required):
  ///
  /// * [UpdateQuiz] updateQuiz (required):
  Future<void> updateQuiz(String quizId, UpdateQuiz updateQuiz,) async {
    final response = await updateQuizWithHttpInfo(quizId, updateQuiz,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }
}
