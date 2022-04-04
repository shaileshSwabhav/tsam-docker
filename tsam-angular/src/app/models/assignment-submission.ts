import { IBatchTopicAssignment } from "./batch-topic-assignment"
import { ITalentAssignmentSubmission } from "./talent-assignment-submission"

export interface IAssignmentSubmission {
    assignment: IBatchTopicAssignment
    talentSubmission: Map<string, ITalentAssignmentSubmission>
}
