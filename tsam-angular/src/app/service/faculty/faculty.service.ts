import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { IFeedbackOptions, IFeedbackQuestion } from '../general/general.service';
import { UserloginService } from '../login/userlogin.service';

@Injectable({
  providedIn: 'root'
})
export class FacultyService {

  private facultyURL: string
  httpHeaders: HttpHeaders

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
  ) {
    this.facultyURL = `${this._constant.BASE_URL}/tenant/${this._constant.TENANT_ID}/faculty`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue("credentialID")
  }

  // Get All Faculty 
  // getAllFaculty(limit: number, offset: number, params?): Observable<any> {
  //   return new Observable<any>((observer) => {
  //     let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //     this._http.get(`${this.facultyURL}/limit/${limit}/offset/${offset}`, { params: params, headers: httpHeaders, observe: 'response' }).
  //       subscribe(data => {
  //         observer.next(data)
  //       }, (error) => {
  //         observer.error(error)
  //       })
  //   })
  // }

  // Get All Faculty 
  // getSearchedFaculty(conditions: any, limit: number, offset: number): Observable<any> {
  //   return this._http.post(`${this.facultyURL}/search/limit/${limit}/offset/${offset}`, conditions,
  //   { headers: this.httpHeaders, observe: "response" });

    // return new Observable<any>((observer) => {
    //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    //   this._http.post(`${this.facultyURL}/search/limit/${limit}/offset/${offset}`, conditions, { headers: httpHeaders, observe: 'response' }).
    //     subscribe(data => {
    //       observer.next(data)
    //     }, (error) => {
    //       observer.error(error)
    //     })
    // })
  // }

  getAllFaculty(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.facultyURL}`, 
    { params: params, headers: httpHeaders, observe: "response" });

  }

  // Get Faculty By ID.
  getFaculty(facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.facultyURL}/${facultyID}`, 
    { headers: httpHeaders, observe: "response" });
  }

  // Get faculty for batch.
  getFacultyForBatch(batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this._constant.BASE_URL}/tenant/${this._constant.TENANT_ID}/batch/${batchID}/faculty`, 
    { headers: httpHeaders, observe: "response" });
  }

  // Add New Faculty
  addFaculty(faculty: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.post(`${this.facultyURL}`, faculty,
    { headers: httpHeaders, observe: "response" });

  }

  // Update Faculty
  updateFaculty(faculty: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.put(`${this.facultyURL}/${faculty.id}`, faculty,
    { headers: httpHeaders, observe: "response" });

  }

  // Delete Faculty
  deleteFaculty(facultyid: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.delete(`${this.facultyURL}/${facultyid}`,
    { headers: httpHeaders, observe: "response" });

  }

  getFacultyCredentialList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.facultyURL}/credential/list`, 
    { headers: httpHeaders, observe: "response" })
  }

  // ============================================================= FACULTY ASSESSMENT =============================================================
  
  getFacultyAssessment(facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.facultyURL}/${facultyID}/assessment`, 
    { headers: httpHeaders, observe: "response" })
  }

  addFacultyAssessment(facultyAssessment: IFacultyAssessment[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.post(`${this.facultyURL}/assessments`, facultyAssessment,
    { headers: httpHeaders, observe: "response" })
  }

  deleteFacultyAssessment(facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.delete(`${this.facultyURL}/${facultyID}/assessment/${credentialID}`, 
    { headers: httpHeaders, observe: "response" })
  }

}

export interface IFacultyCredentialList {
  id: string
  firstName: string
  lastName: string
  facultyID: string
}

export interface IFacultyAssessment {
  id?: string
  credentialID: string
  facultyID?: string
  questionID: string
  optionID: string
  answer: string
  assessmentType: string
  faculty: any
  credential: any
  question: IFeedbackQuestion
  option: IFeedbackOptions
  averageScore: number
}