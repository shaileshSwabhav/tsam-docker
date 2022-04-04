import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { IProgrammingQuestion, IProgrammingQuestionOption, IProgrammingQuestionTalentAnswerQuestionDTO, IProgrammingQuestionType } from '../programming-question/programming-question.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProgrammingQuestionTalentAnswerService {

  answerURL: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.answerURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-talent-answer`
  }

  // Get all programming question talent answers.
  getAllAnswers(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.answerURL}/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Get programming question talent answer by id.
  getAnswer(answerID: string, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.answerURL}/${answerID}`, 
    {params: params, headers: httpHeaders});
  }

  // Add new programming question talent answer.
  addAnswer(answer: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.answerURL}/credential/${credentialID}`, 
    answer, { headers: httpHeaders });
  }

  // Add new solution is viewed.
  addSolutionIsViewed(solutionIsViewd: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-solution-is-viewed/credential/${credentialID}`, 
    solutionIsViewd, { headers: httpHeaders });
  }

  // Update programming question talent answer.
  updateAnswer(answer: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.answerURL}/${answer.id}/credential/${credentialID}`,
     answer, { headers: httpHeaders });
  }

  // Update isCorrect of programming question talent answer.
  updateAnswerScore(answer: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.answerURL}/score/${answer.id}/credential/${credentialID}`,
     answer, { headers: httpHeaders });
  }

  // Delete programming question talent answer.
  deleteAnswer(answerID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.answerURL}/${answerID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }
}

export interface IProgrammingQuestionTalentAnswer {
  id?: string
  answer: string
  score: number

	// Related model IDs.
  programmingQuestionID: string      
  talentID: string                   
  programmingQuestionOptionID: string
  programmingLanguageID: string
  programmingQuestionTypeID: string

  // Flags.
  isCorrect: boolean
}

export interface IProgrammingQuestionTalentAnswerDTO {
  id?: string
  answer: string
  score: number
  totalAnswers: number
  totalNotChecked: number

	// Related models.
  programmingQuestion: IProgrammingQuestionTalentAnswerQuestionDTO      
  talent: IProgrammingQuestionTalentAnswerTalentDTO                   
  programmingQuestionOption: IProgrammingQuestionOption
  programmingLanguage: IProgrammingLanguage
  programmingQuestionType: IProgrammingQuestionType

  // Flags.
  isCorrect: boolean
  solutonIsViewed: boolean
}

export interface IProgrammingQuestionTalentAnswerTalentDTO {
  id?: string
  firstName: string
  lastName: string
}

export interface IProgrammingQuestionTalentAnswerWithFullQuestionDTO {
  id?: string
  answer: string
  score: number

	// Related models.
  programmingQuestion: IProgrammingQuestion   
  programmingLanguage: IProgrammingLanguage
  talent: IProgrammingQuestionTalentAnswerTalentDTO                   

  // Flags.
  isCorrect: boolean
}

export interface IProgrammingQuestionSolutionIsViewed {
  id?: string
  programmingQuestionID: string      
  talentID: string  
}


export interface IProgrammingLanguage {
  id?: string,
  name: string
  rating: number
}