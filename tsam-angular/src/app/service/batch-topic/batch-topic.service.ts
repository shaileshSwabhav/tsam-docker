import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { IBatchList } from '../batch/batch.service';
import { Constant } from '../constant';
import { IProgrammingQuestion } from '../programming-question/programming-question.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BatchTopicService {

  private batchURL: string

  constructor(
    private http: HttpClient,
    private localService: LocalService,
    private constant: Constant,
  ) {
    this.batchURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`
  }

  addBatchTopic(batchID: string, batchTopic: IBatchTopic): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/topic`, batchTopic,
      { headers: httpHeaders, observe: "response" })
  }

  addBatchTopics(batchID: string, batchTopics: IBatchTopic[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/topics`, batchTopics,
      { headers: httpHeaders, observe: "response" })
  }

  updateBatchTopic(batchID: string, batchTopic: IBatchTopic): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchURL}/${batchID}/topic/${batchTopic.id}`, batchTopic,
      { headers: httpHeaders, observe: "response" })
  }

  updateBatchTopics(batchID: string, batchTopic: IBatchTopic[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchURL}/${batchID}/topics`, batchTopic,
      { headers: httpHeaders, observe: "response" })
  }

  deleteBatchTopic(batchID: string, batchTopicID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.batchURL}/${batchID}/topic/${batchTopicID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchModules(batchID: string, params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/module`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getBatchTopic(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
}

export interface IBatchTopic {
  id?: string
  batchID?: string
  batch?: any
  moduleID: string
  moduleTopicID: string
  order: number
  totalTime: string
  isCompleted: boolean
  completedDate: string 
}

export interface IBatchSessionPlan {
  id?: string
  batchID?: string
  batchTopicID?: string
  sessionDate: string
  batch?: IBatchList
  batchTopic?: IBatchTopic
}

export interface IBatchTopicAssignment {
  id?: string
  batchID?: string
  topicID?: string
  topic?: IModuleTopic
  programmingQuestion?: IProgrammingQuestion
  programmingQuestionID?: string  
  dueDate?: string
  assignedDate?: string
  totalTopics?: number
  moduleID?: string

  // flags
  showDetails?: boolean
  isMarked?: boolean
  checked?: string
}
