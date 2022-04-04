import { IProgrammingConcept } from 'src/app/models/programming/concept';

export interface IModuleConcept {
    id?: string
    programmingConcept: IProgrammingConcept
    programmingConceptID: string
    moduleID: string
    level: number
    score?: number
}