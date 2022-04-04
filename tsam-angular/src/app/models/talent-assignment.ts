import { IProgrammingQuestion } from "../service/programming-question/programming-question.service"
import { ITalentAssignmentSubmission } from "../service/talent/talent.service"

export interface ITalentAssignment {
    id?: string
    batchID: string
    programmingQuestion: IProgrammingQuestion
    dueDate?: string
    assignedDate: string
    showDetails?: boolean
    submission?: ITalentAssignmentSubmission
    submittedDate?: string | Date
    score?: number
}