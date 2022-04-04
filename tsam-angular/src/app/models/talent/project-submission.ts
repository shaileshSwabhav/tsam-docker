import { IProjectSubmissionUpload } from "./project-submission-upload";
import { ITalent } from "./talent";

export interface ITalentProjectSubmission {
    id?: string
    talentID: string
    talent?: ITalent
    batchProjectID: string
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
    isLatestSubmission?: boolean
    websiteLink?: string
    projectUpload?: string
    projectSubmissionUpload: IProjectSubmissionUpload[]
    facultyID?: string
    assignedDate?: string
    dueDate?: string
}