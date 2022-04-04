import { HttpClient, HttpHeaders, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { Constant } from '../constant';
import { ICourse, IResource } from '../course/course.service';
import { IModule } from '../module/module.service';
import { IProgrammingAssignment } from '../programming-assignment/programming-assignment.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class CourseModuleService {

  private courseURL: string

  constructor(
    private http: HttpClient,
    private localService: LocalService,
    private constant: Constant,
  ) { 
    this.courseURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/course`
  }

  // Add New Course Module
  addCourseModule(courseID: string, courseModule: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.courseURL}/${courseID}/course-module`,
      courseModule, { headers: httpHeaders })
  }

  // Update Course Module
  updateCourseModule(courseID: string, courseModule: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.courseURL}/${courseID}/course-module/${courseModule.id}`, 
      courseModule, { headers: httpHeaders })
  }

  // Delete Course Module
  deleteCourseModule(courseID: string, courseModuleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.courseURL}/${courseID}/course-module/${courseModuleID}`,
      { headers: httpHeaders })
  }

  getCourseModules(courseID: any, params?: HttpParams): Observable<HttpResponse<ICourseModule[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get<ICourseModule[]>(`${this.courseURL}/${courseID}/course-module`, 
      { headers: httpHeaders, observe: "response", params: params })
  }

  // ========================================== COURSE TOPIC ASSIGNMENTS ==========================================

  addCourseTopicAssignment(courseID: string, courseProgrammingAssignmnet: ICourseProgrammingAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(`${this.courseURL}/${courseID}/programming-assignment`,  courseProgrammingAssignmnet,
      { headers: httpHeaders, observe: "response" })
  }

  updateCourseTopicAssignment(courseID: string, courseProgrammingAssignmnet: ICourseProgrammingAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.courseURL}/${courseID}/programming-assignment/${courseProgrammingAssignmnet.id}`,  
      courseProgrammingAssignmnet, { headers: httpHeaders, observe: "response" })
  }

  deleteCourseTopicAssignment(courseID: string, courseProgrammingAssignmnetID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.courseURL}/${courseID}/programming-assignment/${courseProgrammingAssignmnetID}`,  
      { headers: httpHeaders, observe: "response" })
  }

  getCourseTopicAssignment(courseID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/${courseID}/programming-assignment`, 
      { headers: httpHeaders, observe: "response", params: params })
  }

  getCourseTopicAssignmentList(courseID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/${courseID}/programming-assignment-list`, 
      { headers: httpHeaders, observe: "response", params: params })
  }


}

export interface ICourseProgrammingAssignment {
  id?: string
  courseID: string
  // courseSessionID: string
  moduleTopicID: string
  programmingAssignmentID?: string
  programmingAssignment: IProgrammingAssignment
  order: number
  isActive: boolean
  isMarked?: boolean
  showAssignmentDetails?: boolean
}

export interface ICourseSession {
  id: string
  name: string
  hours: number
  order: number
  studentOutput: string
  sessionID: string
  subSessions: ICourseSession[]
  courseID: string
  resource: IResource[]
  courseProgrammingAssignment: ICourseProgrammingAssignment[]

  // extra fields
  viewSubSessionClicked?: boolean
  cardColumn?: string
  isChecked?: boolean
  showAssignments: boolean
}

export interface ICourseModule {
  id?: string
  courseID?: string
  moduleID?: string
  order: number
  isActive: boolean
  course?: ICourse
  module?: IModule
  isSelected?: boolean
  moduleTopics?: IModuleTopic[]
}