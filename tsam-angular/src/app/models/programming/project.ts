import { ITechnology } from "../list/technology";
import { IResource } from "../resource/resource";

export interface IProgrammingProject {
    id?: string
    projectName: string
    description: string
    projectType: string
    code: string
    isActive: boolean
    complexityLevel: number
    requiredHours: number
    sampleUrl?: string
    resourceType: string
    document?: string
    score?: number
    technologies: ITechnology[]
    resources: IResource[]
  }