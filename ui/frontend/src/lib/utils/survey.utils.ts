import type {
    SurveyDimensionsResponse,
    SurveyIndicatorsResponse,
    SurveyQuestionsResponse,
    SurveysResponse,
    SurveyThemesResponse
} from "$lib/pocketbase/generated-types";
import type { QuestionType } from "$lib/types/generic";

export function getQuestions(survey: SurveysResponse): SurveyQuestionsResponse[] {
    // @ts-ignore
    return survey.expand.survey_questions;
}

export function getQuestionType(question: SurveyQuestionsResponse): QuestionType {
    // @ts-ignore
    return question.expand.survey_question_type.name;
}

export function getQuestionIndicator(question: SurveyQuestionsResponse): SurveyIndicatorsResponse {
    // @ts-ignore
    return question.expand.survey_indicator;
}

export function getQuestionTheme(question: SurveyQuestionsResponse): SurveyThemesResponse {
    // @ts-ignore
    return question.expand.survey_indicator.expand.survey_theme;
}

export function getQuestionDomain(question: SurveyQuestionsResponse): SurveyDimensionsResponse {
    // @ts-ignore
    return question.expand.survey_indicator.expand.survey_theme.expand.survey_dimension;
}
