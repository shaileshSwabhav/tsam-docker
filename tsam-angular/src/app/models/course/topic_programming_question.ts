import { IProgrammingQuestion } from "../programming/question";
import { IModuleTopic } from "./module_topic";
import { ITopicProgrammingConcept } from "./topic_programming_concept";

export interface ITopicProgrammingQuestion {
  id?: string
  programmingQuestionID: string
  topicID?: string
  isActive: boolean
  order?: number
  topic?: IModuleTopic
  programmingQuestion?: IProgrammingQuestion
  programmingConcept?: ITopicProgrammingConcept[]

  // flag
  isMarked?: boolean
  checked?: string
}