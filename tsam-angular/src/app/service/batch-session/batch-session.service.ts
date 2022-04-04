import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';
import { Observable } from 'rxjs';
import { IBatch, IBatchSession } from '../batch/batch.service';
import { IModule } from '../module/module.service';

@Injectable({
  providedIn: 'root'
})
export class BatchSessionService {

  batchUrl: string
  httpHeaders: HttpHeaders

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.batchUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`
  }

  // Add batch session plan.
  addBatchSessionPlan(batchID: string, sessionPlanList: any[], params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.batchUrl}/${batchID}/generate-session`,
      sessionPlanList, { headers: httpHeaders, params: params });
  }

  skipPendingSession(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.batchUrl}/${batchID}/skip-session`, null,
      { headers: httpHeaders, params: params });
  }

  updateBatchSession(batchID: string, batchSession: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.batchUrl}/${batchID}/sessions`,
      batchSession, { headers: httpHeaders, params: params, observe: "response" });
  }

  deleteSessionPlan(batchID: string, facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.batchUrl}/${batchID}/faculty/${facultyID}/sessions`,
      { headers: httpHeaders });
  }

  // Get batch session plan.
  getBatchSessionPlan(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/session-plan`,
      { params: params, headers: httpHeaders });
  }

  // Get batch session dates.
  getBatchSessionDates(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/session`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get batch sessions counts.
  getBatchSessionsCounts(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/session-plan-counts`,
      { headers: httpHeaders })
  }

  // Get batch session with topic and sub topics names.
  getBatchSessionWithTopicNameList(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/session-topic-name-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  //Mark Session As Complete.
  markSubTopicAsComplete(batchID: string, subTopicID: string, subTopic: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put(`${this.batchUrl}/${batchID}/sub-topic/${subTopicID}/mark-as-complete`, subTopic,
      { headers: httpHeaders })
  }

  addSessionPlanAttendance(batchID: string, sessionID: string, talent: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.batchUrl}/${batchID}/batch-session/${sessionID}/talent-attendance`,
      talent, { headers: httpHeaders });
  }

  //Get All Batch Sessions
  getAllBatchSessionPlan(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/all-session-plan`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  //Add PreRequisites
  addPreRequisite(batchID: string, batchSessionID: string, preRequisite: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.batchUrl}/${batchID}/batch-session/${batchSessionID}/pre-requisite`,
      preRequisite, { headers: httpHeaders, observe: "response" });
  }

  //Get PreRequisites
  getPreRequisite(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.batchUrl}/${batchID}/pre-requisite`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  //Update PreRequisites
  updatePreRequisite(batchID: string, batchSessionID: string, prerequisiteID: string, preRequisite: any): Observable<any> {
    console.log("preReq", preRequisite)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.batchUrl}/${batchID}/batch-session/${batchSessionID}/pre-requisite/${prerequisiteID}`,
      preRequisite, { headers: httpHeaders, observe: "response" });
  }

  //Delete PreRequisites
  deletePreRequisite(batchID: string, batchSessionID: string, prerequisiteID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.batchUrl}/${batchID}/batch-session/${batchSessionID}/pre-requisite/${prerequisiteID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getSessionModuleTopics(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.batchUrl}/${batchID}/module-topics`,
      { headers: httpHeaders, params: params, observe: "response" })
  }



  getModuleTopics(batchID: string, moduleID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.batchUrl}/${batchID}/module/${moduleID}/topics`,
      { headers: httpHeaders, params: params, observe: "response" })
  }

}
export interface IBatchSessionDetail {
  id: string
  batch?: IBatch
  batchSessionPrerequisite?: IBatchSessionPrerequisite
  date: string
  isAttendanceMarked: boolean
  isCompleted: boolean
  isSessionTaken: boolean
  module: IModule
  pendingModule: IModule
}

export interface IBatchSessionPrerequisite {
  id: string
  batchID?: string
  batchSessionID?: string
  prerequisite: string
}


