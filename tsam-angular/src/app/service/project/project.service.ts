import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProjectService {

  projectUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.projectUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/project`
  }

  // Get all projects.
  getAllProjects(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.projectUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new project.
  addProject(project: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.projectUrl}`, project, 
      { headers: httpHeaders });
  }

  // Update project.
  updateProject(project: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.projectUrl}/${project.id}`, project, 
      { headers: httpHeaders });
  }

  // Delete project.
  deleteProject(projectID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.projectUrl}/${projectID}`, 
      { headers: httpHeaders });
  }
}

export interface IProject{
  id?: string
  name: string         
  projectID: string            
  subProjects: IProject[]  
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}
