import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class DegreeService {

  degreeUrl: string
  httpHeaders: HttpHeaders
  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.degreeUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/degree`
  }

  // Get all degrees.
  getAllDegrees(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.degreeUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new degree.
  addDegree(degree: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.degreeUrl}`, 
      degree, { headers: httpHeaders });
  }

  // Update degree.
  updateDegree(degree: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.degreeUrl}/${degree.id}`, 
      degree, { headers: httpHeaders });
  }

  // Delete degree.
  deleteDegree(degreeID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.degreeUrl}/${degreeID}`, 
      { headers: httpHeaders });
  }
}

export interface IDegree {
  id?: string
  name: string
}

export interface ISpecialization {
  id?: string
  branchName: string
  degreeID: string
}

export interface ISpecializationDTO {
  id?: string
  branchName: string
  degreeID: string
  degree: IDegree
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}

