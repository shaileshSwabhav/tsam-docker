import { IProgrammingLanguage } from "./language";
import { IProgrammingQuestionOption } from "./question_option";
import { IProgrammingQuestionSolutionDTO } from "./question_solution";
import { IProgrammingQuestionTestCase } from "./question_test_case";
import { IProgrammingQuestionType } from "./question_type";

export interface IProgrammingQuestion {
    id?: string
    label: string
    question: string
    inputFormat: string        
    example: string   
    constraints: string        
    comment: string    
    hasOptions: boolean
    isActive: boolean
    level: number
    levelName: string
    levelClass: string
    score: number
    timeRequired: number
    programmingQuestionTypes: IProgrammingQuestionType[]
    options: IProgrammingQuestionOption[]
    testCases: IProgrammingQuestionTestCase[]
    isAnswered: boolean
    answer: string
    programmingQuestionOptionID: string
    attemptedByCount: number
    solvedByCount: number
    successRatio: string
    solutionCount: number
    solutions: IProgrammingQuestionSolutionDTO[]
    programmingLanguage: IProgrammingLanguage
    solutonIsViewed: boolean
    hasAnyTalentAnswered: boolean
    isLanguageSpecific: boolean
    programmingLanguages: IProgrammingLanguage[]
  
    // flags
    isMarked: boolean
  }