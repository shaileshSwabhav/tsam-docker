import { IResource } from "../resource/resource";
import { IModuleTopic } from "./module_topic";

export interface IModule {
    id?: string
    moduleName: string
    logo?: string
    moduleTopics?: IModuleTopic[]
    resources?: IResource[]
    batchTopicAssignment: any
    totalTopics: number
    totalSubTopics: number
    totalProgrammingQuestions: number

    // UI fields.
    showDetails: boolean
    isResources: boolean
    // close: string
    // totalQuestions: number
  }