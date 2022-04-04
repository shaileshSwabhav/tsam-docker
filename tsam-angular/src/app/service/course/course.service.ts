import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { ITechnology } from '../technology/technology.service';
import { IProgrammingAssignment } from '../programming-assignment/programming-assignment.service';


@Injectable({
  providedIn: 'root'
})
export class CourseService {

  private generalURL: string
  private courseURL: string
  credentialID: string
  httpHeaders: HttpHeaders

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.generalURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.courseURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/course`
    // this.credentialID = localService.getJsonValue('credentialID')
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  }

  // Add New Course
  addCourse(course: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.courseURL}`,
      course, { headers: httpHeaders })
  }

  // Update Course
  updateCourse(course: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.courseURL}/${course.id}`, course,
      { headers: httpHeaders })
  }

  // Delete Course
  deleteCourse(courseID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.courseURL}/${courseID}`,
      { headers: httpHeaders })
  }

  // Get All Course
  getCourses(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.courseURL}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // getAllSearchedCourses(conditions, limit: number, offset: number): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

  //   return this.http.post(`${this.courseURL}/search/limit/${limit}/offset/${offset}`,
  //     conditions, { headers: httpHeaders, observe: "response" })
  // }

  // Get Course By Course ID
  getCourse(courseID: any) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.courseURL}/${courseID}`, { headers: httpHeaders })
  }
  
  // getCourseList is temp function to get all Course List with Code, Name and Eligibility 
  getCourseList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/course-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ============================================================= SESSIONS =============================================================

  getSessionsForCourse(courseID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.courseURL}/${courseID}/session`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // getCourseSessionWiseProgrammingAssignment(courseID: string, params?: HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

  //   return this.http.get(`${this.courseURL}/${courseID}/session-assignment`,
  //     { params: params, headers: httpHeaders, observe: "response" })
  // }

  addSession(session: ICourseSession[], courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.courseURL}/${courseID}/sessions`
      , session, { headers: httpHeaders })
  }

  updateSession(session: ICourseSession, courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.courseURL}/${courseID}/session/${session.id}`
      , session, { headers: httpHeaders })
  }

  deleteSession(courseID: string, sessionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.courseURL}/${courseID}/session/${sessionID}`
      , { headers: httpHeaders })
  }

  // getCourseSessionList returns all the sessions for a specified course
  getCourseSessionList(courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/course/${courseID}/session-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ============================================================= RESOURCE =============================================================

  addSessionResource(sessionID: string, resource: IResource): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.generalURL}/session/${sessionID}/resource`
      , resource, { headers: httpHeaders })
  }

  addSessionResources(resources: IResource[], sessionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.generalURL}/session/${sessionID}/resources`
      , resources, { headers: httpHeaders })
  }

  updateSessionResource(sessionID: string, resource: IResource): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.generalURL}/session/${sessionID}/resources`
      , resource, { headers: httpHeaders })
  }

  deleteSessionResources(sessionID: string, resourceID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.generalURL}/session/${sessionID}/resource/${resourceID}`,
      { headers: httpHeaders })
  }

  // ============================================================= COURSE TECHNICAL ASSESSMENT =============================================================

  getAllCourseTechnicalAssessment(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.generalURL}/course-technical-assessment`,
      { headers: httpHeaders, observe: "response" })
  }

  getCourseTechnicalAssessmentForFaculty(facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.generalURL}/course-technical-assessment/faculty/${facultyID}`,
      { headers: httpHeaders, observe: "response" })
  }

  addCourseTechnicalAssessmnets(assessments: ICourseTechnicalAssessment[], facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.generalURL}/course-technical-assessment/faculty/${facultyID}`
      , assessments, { headers: httpHeaders, observe: "response" })
  }

  updateCourseTechnicalAssessmnet(assessment: ICourseTechnicalAssessment, facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.generalURL}/course-technical-assessment/${assessment.id}/faculty/${facultyID}`
      , assessment, { headers: httpHeaders, observe: "response" })
  }

  deleteCourseTechnicalAssessmnet(assessmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.generalURL}/course-technical-assessment/${assessmentID}`
      , { headers: httpHeaders, observe: "response" })
  }

  // ======================================= COURSE DETAILS FOR STUDENT LOGIN =============================================================
  // Get my courses by talent id.
  
  getMyCoursesByTalent(talentID: string, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/talent/${talentID}/my-courses`,
      { params: params, headers: httpHeaders})
  }

  // Get course details by course id.
  getCoursesDetails(courseID: string, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/details/${courseID}`,
      { params: params, headers: httpHeaders})
  }

  // Get course minimum details by course id.
  getCourseMinimumDetails(talentID: string, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/talent/${talentID}`,
      { params: params, headers: httpHeaders})
  }

  // ======================================= COURSE PROGRAMMING ASSIGNMENT =============================================================

  addCourseProgrammingAssignment(courseAssignment: ICourseProgrammingAssignment, courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.courseURL}/${courseID}/programming-assignment`
      , courseAssignment, { headers: httpHeaders, observe: "response" })
  }

  updateCourseProgrammingAssignment(courseAssignment: ICourseProgrammingAssignment, courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.courseURL}/${courseID}/programming-assignment/${courseAssignment.id}`
      , courseAssignment, { headers: httpHeaders, observe: "response" })
  }

  deleteCourseProgrammingAssignment(courseAssignmentID: string, courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.courseURL}/${courseID}/programming-assignment/${courseAssignmentID}`
      , { headers: httpHeaders, observe: "response" })
  }

  getCourseProgrammingAssignmentList(courseID: string, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/${courseID}/programming-assignment-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getCourseProgrammingAssignment(courseID: string, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.courseURL}/${courseID}/programming-assignment`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

}

export interface ICourse {
  id: string
  name: string
  code: string
  courseType: string
  courseLevel: string
  description: string
  price: number
  durationInMonths: number
  totalHours: number
  sessions: ICourseSession[]
  technologies: ITechnology[]
  eligibility: IEligibility
  brochure: string
  logo: string
  totalModules: number
  totalTopics: number
  totalConcepts: number
  totalAssignments: number
}

export interface ICourseList {
  id: string
  code: string
  name: string
  courseType: string
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

export interface IResource {
  id: string
  resourceType: string
  resourceURL: string
  description: string
  moduleTopicID: string
}

export interface IEligibility {
  id: string
  technologies: ITechnology[]
  studentRating: string
  experience: boolean
  academicYear: string
}

export interface ICourseTechnicalAssessment {
  id?: string
  courseID: string
  course: ICourse
  faculty: any
  rating: number
}

export interface ICourseProgrammingAssignment {
  id?: string
  courseID: string
  courseSessionID: string
  programmingAssignmentID?: string
  programmingAssignment: IProgrammingAssignment
  order: number
  isActive: boolean
  isMarked?: boolean
  showAssignmentDetails?: boolean
}

export interface ICourseModule {
  id: string
  moduleName: string
  courseID: string
  order: number
  isActive: boolean
  course: ICourse
  courseSessions: ICourseSession[]
}