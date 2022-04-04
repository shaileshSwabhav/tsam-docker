import { HttpClient, HttpHeaders, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IModuleConcept } from 'src/app/models/module-concept';
import { Constant } from '../constant';
import { DagModelItem } from '../dag-manager/dag-manager.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ConceptModuleService {

  conceptModuleURL: string
  tenantURL: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.conceptModuleURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-concept-module`
  }

  // Get all concept modules for concept tree.
  GetAllModuleProgrammingConceptsForConceptTree(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.conceptModuleURL}-tree`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get all concept modules.
  getAllModuleProgrammingConcepts(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.conceptModuleURL}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get all concept modules for assignment.
  // bta - batchTopicAssignment
  getConceptModulesForAssignment(btaID: string, params?: HttpParams): Observable<HttpResponse<IModuleConcept[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<IModuleConcept[]>(`${this.tenantURL}/batch-topic-assignment/${btaID}/programming-concepts-modules`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get all concept modules for talent score.
  getAllConceptModulesForTalentScore(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.tenantURL}/talent/${talentID}/programming-concept-module`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get one concept modules .
  getConceptModule(conceptModuleID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.conceptModuleURL}/${conceptModuleID}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Add new concept modules.
  addConceptModules(conceptModules: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.conceptModuleURL}`,
      conceptModules, { headers: httpHeaders })
  }

  // Update concept modules.
  updateConceptModules(moduleID: string, conceptModules: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.conceptModuleURL}/module/${moduleID}`,
      conceptModules, { headers: httpHeaders })
  }

  // Delete concept module.
  deleteConceptModule(conceptModuleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.conceptModuleURL}/${conceptModuleID}`,
      { headers: httpHeaders })
  }
}

export interface WorkflowItem extends DagModelItem {
  [x: string]: any;
  conceptID: string
  id: string
  toBeRemoved: boolean
  parentConceptIDs: string[]
  currentParentConceptID?: string
  isConceptIDDisabled: boolean
  showParentConcepts: boolean
}

export interface WorkflowItemScore extends DagModelItem {
  [x: string]: any;
  conceptID: string
  id: string
  score: number
  conceptName: string
  showParentConcepts: boolean
  complexityName: string
}


