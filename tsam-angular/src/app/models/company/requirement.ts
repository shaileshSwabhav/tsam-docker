import { IDegree } from "../general/degree";
import { IDesignation } from "../general/designation";
import { ISalaryTrend } from "../general/salary_trend";
import { IUniversity } from "../general/university";
import { ITechnology } from "../list/technology";
import { ITalent } from "../talent/talent";
import { IJobLocation } from "./job_location";

export interface IRequirement {
    id?: string
  
    // Address.
    jobLocation: IJobLocation
  
    // Multiple.
    qualifications: IDegree[]
    universities: IUniversity[]
    technologies: ITechnology[]
    selectedTalents: ITalent[]
    designation: IDesignation
  
    // Single objects.
    salesPersonID: string
    companyBranchID: string
    designationID: string
  
    // Other fields.
    isActive: boolean
    code: string
    talentRating: string
    personalityType: string
    minimumExperience: number
    maximumExperience: number
    jobRole: string
    jobDescription: string
    jobRequirement: string
    jobType: string
    // packageOffered: number
    requiredBefore: string
    requiredFrom: string
    vacancy: number
    comment: string
    termsAndConditions: string
  
    // Criteria
    criteria1: string
    criteria2: string
    criteria3: string
    criteria4: string
    criteria5: string
    criteria6: string
    criteria7: string
    criteria8: string
    criteria9: string
    criteria10: string
  
    // Rating
    rating: number
    
    // Salary trend
    salaryTrend: ISalaryTrend
  }