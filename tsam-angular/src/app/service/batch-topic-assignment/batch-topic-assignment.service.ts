import { HttpHeaders, HttpClient, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';
import { IBatchTopicAssignment } from 'src/app/models/batch-topic-assignment';
import { ITalentAssignmentSubmission } from 'src/app/models/talent-assignment-submission';


@Injectable({
  providedIn: 'root'
})
export class BatchTopicAssignmentService {

  batchUrl: string
  assignmentSubmissionURL: string
  httpHeaders: HttpHeaders
  tenantUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.batchUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`
    this.assignmentSubmissionURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/topic-assignment`
    this.tenantUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }


  // will get all topic assignments with submissions
  // future change can be using qp to decide whether to load submissions. #niranjan
  getTopicAssignmentWithSubmissions(batchID: string, params?: HttpParams): Observable<HttpResponse<IBatchTopicAssignment[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get<IBatchTopicAssignment[]>(`${this.tenantUrl}/batch/${batchID}/topic-assignments`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Add batch topic assignment.
  addBatchTopicAssignment(batchID: string, topicID: string, batchTopicAssignment: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.batchUrl}/${batchID}/topic/${topicID}/assignment`,
      batchTopicAssignment, { headers: httpHeaders })
  }

  // Get batch assignment List
  getBatchAssignment(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.batchUrl}/${batchID}/talent/${talentID}/topic-assignment-submissions`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get talent assignment List
  getTalentAssignmentSubmission(batchTopicAssignmentID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.assignmentSubmissionURL}/${batchTopicAssignmentID}/talent/${talentID}/talent-submission`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  deleteBatchTopicAssignment(batchID: string, topicID: string, topicAssignmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.batchUrl}/${batchID}/topic/${topicID}/assignment/${topicAssignmentID}`,
      { headers: httpHeaders })
  }


  updateTopicAssignment(batchID: string, topicID: string, topicAssignmentID: string, assignment: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.batchUrl}/${batchID}/topic/${topicID}/assignment/${topicAssignmentID}`,
      assignment, { headers: httpHeaders })
  }

  getAllTopicAssignments(batchID:string,params?:any) :Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.batchUrl}/${batchID}/assignments`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
  //********************************************* TALENT ASSIGNMENT SUBMISSION FUNCTIONS ************************************************************

  // Assignment submissions of a particular talent.
  // needs qp and batchID change. #Niranjan
  getTalentAssignmentSubmissions(topicAssignmentID: string, talentID: string, params?: HttpParams): Observable<HttpResponse<ITalentAssignmentSubmission[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // tenant/{id}/topic-assignments/{id}/talents/{id}/submissions
    return this._http.get<ITalentAssignmentSubmission[]>(`${this.tenantUrl}/topic-assignments/${topicAssignmentID}/talents/${talentID}/submissions`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get batch topic assignment list.
  getAllBatchTopicAssignmentsForTalent(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/talent/${talentID}/talent-submission`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Add talent assignment submission.
  addAssignmentSubmission(batchTopicAssignmentID: string, talentID: string, assignmentSubmission: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.assignmentSubmissionURL}/${batchTopicAssignmentID}/talent/${talentID}/talent-submission`,
      assignmentSubmission, { headers: httpHeaders, observe: "response" })
  }

  scoreTalentAssignment(assignmentID: string, talentID: string, submissionID: string, scoredAssignment: any): Observable<any> {
    // /tenant/{id}/topic-assignments/{id}/talents/{id}/submissions/{id}/score
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.tenantUrl}/topic-assignments/${assignmentID}/talents/${talentID}/submissions/${submissionID}/score`,
      scoredAssignment, { headers: httpHeaders })
  }
}