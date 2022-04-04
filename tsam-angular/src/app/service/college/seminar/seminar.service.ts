import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../../constant';
import { LocalService } from '../../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class SeminarService {

  generalURL: string

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
  ) { 
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`
  }

  //*******************************SEMINAR API CALLS********************************************* 

  // Add one seminar.
  addSeminar(seminar: any): Observable<any> {
    // let talentJson:string = JSON.stringify(seminar)
    // console.log("talent" + talentJson)
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/seminar/credential/${credentialID}`, 
      seminar, { headers: httpHeaders })
  }

  // Get all seminars.
  getSeminars(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/seminar/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" })
  }

  // Update seminar.
  updateSeminar(seminar: any): Observable<any> {
    // let talentJson: string = JSON.stringify(seminar)
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/seminar/${seminar.id}/credential/${credentialID}`, 
    seminar, { headers: httpHeaders })
  }

  // Delete seminar.
  deleteSeminar(id: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/seminar/${id}/credential/${credentialID}`, 
    { headers: httpHeaders })
  }

  // *******************************SEMINAR TOPIC API CALLS********************************
  // Get all topics by seminar.
  getTopicsBySeminar(seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/seminar-topic/seminar/${seminarID}`, 
    { headers: httpHeaders });
  }

  // Add topic.
  addTopic(topic: any, seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/seminar-topic/seminar/${seminarID}/credential/${credentialID}`, 
    topic, { headers: httpHeaders});
  }

  // Update topic.
  updateTopic(topic: any, seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/seminar-topic/${topic.id}/seminar/${seminarID}/credential/${credentialID}`, 
    topic, { headers: httpHeaders });
  }

  // Delete topic.
  deleteTopic(topicID: any, seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/seminar-topic/${topicID}/seminar/${seminarID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  // *******************************STUDENT API CALLS********************************
  // Get all students by seminar.
  getStudentsBySeminar(seminarID: any, limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/student/seminar/${seminarID}/limit/${limit}/offset/${offset}`, 
    { params: params, headers: httpHeaders, observe: "response"  });
  }

  // Add student.
  addStudent(student: any, seminarID: any): Observable<any> {
    // let talentJson:string = JSON.stringify(student);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/student/seminar/${seminarID}/credential/${credentialID}`, 
    student, { headers: httpHeaders});
  }

  // Update student.
  updateStudent(student: any, seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/student/${student.seminarTalentRegistrationID}/seminar/${seminarID}/credential/${credentialID}`, 
    student, { headers: httpHeaders });
  }

  // Delete student.
  deleteStudent(studentID: any, seminarID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/student/${studentID}/seminar/${seminarID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  // Update multiple fields of students.
  updateMultipleStudent(updateMultipleStudent: IUpdateMultipleStudent) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })  
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/student/seminar/${updateMultipleStudent.seminarID}/credential/${credentialID}`, 
    updateMultipleStudent, { headers: httpHeaders });
  }
}

export interface ISeminar {
  id?: string

	// Maps.
  collegeBranches: ICollegeBranch[]  
  salesPeople: ISalesperson[]        
  speakers: ISpeaker[]        

	// Flags.
  isActive: boolean

  // Other fields.
  seminarName: string   
  description: string    
  location: string       
  code: string           
  seminarDate: string    
  fromTime: string       
  toTime: string         
  registrationLink: string
}

export interface ISeminarDTO {
  id?: string

	// Maps.
  collegeBranches: ICollegeBranch[]  
  salesPeople: ISalesperson[]   
  speakers: ISpeaker[]        

	// Flags.
  isActive: boolean

  // Other fields.
  seminarName: string   
  description: string    
  location: string       
  code: string           
  seminarDate: string    
  fromTime: string       
  toTime: string         
  registrationLink: string
  totalRegisteredStudents: number
  totalVisitedStudents: number
}

export interface ICollegeBranch {
  id?: string
  branchName: string
  code: string
}

export interface ISalesperson {
  id?: string
  firstName: string
  lastName: string
}

export interface ISpeaker {
  id?: string
  firstName: string
  lastName: string
  company: string
  experienceInYears: number
  designationID: string
}

export interface ISeminarTopic {
  id?: string

	// Related table IDs.
  speakerID: string
  seminarID: string

	// Other fields.
  topicName: string
  date: string
  fromTime: number
  toTime: string
  description: string
}

export interface ISeminarTopicDTO {
  id?: string

	// Related tables.
  speaker: ISpeaker
  seminarID: string

	// Other fields.
  topicName: string
  date: string
  fromTime: number
  toTime: string
  description: string
}

export interface IUpdateMultipleStudent {
  hasVisited: boolean
  seminarID: string
  seminarTalentRegistrationIDs: string[]
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}
