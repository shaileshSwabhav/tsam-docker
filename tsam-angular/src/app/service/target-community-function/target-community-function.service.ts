import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class TargetCommunityFunctionService {

  targetCommunityFunctionUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.targetCommunityFunctionUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/target-community-function`
  }

  // Get all target community functions.
  getAllTargetCommunityFunctions(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.targetCommunityFunctionUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new target community function.
  addTargetCommunityFunction(targetCommunityFunction: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.targetCommunityFunctionUrl}`, 
      targetCommunityFunction, { headers: httpHeaders });
  }

  // Update target community function.
  updateTargetCommunityFunction(targetCommunityFunction: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.targetCommunityFunctionUrl}/${targetCommunityFunction.id}`, 
      targetCommunityFunction, { headers: httpHeaders });
  }

  // Delete target community function.
  deleteTargetCommunityFunction(targetCommunityFunctionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.targetCommunityFunctionUrl}/${targetCommunityFunctionID}`, 
      { headers: httpHeaders });
  }
}
