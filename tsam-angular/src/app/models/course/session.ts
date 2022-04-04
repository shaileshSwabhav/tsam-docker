import { IResource } from "../resource/resource";
import { ICourseProgrammingAssignment } from "./programming_assignment";

export interface ICourseSession {
    id: string
    name: string
    hours: number
    order: number
    studentOutput: string
    sessionID: string
    subSessions: ICourseSession[]
    courseID: string
    resource: IResource[]
    courseProgrammingAssignment: ICourseProgrammingAssignment[]
  
    // extra fields
    viewSubSessionClicked?: boolean
    cardColumn?: string
    isChecked?: boolean
    showAssignments: boolean
  }
  