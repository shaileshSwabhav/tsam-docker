import { IBatchSessionTopic } from "../batch/session_topic";
import { IBatchTopicAssignment } from "../list/batch_topic_assignment";
import { IModule } from "./module";
import { ITopicProgrammingConcept } from "./topic_programming_concept";
import { ITopicProgrammingQuestion } from "./topic_programming_question";

export interface IModuleTopic {
    id?: string
    topicName: string
    totalTime: number
    order: number
    studentOutput: string
    subTopics?: IModuleTopic[]
    topicID: string
    moduleID?: string
    module?: IModule
    topicProgrammingConcept?: ITopicProgrammingConcept[]
    topicProgrammingQuestions?: ITopicProgrammingQuestion[]
  
    // extra
    programmingConceptIDs?: string[]
    batchTopicAssignment: IBatchTopicAssignment[]
    batchSessionTopic: IBatchSessionTopic[]
  
    // flags
    isAddQuestionClick?: boolean
    isEditQuestionClick?: boolean
    isTopicAdded?: boolean
    showDetails?: boolean
    isMarked: boolean
  }