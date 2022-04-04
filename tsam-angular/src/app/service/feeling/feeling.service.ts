import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class FeelingService {

  feelingUrl: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    this.feelingUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/feeling`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
  }

  // Get all feelings.
  getAllFeelings(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.feelingUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new feeling.
  addFeeling(feeling: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.feelingUrl}`, 
      feeling, { headers: httpHeaders });
  }

  // Update feeling.
  updateFeeling(feeling: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.feelingUrl}/${feeling.id}`, 
      feeling, { headers: httpHeaders });
  }

  // Delete feeling.
  deleteFeeling(feelingID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.feelingUrl}/${feelingID}`, 
      { headers: httpHeaders });
  }
}

export interface IFeelingLevel{
  id?: string
  feelingID: string         
  description: string            
  levelNumber: number  
}

export interface IFeeling{
  id?: string
  feelingName: string
  feelingLevels: IFeelingLevel[]          
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}