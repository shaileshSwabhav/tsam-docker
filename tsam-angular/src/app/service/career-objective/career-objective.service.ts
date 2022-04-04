import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class CareerObjectiveService {

  careerObjectiveUrl: string
  httpHeaders: HttpHeaders
  // credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    this.careerObjectiveUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/career-objective`
    this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
  }

  // Get all career objectives.
  getAllCareerObjectives(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.careerObjectiveUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new career objective.
  addCareerObjective(careerObjective: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.careerObjectiveUrl}`, careerObjective, 
      { headers: httpHeaders });
  }

  // Update career objective.
  updateCareerObjective(careerObjective: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.careerObjectiveUrl}/${careerObjective.id}`, 
      careerObjective, { headers: httpHeaders });
  }

  // Delete career objective.
  deleteCareerObjective(careerObjectiveID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.careerObjectiveUrl}/${careerObjectiveID}`, 
      { headers: httpHeaders });
  }
}

export interface ICareerObjectiveCourse{
  id?: string
  careerObjectiveID: string
  courseID: string         
  order: number            
  technicalAspect: string  
}

export interface ICareerObjective{
  id?: string
  name: string
  courses: ICareerObjectiveCourse[]          
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}
