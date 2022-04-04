import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class SpecializationService {

  specializationUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.specializationUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/specialization`
  }

  // Get all specializations.
  getAllSpecializations(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.specializationUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new specialization.
  addSpecialization(specialization: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.specializationUrl}`, 
      specialization, { headers: httpHeaders });
  }

  // Update specialization.
  updateSpecialization(specialization: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.specializationUrl}/${specialization.id}`, 
      specialization, { headers: httpHeaders });
  }

  // Delete specialization.
  deleteSpecialization(specializationID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.specializationUrl}/${specializationID}`, 
      { headers: httpHeaders });
  }
}
