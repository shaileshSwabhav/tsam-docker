import { IModuleTopic } from "../course/module_topic";
import { IProgrammingQuestion } from "../programming/question";

export interface IBatchTopicAssignment {
    id?: string
    batchID?: string
    topicID?: string
    topic?: IModuleTopic
    programmingQuestion?: IProgrammingQuestion
    programmingQuestionID?: string  
    dueDate?: string
    assignedDate?: string
    totalTopics?: number
    moduleID?: string
  
    // flags
    showDetails?: boolean
    isMarked?: boolean
    checked?: string
  }