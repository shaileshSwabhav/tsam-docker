import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';
import { ICredential } from '../target-community/target-community.service';

@Injectable({
  providedIn: 'root'
})
export class ResourceService {

  generalURL: string
  resourceURL: string

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.generalURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.resourceURL = `${this.generalURL}/resource`
  }

  
  // ================================================ RESOURCE API ================================================

  getAllResources(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.resourceURL}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getResourceCount(param: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.resourceURL}/file-type/count`,
      { params: param, headers: httpHeaders, observe: "response" })
  }

  getResourcesList(param?: HttpParams) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.generalURL}/resource-list`,
      { params: param, headers: httpHeaders, observe: "response" })
  }

  addResource(resource: IResource): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    
    return this.http.post(`${this.resourceURL}`, resource,
      { headers: httpHeaders, observe: "response" })
  }

  updateResource(resource: IResource): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    
    return this.http.put(`${this.resourceURL}/${resource.id}`, resource,
      { headers: httpHeaders, observe: "response" })
  }

  deleteResource(resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    
    return this.http.delete(`${this.resourceURL}/${resourceID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ RESOURCE DOWNLOAD API ================================================

  addResourceDownload(resourceDownload: IResourceDownload): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.resourceURL}/download/credential/${credentialID}`, resourceDownload,
      { headers: httpHeaders, observe: "response" })
  }

  getResourceDownloadCount(resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.resourceURL}/${resourceID}/download/count`,
      { headers: httpHeaders, observe: "response" })
  }

  getResourceDownload(resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.resourceURL}/${resourceID}/download`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ RESOURCE LIKE API ================================================
  
  addResourceLike(resourceLike: IResourceLike): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.resourceURL}/like/credential/${credentialID}`, resourceLike,
      { headers: httpHeaders, observe: "response" })
  }
  
  updateResourceLike(resourceLike: IResourceLike): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.resourceURL}/like/${resourceLike.id}/credential/${credentialID}`, resourceLike,
      { headers: httpHeaders, observe: "response" })
  }

  getResourceLike(resourceID: string): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.resourceURL}/${resourceID}/like/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getAllResourceLike(resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.resourceURL}/${resourceID}/like`,
      { headers: httpHeaders, observe: "response" })
  }
}

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

export interface IResourceDownload {
  id?: string
  resourceID: string
  credentialID: string
}

export interface IResourceLike {
  id?: string
  resourceID?: string
  credentialID?: string
  resource: IResource
  credential: ICredential
  isLiked: boolean
}