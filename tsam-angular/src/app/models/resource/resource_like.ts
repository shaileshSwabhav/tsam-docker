import { ICredential } from "../general/credential";
import { IResource } from "./resource";

export interface IResourceLike {
    id?: string
    resourceID?: string
    credentialID?: string
    resource: IResource
    credential: ICredential
    isLiked: boolean
  }