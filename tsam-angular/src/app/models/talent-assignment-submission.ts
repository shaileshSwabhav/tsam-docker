import { IBatchTopicAssignment } from "./batch-topic-assignment";

export interface ITalentAssignmentSubmission {
    id?: string
    talentID: string
    talent?: any
    batchTopicAssignmentID: string
    batchTopicAssignment?: IBatchTopicAssignment
    isAccepted?: boolean
    isChecked?: boolean
    acceptanceDate?: string
    facultyRemarks?: string
    facultyVoiceNote?: string
    score?: number
    solution?: string
    submittedOn?: string | Date
    githubURL?: string
    voiceNote?: string
    isVoiceNoteVisible?: boolean
    talentConceptRatings?: ITalentConceptRating
    isLatestSubmission?: boolean
    assignmentSubmissionUploads: IAssignmentSubmissionUpload[]
}

export interface IAssignmentSubmissionUpload {
    id?: string
    assignmentSubmissionID?: string
    imageURL?: string
    description?: string
}

export interface ITalentAllSubmissions {
    talentID: string
    allSubmissions: ITalentAssignmentSubmission[]
}

export interface ITalentConceptRating {
    programmingConceptModuleID: string
    talentID: string
    talentSubmissionID: string
    score: number
}