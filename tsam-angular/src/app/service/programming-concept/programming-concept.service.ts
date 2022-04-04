import { HttpHeaders, HttpClient, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';
import { IProgrammingConcept } from 'src/app/models/programming/concept';

@Injectable({
  providedIn: 'root'
})
export class ProgrammingConceptService {

  conceptURL: string
  tenantURL: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.conceptURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-concept`
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // Get all concepts.
  getAllConcepts(params?: HttpParams): Observable<HttpResponse<IProgrammingConcept[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<IProgrammingConcept[]>(`${this.conceptURL}s`,
      { params: params, headers: httpHeaders, observe: "response" });
  }
  // // Get all parent concepts.
  // getAllParentConcepts(params?:HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this._http.get(`${this.conceptUrl}`, 
  //   {params: params, headers: httpHeaders});
  // }

  // Get concept by id.
  getConcept(conceptID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.conceptURL}/${conceptID}`,
      { params: params, headers: httpHeaders });
  }

  // Add new concept.
  addConcept(concept: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.conceptURL}/credential/${credentialID}`,
      concept, { headers: httpHeaders });
  }

  // Update concept.
  updateConcept(concept: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.conceptURL}/${concept.id}/credential/${credentialID}`,
      concept, { headers: httpHeaders });
  }

  // Delete concept.
  deleteConcept(conceptID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.conceptURL}/${conceptID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }
}

export interface IConceptModules {
  id?: string
  programmingConceptID: string
  moduleID: string
  level: number
}


