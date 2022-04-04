import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class DesignationService {

  designationUrl: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    this.designationUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/designation`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
  }

  // Get all designations.
  getAllDesignations(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.designationUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new designation.
  addDesignation(designation: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.designationUrl}`, 
     designation, { headers: httpHeaders });
  }

  // Update designation.
  updateDesignation(designation: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.designationUrl}/${designation.id}`,
      designation, { headers: httpHeaders });
  }

  // Delete designation.
  deleteDesignation(designationID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.designationUrl}/${designationID}`, 
      { headers: httpHeaders });
  }
}

export interface IDesignation {
  id?: string
  position: string
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}

