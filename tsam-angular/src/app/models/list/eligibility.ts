import { ITechnology } from "./technology";

export interface IEligibility {
  id: string
  technologies: ITechnology[]
  studentRating: string
  experience: boolean
  academicYear: string
}