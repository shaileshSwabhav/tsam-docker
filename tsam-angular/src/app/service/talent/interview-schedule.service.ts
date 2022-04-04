import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { UrlWithStringQuery } from 'url';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class InterviewScheduleService {

  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService
  ) {
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`;
    }

  
// *********************************************Talent Interview Schedule API calls********************************************
  // Get all interview schedules by talent
  getInterviewSchedulesByTalent(talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview-schedule/talent/${talentID}`, { headers: httpHeaders });
  }

  // Add interview schedule
  addInterviewSchedule(interviewSchedule: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.post(`${this.generalURL}/interview-schedule/talent/${talentID}/credential/${credentialID}`, 
    interviewSchedule, { headers: httpHeaders});
  }

  getInterviewSchedule(interviewScheduleID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview-schedule/${interviewScheduleID}/talent/${talentID}`, 
    { headers: httpHeaders });
  }

  // Update talent interview schedule
  updateInterviewSchedule(interviewSchedule: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.put<any>(`${this.generalURL}/interview-schedule/${interviewSchedule.id}/talent/${talentID}/credential/${credentialID}`, 
    interviewSchedule, { headers: httpHeaders });
  }

  // Delete talent intreview schedule
  deleteInterviewSchedule(interviewScheduleID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.delete<any>(`${this.generalURL}/interview-schedule/${interviewScheduleID}/talent/${talentID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

// **************************************Talent Interview API calls************************************************
  // Get all interview by interview schedule
  getInterviewsByInterviewSchedule(interviewScheduleID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview/interview-schedule/${interviewScheduleID}`, 
    { headers: httpHeaders });
  }

  getAllInterviewersList(): Observable<any>{
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview/interviewers-list`, 
    { headers: httpHeaders });
  }

  // Add interview
  addInterview(interview: any, interviewScheduleID: any): Observable<any> {
    // let talentJson:string = JSON.stringify(interview);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.post(`${this.generalURL}/interview/interview-schedule/${interviewScheduleID}/credential/${credentialID}`, 
    interview, { headers: httpHeaders});
  }

  getInterview(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview`, 
    { params: params, headers: httpHeaders, observe: "response" })
  }

  // Update interview
  updateInterview(interview: any, interviewScheduleID: any): Observable<any> {
    // let talentJson:string = JSON.stringify(interview);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.put<any>(`${this.generalURL}/interview/${interview.id}/interview-schedule/${interviewScheduleID}/credential/${credentialID}`, 
    interview, { headers: httpHeaders });
  }

  // Delete interview
  deleteInterview(interviewID: any, interviewScheduleID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue("credentialID");
    return this._http.delete<any>(`${this.generalURL}/interview/${interviewID}/interview-schedule/${interviewScheduleID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

// **************************************Talent Interview Round API calls************************************************
  // Get all interview rounds
  getInterviewRounds(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/interview-round`, 
    { headers: httpHeaders });
  }

}

export interface IInterview{
  id?: string
  talentID: string
  roundID: string
  rating: number
  comment: string
  status: string
}

