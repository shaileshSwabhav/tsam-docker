import { HttpClient, HttpHeaders, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { IDay } from '../batch/batch.service';
import { Constant } from '../constant';
import { IProgrammingConcept } from 'src/app/models/programming/concept';
import { IProgrammingQuestion } from '../programming-question/programming-question.service';
import { IResource } from '../resource/resource.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ModuleService {

  private moduleURL: string
  private topicURL: string

  constructor(
    private http: HttpClient,
    private localService: LocalService,
    private constant: Constant
  ) {
    this.moduleURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/modules`
    this.topicURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/topic`
  }

  addModule(module: IModule): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.moduleURL}`, module,
      { headers: httpHeaders, observe: "response" })
  }

  updateModule(module: IModule): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.moduleURL}/${module.id}`, module,
      { headers: httpHeaders, observe: "response" })
  }

  deleteModule(moduleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.moduleURL}/${moduleID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getModule(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.moduleURL}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ========================================== MODULE RESOURCE ==========================================

  addResource(moduleID: string, resource: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.moduleURL}/${moduleID}/resources`, resource,
      { headers: httpHeaders, observe: "response" })
  }

  deleteResource(moduleID: string, resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.moduleURL}/${moduleID}/resources/${resourceID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getModuleResource(moduleID: string): Observable<HttpResponse<IResource[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get<IResource[]>(`${this.moduleURL}/${moduleID}/resources`,
      { headers: httpHeaders, observe: "response" })
  }

  // ========================================== MODULE TOPICS ==========================================


  addModuleTopic(moduleID: string, moduleTopic: IModuleTopic): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.moduleURL}/${moduleID}/topic`, moduleTopic,
      { headers: httpHeaders, observe: "response" })
  }

  addModuleTopics(moduleID: string, moduleTopics: IModuleTopic[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.moduleURL}/${moduleID}/topics`, moduleTopics,
      { headers: httpHeaders, observe: "response" })
  }

  updateModuleTopic(moduleID: string, moduleTopic: IModuleTopic): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.moduleURL}/${moduleID}/topic/${moduleTopic.id}`, moduleTopic,
      { headers: httpHeaders, observe: "response" })
  }

  deleteModuleTopic(moduleID: string, moduleTopicID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.moduleURL}/${moduleID}/topic/${moduleTopicID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getModuleTopic(moduleID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.moduleURL}/${moduleID}/topic`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ========================================== TOPIC ASSIGNMENT ==========================================

  addTopicProgrammingQuestion(topicID: string, topicAssignment: ITopicProgrammingQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.topicURL}/${topicID}/programming-question`, topicAssignment,
      { headers: httpHeaders, observe: "response" })
  }

  updateTopicProgrammingAssignment(topicID: string, topicAssignment: ITopicProgrammingQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.topicURL}/${topicID}/programming-question/${topicAssignment.id}`, topicAssignment,
      { headers: httpHeaders, observe: "response" })
  }

  deleteTopicProgrammingAssignment(moduleID: string, topicAssignmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.moduleURL}/${moduleID}/topic/${topicAssignmentID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getTopicProgrammingQuestions(topicID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    console.log("topicID -> ", topicID);

    return this.http.get(`${this.topicURL}/${topicID}/programming-question`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getTopicProgrammingQuestionsList(topicID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.topicURL}/${topicID}/programming-question-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
}

export interface IModule {
  id?: string
  moduleName: string
  logo?: string
  moduleTopics?: IModuleTopic[]
  resources?: IResource[]
  close: string

  // UI fields.
  totalSubTopics: number
  totalQuestions: number
}

// export interface IModuleTopic {
//   id?: string
//   topicName: string
//   totalTime: string
//   order: number
//   studentOutput: string
//   subTopics?: IModuleTopic[]
//   topicID: string
//   moduleID?: string
//   module?: IModule
//   topicProgrammingConcept?: ITopicProgrammingConcept[]
//   topicProgrammingQuestions?: ITopicProgrammingQuestion[]

//   // extra
//   programmingConceptIDs?: string[]
//   batchTopicAssignment: IBatchTopicAssignment[]
//   batchSessionTopic: any[]

//   // flags
//   isAddQuestionClick?: boolean
//   isEditQuestionClick?: boolean
//   showDetails?: boolean
//   isMarked: boolean
// }

export interface ITopicProgrammingQuestion {
  id?: string
  programmingQuestionID: string
  topicID?: string
  isActive: boolean
  order?: number
  topic?: IModuleTopic
  programmingQuestion?: IProgrammingQuestion
  programmingConcept?: ITopicProgrammingConcept[]

  // flag
  isMarked?: boolean
  checked?: string
}

export interface ITopicProgrammingConcept {
  id?: string
  programmingConceptID?: string
  programmingConcept?: IProgrammingConcept
  moduleTopicID?: string
  moduleTopic?: IModuleTopic
}

export interface IModuleTiming {
  id?: string
  batchID?: string
  moduleID?: string
  facultyID?: string
  dayID?: string
  day: IDay
  fromTime?: string
  toTime?: string

  // flag
  isMarked: boolean
}