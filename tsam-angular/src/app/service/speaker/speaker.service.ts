import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class SpeakerService {

  speakerUrl: string
  httpHeaders: HttpHeaders

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.speakerUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/speaker`
   }

  // Get all speakers.
  getAllSpeakers(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.speakerUrl}/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new speaker.
  addSpeaker(speaker: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.speakerUrl}/credential/${credentialID}`, 
    speaker, { headers: httpHeaders });
  }

  // Update speaker.
  updateSpeaker(speaker: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.put<any>(`${this.speakerUrl}/${speaker.id}/credential/${credentialID}`, 
    speaker, { headers: httpHeaders });
  }

  // Delete speaker.
  deleteSpeaker(speakerID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.delete<any>(`${this.speakerUrl}/${speakerID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }
}

export interface ISpeaker {
  id?: string

	// Related table IDs.
  designationID: string

	// Other fields.
  firstName: string
  lastName: string
  company: string
  experienceInYears: number
}

export interface ISpeakerDTO {
  id?: string

	// Related tables.
  designation: IDesignation

	// Other fields.
  firstName: string
  lastName: string
  company: string
  experienceInYears: number
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