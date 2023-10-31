// Municipalities types
export interface Municipality {
    id: string;
    name: string;
    region: string;
    canton: string;
    district: string;
}

export interface ParticipatingMunicipality extends Municipality {
    avgQuantitativeRating: number;
    avgQualitativeRating: number;
}

export interface MunicipalitiesStats {
    totalMunicipalities: number;
    avgQuantitativeRating: number;
    avgQualitativeRating: number;
}

export interface MunicipalityDetails extends Municipality {
    quantitativeData: Data[];
    qualitativeData: Data[];
}

// Survey types
export interface Survey {
    id: string;
    name: string;
}

export interface ParticipatingSurvey extends Survey {
    totalParticipatingMunicipalities: number;
    avgQualitativeRating: number;
}

export interface SurveyStats {
    totalSurveys: number;
    avgSurveysPerDay: number;
    avgQualitativeRating: number;
}

export interface Question {
    id: string;
    name: string;
    description: string;
    dimension: string;
    theme: string;
    indicator: string;
}

export interface SurveyDetails extends Survey {
    questions: Question[];
    quantitativeData: Data[];
    participatingMunicipalities: ParticipatingMunicipality[];
}

// Generic types
export interface Data {
    dimension: string;
    themes: Theme[];
}

export interface Theme {
    name: string;
    indicators: Indicator[];
}

export interface Indicator {
    name: string;
    rating: number;
}

// Responses

// @path: /api/municipalities
export interface ApiMunicipalitiesResponse {
    municipalities: ParticipatingMunicipality[];
}

// @path: /api/municipalities/stats
export interface ApiMunicipalitiesStatsResponse {
    municipalitiesStats: MunicipalitiesStats;
}

// @path: /api/municipalities/:id
export interface ApiMunicipalityDetailsResponse {
    municipalityDetails: MunicipalityDetails;
}

// @path: /api/surveys
export interface ApiSurveysResponse {
    surveys: ParticipatingSurvey[];
}

// @path: /api/surveys/stats
export interface ApiSurveysStatsResponse {
    surveysStats: SurveyStats;
}

// @path: /api/surveys/:id
export interface ApiSurveyDetailsResponse {
    surveyDetails: SurveyDetails;
}
