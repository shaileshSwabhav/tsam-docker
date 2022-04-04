import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { IProgrammingQuestion } from '../programming-question/programming-question.service';
import { IResource } from '../resource/resource.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProgrammingAssignmentService {
  
  assignmentURL: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) { 
    this.assignmentURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-assignment`
  }

  // Add new programming-assignment.
  addAssignment(assignment: IProgrammingAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.assignmentURL}`, 
    assignment, { headers: httpHeaders });
  }

  // Update programming-assignment.
  updateAssignment(assignment: IProgrammingAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.assignmentURL}/${assignment.id}`, 
    assignment, { headers: httpHeaders });
  }

  // Delete programming-assignment.
  deleteAssignment(assignmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.assignmentURL}/${assignmentID}`,
		{ headers: httpHeaders });
  }

	getProgrammingAssignments(params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.assignmentURL}`, 
		{ headers: httpHeaders, observe: "response", params: params });
	}

	getProgrammingAssignmentList(params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.assignmentURL}-list`, 
		{ headers: httpHeaders, observe: "response", params: params });
	}
}

export interface IProgrammingAssignment {
	id?: string
	title: string
	taskDescription: string
	timeRequired: string
	complexityLevel: number
	score: number
	additionalComments: string
	sourceURL: string
	source?: string
  programmingAssignmentType: string
	programmingQuestion: IProgrammingQuestion[]
	programmingAssignmentSubTask: IProgrammingAssignmentSubTask[]
}

export interface IProgrammingAssignmentSubTask {
	id?: string
	programmingAssignmentID: string
  resourceID: string
  resource: IResource
  description: string
}

export interface IProgrammingAssignmentList {
	id?: string
	title: string
}
