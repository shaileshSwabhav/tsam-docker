import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BatchProjectService {
  // getTalentProjectSubmissions(projectID: string, talentID: string) {
  //   throw new Error('Method not implemented.');
  // }

  batchUrl: string
  projectSubmissionURL: string
  httpHeaders: HttpHeaders
  tenantUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.batchUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`
    this.projectSubmissionURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/project-submission`
    this.tenantUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  //********************************************* TALENT PROJECT SUBMISSION FUNCTIONS ************************************************************

  // Get batch topic project list.
  getAllBatchProjectForTalent(batchID: string, talentID: string,projectID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.batchUrl}/${batchID}/talent/${talentID}/project-submission`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Add talent project submission.
  addProjectSubmission(batchID: string, talentID: string, projectID: string, assignmentSubmission: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.batchUrl}/${batchID}/project/${projectID}/talent/${talentID}/project-submission`,
      assignmentSubmission, { headers: httpHeaders, observe: "response" })
  }

  // get talent project submission with score
  getProjectWithSubmissions(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.tenantUrl}/batch/${batchID}/projects`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  getTalentProjectSubmissions(batchID: string,projectID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // tenant/{id}/project/{id}/talent/{id}/project-submission
    return this._http.get(`${this.tenantUrl}/batch/${batchID}/project/${projectID}/talent/${talentID}/project-submission`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  scoreTalentProject(batchID: string,projectID: string, talentID: string, projectSubmissionID: string, scoredProject: any): Observable<any> {
    // /tenant/{id}/topic-assignments/{id}/talents/{id}/submissions/{id}/score
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.tenantUrl}/batch/${batchID}/project/${projectID}/talent/${talentID}/project-submission/${projectSubmissionID}/score`,
    scoredProject, { headers: httpHeaders })
  }

  //get project rating parameter
  getProjectRatingParameter(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.tenantUrl}/project-rating-parameter`,
      { headers: httpHeaders, observe: "response" })
  }
}
