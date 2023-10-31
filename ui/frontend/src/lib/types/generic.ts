export type QuestionType = "number" | "bool";

export type QuestionUi = {
    questionId: string;
    question: string;
    type: QuestionType;
    indicator: string;
    theme: string;
    dimension: string;
    answer: number;
};
