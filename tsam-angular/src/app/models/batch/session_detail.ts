import { IModule } from "../course/module";
import { IBatchSessionTopic } from "./session_topic";
import { IBatch } from "./batch";
import { IBatchSessionPrerequisite } from "./session_prerequisites";

export interface IBatchSessionDetail {
    id: string
    batch?: IBatch
    batchSessionPrerequisite?: IBatchSessionPrerequisite
    date: string
    isAttendanceMarked: boolean
    isCompleted: boolean
    isSessionTaken: boolean
    module: IModule[]
    pendingModule: IModule[]
    batchSessionTopic: IBatchSessionTopic[]
    totalModuleSubTopic: number
    isResources: boolean
  }