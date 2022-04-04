import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class GeneralService {

  generalURL: string

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    // console.log("New instance of general service");

    this.generalURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // ================================================ ROLE ==========================================================

  getAllRoles(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/role`,
      { headers: httpHeaders })
  }

  // ================================================ GENERAL TYPE ==========================================================

  // Get general type list by type.
  getGeneralTypeByType(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/general-type/type/${type}`,
      { headers: httpHeaders })
  }

  // Get all general types.
  getGeneralTypes(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/general-type`,
      { headers: httpHeaders })
  }

  // Get General Type By ID
  getGeneralType(generaltypeID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/general-type/${generaltypeID}`,
      { headers: httpHeaders })
  }

  // Update General Type
  updateGeneralType(generaltype: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/general-type/${generaltype.id}`,
      generaltype, { headers: httpHeaders })
  }

  // Delete General Type By ID
  deleteGeneralType(generaltypeID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/general-type/${generaltypeID}`,
      { headers: httpHeaders })
  }

  // ================================================ EXAMINATION ==========================================================

  // getExaminationsList gets all examinations list
  getExaminationList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/examination`,
      { headers: httpHeaders });
  }

  // ================================================ UNIVERSITY ==========================================================

  // getUniversityByCountryID gets all university list by country
  getUniversityByCountryID(countryID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/country/${countryID}/university`,
      { headers: httpHeaders });
  }

  /**
   * Gets all universities
   * @returns all universities 
   */
  getAllUniversities(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/university-list`,
      { headers: httpHeaders })
  }

  // ================================================ SOURCE ==========================================================

  // getSources gets all sources
  getSources(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/source`,
      { headers: httpHeaders });
  }

  // ================================================ COLLEGE ==========================================================

  /**
   * Gets all college names list
   * @returns all college names list 
   */
  getAllCollegeNamesList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/college/names`,
      { headers: httpHeaders, observe: "response" })
  }

  getAllCollegesByUniversityID(universityID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/college/${universityID}`,
      { headers: httpHeaders, observe: "response" })
  }

   // Get all Colleges branch list
   getCollegeBranchListWithLimit(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/college/branch/list/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, params: params })
  }

  // Get all Colleges branch list
  getCollegeBranchList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/college/branch/list`,
      { headers: httpHeaders })
  }

  // Get all Colleges
  getColleges(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/college/list`,
      { headers: httpHeaders })
  }

  // ================================================ COURSE ==========================================================

 

  // getCourseList is temp function to get all Course List with Code, Name and Eligibility 
  
  getCourseList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/course-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ SESSION ==========================================================

  // getCourseSessionList returns all the sessions for a specified course
  getCourseSessionList(courseID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/course/${courseID}/session-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ DAY ==========================================================

  // getDaysList returns all the days 
  getDaysList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/day`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ FACULTY ==========================================================

  // getFacultyList is temp function to get all Faculty List. 
  getFacultyList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/faculty-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ EMPLOYEE ==========================================================

  // getEmployeeList is temp function to get all employees List. 
  getEmployeeList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/other-employee-list`,
      { headers: httpHeaders })
  }

  // ================================================ PURPOSE ==========================================================

  // getPurposeList is temp function to get all purposes list. 
  getPurposeListByType(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/purpose/type/${type}`,
      { headers: httpHeaders });
  }

  // ================================================ OUTCOME ==========================================================

  // getOutcomeListByPurpose is temp function to get all outcomes list by purpose
  getOutcomeListByPurpose(purposeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/purpose/${purposeID}/outcome`,
      { headers: httpHeaders });
  }

  // ================================================ STATE ==========================================================

  //  Get State By ID
  getState(stateID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/state/${stateID}`,
      { headers: httpHeaders })
  }

  // Get All State  
  getStates(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/state`,
      { headers: httpHeaders })
  }

  // Delete State By ID
  getStatesByCountryID(countryID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/country/${countryID}/state`,
      { headers: httpHeaders })
  }

  // Update State By ID
  updateState(state: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/state/${state.id}`, state,
      { headers: httpHeaders })
  }

  // Delete State By ID
  deleteState(stateID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/state/${stateID}`,
      { headers: httpHeaders })
  }

  // ================================================ COUNTRY ==========================================================

  //  Get Country By ID
  getCountry(stateID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/state/${stateID}`,
      { headers: httpHeaders })
  }

  // Get All Country  
  getCountries(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/country`,
      { headers: httpHeaders })
  }

  // Update Country By ID
  updateCountry(country: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/country/${country.id}`, country,
      { headers: httpHeaders })
  }

  // Delete Country By ID
  deleteCountry(countryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/country/${countryID}`,
      { headers: httpHeaders })
  }

  // ================================================ COMPANY ==========================================================

  // getCompanies Return All Companies
  getCompanies(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/company`,
      { headers: httpHeaders })
  }

  getCompanyBranchList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/company/branch/list`,
      { headers: httpHeaders, observe: "response" })
  }

  getRequirementList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/company-requirement/list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ================================================ USER ================================================

  // getUserList get list of all salesperson and admin
  getUserList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/user-list`,
      { headers: httpHeaders });
  }

  // GetAllSalesPeople is temp function to get all sales people
  getSalesPersonList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/salesperson-list`,
      { headers: httpHeaders, observe: "response" })
  }

   // Add User.
   addUser(user: IUser): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.post(`${this.generalURL}/user/credential/${credentialID}`, user,
      { headers: httpHeaders })
  }

  // update user.
  updateUser(user: IUser): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/user/${user.id}/credential/${credentialID}`,
      user, { headers: httpHeaders })
  }

  // deletes specified user.
  deleteUser(userID: string): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/user/${userID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  // Get specific users.
  getSpecificUsers(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/user/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ================================================ CREDENTIAL ================================================

  // getUserCredentialList get list of all salesperson and admin
  getUserCredentialList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/user/credential-list`,
      { headers: httpHeaders });
  }

  getAllEmployeeCredentials(limit: number, offset: number, param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/all-employees/limit/${limit}/offset/${offset}`,
      { params: param, headers: httpHeaders, observe: "response" })
  }

  // Get credenyial list by role.
  getCredentialListByRole(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/credential-by-role`,
      { params: params, headers: httpHeaders })
  }

  // ================================================ SPEAKER ================================================

  // Gets speaker list.
  getSpeakerList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/speaker`,
      { headers: httpHeaders})
  }

  // ================================================ FEEDBACK ================================================

  // returns all the feedback questions
  getAllFeedbackQuestion(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feedback-question`,
      { headers: httpHeaders, observe: "response" })
  }

  // returns all the feedback questions with limit and offset
  getFeedbackQuestions(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feedback-question/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // returns feedback questions by the specified type
  getFeedbackQuestionByType(type: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feedback-question/type/${type}`,
      { headers: httpHeaders, observe: "response" })
  }

  // adds feedback questions
  addFeedbackQuestion(feedbackQuestion: IFeedbackQuestion): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(`${this.generalURL}/feedback-question/credential/${credentialID}`,
      feedbackQuestion, { headers: httpHeaders })
  }

  // updates feedback question and its options
  updateFeedbackQuestion(feedbackQuestion: IFeedbackQuestion): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/feedback-question/${feedbackQuestion.id}/credential/${credentialID}`,
      feedbackQuestion, { headers: httpHeaders })
  }

  // deletes feedback questions
  deleteFeedbackQuestion(feedbackQuestionID: string): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/feedback-question/${feedbackQuestionID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  // returns all the feedback options for specified question
  getFeedbackOptions(questionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feedback-question/${questionID}/option`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ NEXT ACTION ================================================

  getNextActionTypeList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/next-action-type`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ TARGET COMMUNITY ================================================

  // Get all target community functions by department id.
  getTargetCommunityFunctionByDepartemnt(departmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/target-community-function/department/${departmentID}`,
      { headers: httpHeaders })
  }

  // ================================================ PROJECT ================================================

  getListOfProjects(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/project-list`,
      { headers: httpHeaders, observe: "response" })
  }

  getListOfSubProjects(projectID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/project/${projectID}/sub-project-list`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ EMPLOYEE ================================================

  getAllEmployees(limit: number, offset: number, param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/other-employee/limit/${limit}/offset/${offset}`,
      { params: param, headers: httpHeaders, observe: "response" })
  }

  addEmployee(employee: IEmployee): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.post(`${this.generalURL}/other-employee/credential/${credentialID}`, employee,
      { headers: httpHeaders, observe: "response" })
  }

  updateEmployee(employee: IEmployee): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.put(`${this.generalURL}/other-employee/${employee.id}/credential/${credentialID}`, employee,
      { headers: httpHeaders, observe: "response" })
  }

  deleteEmployee(employeeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.delete(`${this.generalURL}/other-employee/${employeeID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getAllEmployeeList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/all-employee-list`,
      { headers: httpHeaders, observe: "response" })
  }

  addSupervisor(supervisor: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(`${this.generalURL}/supervisor`, supervisor,
      { headers: httpHeaders, observe: "response" })
  }
  deleteSupervisor(supervisorID: string, employeeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/supervisor/${supervisorID}/employee/${employeeID}`,
      { headers: httpHeaders })
  }

  getDirectReports(supervisorID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/direct-reports/supervisor/${supervisorID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ BATCH ================================================

  getBatchList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch-list`,
      { params: params, headers: httpHeaders});
  }

  // ================================================ FEELING ================================================

  getFeelingList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feeling-list`,
      { headers: httpHeaders, observe: "response" })
  }

  getFeelingLevelList(feelingID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/feeling/${feelingID}/level`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ DEPARTMENT ================================================

  // Get department list.
  getDepartmentList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/department-list`,
      { params: params, headers: httpHeaders })
  }

  // ================================================ DESIGNATION ================================================

  // getDesignations Return All Designations
  getDesignations(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/designation-list`,
      { headers: httpHeaders })
  }

  // Designation By ID
  getDesignation(designationID): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/designation/${designationID}`,
      { headers: httpHeaders })
  }

  // ================================================ TECHNOLOGY ================================================

  // Get Technologies
  getTechnologies(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/technology-list`,
      { headers: httpHeaders, params: params })
  }

  // ================================================ PROGRAMMING LANGUAGE ================================================

  // Get programming language list.
  getProgrammingLanguageList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/programming-language`,
      { headers: httpHeaders })
  }

  // ================================================ BLOG TOPIC ================================================

  // Get blog topic list.
  getBlogTopicList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/blog-topic`,
      { headers: httpHeaders })
  }

  // ================================================ DEGREE ================================================

  // Get all Degree
  getDegrees(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/degree-list`,
      { headers: httpHeaders })
  }

  // Get Degree By ID
  getDegree(degreeID): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/degree/${degreeID}`,
      { headers: httpHeaders })
  }

  // ================================================ SPECIALIZATION ================================================

  // getSpecializationsByDegree get all specializations by dgeree id
  getSpecializations(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/specialization`,
      { headers: httpHeaders, observe: "response" });
  }

  getSpecializationByDegreeID(degreeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/specialization/degree/${degreeID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // ================================================ PROGRAMMING QUESTION ================================================

  // Returns all the programming questions with limit and offset.
  getProgrammingQuestionList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.generalURL}/programming-question`,
      { params: params, headers: httpHeaders })
  }

  // ================================================ PROGRAMMING QUESTION TYPE ================================================

  // Returns programming question list.
  getProgrammingQuestionTypeList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/programming-question-type`,
      { params: params, headers: httpHeaders})
  }

  // // ================================================ ADMIN RESOURCE API ================================================

  // getAllResources(limit: number, offset: number, params?: HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.get(`${this.generalURL}/resource/limit/${limit}/offset/${offset}`,
  //     { params: params, headers: httpHeaders, observe: "response" })
  // }

  // getResourcesByFileType(limit: number, offset: number, fileType: string, params?: HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.get(`${this.generalURL}/resource/file-type/${fileType}/limit/${limit}/offset/${offset}`,
  //     { params: params, headers: httpHeaders, observe: "response" })
  // }

  // getResourceCount(param: HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.get(`${this.generalURL}/resource/file-type/count`,
  //     { params: param, headers: httpHeaders, observe: "response" })
  // }

  // addResource(resource: IResource): Observable<any> {
  //   let credentialID: string = this.localService.getJsonValue('credentialID')
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.post(`${this.generalURL}/resource/credential/${credentialID}`, resource,
  //     { headers: httpHeaders, observe: "response" })
  // }

  // updateResource(resource: IResource): Observable<any> {
  //   let credentialID: string = this.localService.getJsonValue('credentialID')
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.put(`${this.generalURL}/resource/${resource.id}/credential/${credentialID}`, resource,
  //     { headers: httpHeaders, observe: "response" })
  // }

  // deleteResource(resourceID: string): Observable<any> {
  //   let credentialID: string = this.localService.getJsonValue('credentialID')
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.delete(`${this.generalURL}/resource/${resourceID}/credential/${credentialID}`,
  //     { headers: httpHeaders, observe: "response" })
  // }

  // ================================================ SALARY TRENDS ================================================

  getSalaryTrend(limit: number, offset: number, param?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.generalURL}/salary-trend/limit/${limit}/offset/${offset}`,
      { params: param, headers: httpHeaders, observe: "response" })
  }

  addSalaryTrend(salaryTrend: ISalaryTrend): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(`${this.generalURL}/salary-trend/credential/${credentialID}`, salaryTrend,
      { headers: httpHeaders, observe: "response" })
  }

  updateSalaryTrend(salaryTrend: ISalaryTrend): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.put(`${this.generalURL}/salary-trend/${salaryTrend.id}/credential/${credentialID}`, salaryTrend,
      { headers: httpHeaders, observe: "response" })
  }

  deleteSalaryTrend(salaryTrendID: string): Observable<any> {
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.delete(`${this.generalURL}/salary-trend/${salaryTrendID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

}

export interface IAcademics {
  id: string
  degree: IDegree
  college: string
  passout: number
  percentage: number
  degreeID: string
}

export interface IDegree {
  id: string
  name: string
}

export interface ISpecialization {
  id?: string
  branchName: string
  degreeID: string
  isVisible: boolean
}

export interface ITechnologies {
  id: string
  language: string
}

export interface ICountry {
  id?: string
  name: string
}

export interface IState {
  id?: string
  countryID: string
  name: string
}

export interface IDesignation {
  id?: string
  position: string
}

export interface IFeedbackQuestion {
  id: string
  type: string
  question: string
  hasOptions: boolean
  order: number
  maxScore?: number
  keyword: string
  options: IFeedbackOptions[]
  feedbackQuestionGroup?: IFeedbackQuestionGroup
}

export interface IFeedbackOptions {
  id: string
  questionID: string
  key: number
  order: number
  value: string
}

export interface ITargetCommunity {
  id?: string
  targetType: string
  estimatedTargetDate: string
  achievedTargetDate: string
  targetCount: number
  targetAchieved: boolean
}

export interface IDepartment {
  id?: string
  name: string
  roleID: string
  role?: IRole
}

export interface IProject {
  id?: string
  name: string
  subProjects?: IProject[]
}

export interface IEmployee {
  id?: string
  firstName: string
  lastName: string
  email: string
  contact: string
  dateOfBirth?: string
  dateOfJoining?: string
  technologies: ITechnologies[]
  resume?: string
  isActive: boolean
  type: string
  address?: string
  city?: string
  pinCode?: number
  country: ICountry
  state: IState
}

export interface IUser {
  id?: string
  code?: string
  firstName?: string
  lastName?: string
  email?: string
  contact?: string
  dateOfBirth?: string
  dateOfJoining?: string
  role?: IRole
  address?: IAddress
  resume?: string
}

export interface IAddress {
  address?: string
  city?: string
  pinCode?: string
  country?: ICountry
  countryID?: string
  stateID?: string
  state?: IState
}

export interface IRole {
  id?: string
  roleName: string
  level?: string
}

export interface IFeeling {
  id?: string
  feelingName: string
  feelingLevels: IFeelingLevel[]
}

export interface IFeelingLevel {
  id?: string
  description: string
  levelNumer: number
  feelingID: string
}

export interface ITargetCommunityFunction {
  id?: string
  functionName: string
  departmentID: string
  department: IDepartment
}

// export interface IResource {
//   id?: string
//   resourceType: string
//   fileType: string
//   resourceURL: string
//   resourceName: string
//   description?: string
// }

export interface IFileTypeList {
  id?: string
  type: string
  value: string
  key: number
  count: number
}

export interface ISalaryTrend {
  id?: string
  companyRating: number
  date: string
  minimumExperience: number
  maximumExperience: number
  minimumSalary: number
  maximumSalary: number
  medianSalary: number
  technology: ITechnologies
  designation: IDesignation
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}

export interface IFeedbackQuestionGroup {
  id?: string
  groupName: string
  groupDescription: string
  order: number
  feedbackQuestions?: IFeedbackQuestion[]
}