import { IProgrammingProject } from "../programming/project";
import { ITalentProjectSubmission } from "../talent/project-submission";
import { IBatch } from "./batch"

export interface IBatchProject {
    id?: string
    batchID: string
    batch: IBatch
    programmingProjectID: string
    dueDate: string
    assignedDate: string
    submissions: ITalentProjectSubmission[]
    programmingProject: IProgrammingProject
    submittedOn?: string | Date
}