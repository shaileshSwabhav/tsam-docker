import { IBatchProject } from "./batch/project";
import { ITalentProjectSubmission } from "./talent/project-submission";

export interface IProjectSubmission {
    project: IBatchProject
    talentSubmission: Map<string, ITalentProjectSubmission>
}