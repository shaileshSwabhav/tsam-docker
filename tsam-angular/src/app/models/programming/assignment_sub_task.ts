import { IResource } from "../resource/resource";

export interface IProgrammingAssignmentSubTask {
    id?: string
    programmingAssignmentID: string
    resourceID: string
    resource: IResource
    description: string
}