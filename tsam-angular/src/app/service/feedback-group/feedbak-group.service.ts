import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { IFeedbackQuestion } from '../feedback/feedback.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class FeedbackGroupService {

  feedbackQuestionGroupUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.feedbackQuestionGroupUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/feedback-question-group`
  }

  // Get all feedbackQuestionGroups.
  getAllFeedbackQuestionGroups(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.feedbackQuestionGroupUrl}/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  getFeedbackQuestionGroupList(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.feedbackQuestionGroupUrl}/type/${type}`, 
    { headers: httpHeaders, observe: "response" });
  }

  // Add new feedbackQuestionGroup.
  addFeedbackQuestionGroup(feedbackQuestionGroup: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.feedbackQuestionGroupUrl}/credential/${credentialID}`, 
    feedbackQuestionGroup, { headers: httpHeaders });
  }

  // Update feedbackQuestionGroup.
  updateFeedbackQuestionGroup(feedbackQuestionGroup: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.feedbackQuestionGroupUrl}/${feedbackQuestionGroup.id}/credential/${credentialID}`, 
    feedbackQuestionGroup, { headers: httpHeaders });
  }

  // Delete feedbackQuestionGroup.
  deleteFeedbackQuestionGroup(feedbackQuestionGroupID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.feedbackQuestionGroupUrl}/${feedbackQuestionGroupID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  getFeedbackQuestionGroupByType(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.feedbackQuestionGroupUrl}/feedback-question/type/${type}`,
    { headers: httpHeaders, observe: "response" })
  }


  getGroupwiseKeywordNames(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.feedbackQuestionGroupUrl}/group-wise-keyword`,
    { headers: httpHeaders, observe: "response" })
  }

}

export interface IFeedbackQuestionGroup {
  id?: string
  groupName: string
  groupDescription: string
  order: number
  feedbackQuestions?: IFeedbackQuestion[]
}
