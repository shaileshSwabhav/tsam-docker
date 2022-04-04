import { IConceptModules } from "./concept_module";

export interface IProgrammingConcept{
    id?: string
    name: string 
    difficultyClass?: string
    complexity: number
    levelName: string
    modules: IConceptModules[]
  }