import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class FeedbackService {

  feedbackURL: string
  generalURL: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.feedbackURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/feedback-question`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
    
  }

  // Returns all the feedback questions with limit and offset.
  getFeedbackQuestions(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.feedbackURL}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
  
  // getFeedbackQuestionGroupByType(type: string): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

  //   return this._http.get(`${this.feedbackURL}/feedback-question-group/type/${type}`,
  //   { headers: httpHeaders, observe: "response" })
  // }

  // returns feedback questions by the specified type
  getFeedbackQuestionByType(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.feedbackURL}/type/${type}`, 
    { headers: httpHeaders, observe: "response" })
  }

  // Adds feedback questions.
  addFeedbackQuestion(feedbackQuestion: IFeedbackQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    
    return this._http.post(`${this.feedbackURL}/credential/${credentialID}`,
      feedbackQuestion, { headers: httpHeaders })
  }

  // Updates feedback question and its options.
  updateFeedbackQuestion(feedbackQuestion: IFeedbackQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.put(`${this.feedbackURL}/${feedbackQuestion.id}/credential/${credentialID}`,
      feedbackQuestion, { headers: httpHeaders })
  }

  // Updates the status of the feedback question to active/inactive
  updateFeedbackQuestionStatus(feedbackQuestion: IFeedbackQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.put(`${this.feedbackURL}/${feedbackQuestion.id}/status/credential/${credentialID}`,
      feedbackQuestion, { headers: httpHeaders })
  }

  // Deletes feedback questions.
  deleteFeedbackQuestion(feedbackQuestionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.delete(`${this.feedbackURL}/${feedbackQuestionID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }


}

export interface IFeedbackQuestion {
  id: string
  type: string
  question: string
  hasOptions: boolean
  order: number
  maxScore?: number
  keyword: string
  isActive: boolean
  options: IFeedbackOptions[]
  feedbackQuestionGroup?: IFeedbackQuestionGroup
}

export interface IFeedbackOptions {
  id: string
  questionID: string
  key: number
  order: number
  value: string
}

export interface IFeedbackQuestionGroup {
  id?: string
  groupName: string
  groupDescription: string
  order: number
  answer?: number
  maxScore: number
  minScore: number
  feedbackQuestions?: IFeedbackQuestion[]
}

export interface IResourceCount {
  fileType: string
  totalCount: number
}