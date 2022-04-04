import { IProgrammingAssignment } from "../programming/assignment";

export interface ICourseProgrammingAssignment {
    id?: string
    courseID: string
    courseSessionID: string
    programmingAssignmentID?: string
    programmingAssignment: IProgrammingAssignment
    order: number
    isActive: boolean
    isMarked?: boolean
    showAssignmentDetails?: boolean
  }