import { IProgrammingConcept } from "../programming/concept";
import { IModuleTopic } from "./module_topic";


export interface ITopicProgrammingConcept {
    id?: string
    programmingConceptID?: string
    programmingConcept?: IProgrammingConcept
    moduleTopicID?: string
    moduleTopic?: IModuleTopic
  }