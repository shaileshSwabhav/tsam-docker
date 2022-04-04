import { IProgrammingLanguage } from "./language";

export interface IProgrammingQuestionSolutionDTO{
    id?: string
    solution: string
    programmingQuestionID: string
    programmingLanguage: IProgrammingLanguage
  }