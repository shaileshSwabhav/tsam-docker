import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProgrammingLanguageService {

  languageUrl: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.languageUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-language`
  }

  // Get all programming languages.
  getAllProgrammingLanguages(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.languageUrl}/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new programming language.
  addProgrammingLanguage(language: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.languageUrl}/credential/${credentialID}`, 
    language, { headers: httpHeaders });
  }

  // Update programming language.
  updateProgrammingLanguage(language: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.languageUrl}/${language.id}/credential/${credentialID}`, 
    language, { headers: httpHeaders });
  }

  // Delete programming language.
  deleteProgrammingLanguage(designationID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.languageUrl}/${designationID}/credential/${credentialID}`, 
      { headers: httpHeaders });
  }
}

export interface IProgrammingLanguage {
  id?: string,
  name: string
  rating: number
}
