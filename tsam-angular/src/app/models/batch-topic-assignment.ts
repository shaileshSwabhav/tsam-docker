import { IProgrammingQuestion } from "../service/programming-question/programming-question.service"
import { ITalentAssignmentSubmission } from "./talent-assignment-submission"

export interface IBatchTopicAssignment {
    id?: string
    batchID: string
    programmingQuestionID: string
    programmingQuestion?: IProgrammingQuestion
    topicID?: string
    moduleID?: string
    dueDate?: string
    batchSessionID?: string
    assignedDate?: string
    showDetails?: boolean
    submissions?: ITalentAssignmentSubmission[]
    submittedOn?: string | Date
}