import { IDesignation } from "./designation";
import { ITechnologies } from "./technology";

export interface ISalaryTrend {
    id?: string
    companyRating: number
    date: string
    minimumExperience: number
    maximumExperience: number
    minimumSalary: number
    maximumSalary: number
    medianSalary: number
    technology: ITechnologies
    designation: IDesignation
}