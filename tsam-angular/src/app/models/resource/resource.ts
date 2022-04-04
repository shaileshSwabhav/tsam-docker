import { IResourceLike } from "./resource_like";

export interface IResource {
    id?: string
    resourceType: string
    fileType: string
    resourceURL: string
    previewURL: string
    resourceName: string
    description?: string
    totalDownload: number
    totalLike: number
    resourceLike: IResourceLike
  
    isBook: string
    author: string
    publication: string
  }