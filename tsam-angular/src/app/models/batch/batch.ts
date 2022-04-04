import { IRequirement } from "../company/requirement";
import { ICourse } from "../course/course";
import { IEligibility } from "../list/eligibility";
import { IFaculty } from "../list/faculty";
import { ISalesperson } from "../list/sales_person";
import { IBatchSession } from "./session";
import { IBatchTiming } from "./timing";

export interface IBatch {
    id?: string
    batchName: string
    code: string
    startDate?: string
    endDate?: string
    totalStudents?: number
    totalIntake: number
    batchStatus?: string
    isActive: boolean
    eligibility: IEligibility
    course: ICourse  
    salesPerson: ISalesperson
    faculty: IFaculty
    sessions: IBatchSession[]
    batchTimings: IBatchTiming[]
    requirement: IRequirement
    isB2B: boolean
    brochure?: string
    logo?: string
    batchObjective: string
    totalSessionCount: number
    completedSessionCount: number
  }