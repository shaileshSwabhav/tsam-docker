import { IEligibility } from "../list/eligibility";
import { ITechnology } from "../list/technology";
import { ICourseSession } from "./session";

export interface ICourse {
    id: string
    name: string
    code: string
    courseType: string
    courseLevel: string
    description: string
    price: number
    durationInMonths: number
    totalHours: number
    sessions: ICourseSession[]
    technologies: ITechnology[]
    eligibility: IEligibility
    brochure: string
    logo: string
    totalModules: number
    totalTopics: number
    totalConcepts: number
    totalAssignments: number
  }