import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { ICourse } from '../course/course.service';
import { IEvent } from '../admin/admin.service';
import { IBatchTopicAssignment } from 'src/app/models/batch-topic-assignment';

@Injectable({
  providedIn: 'root'
})
export class TalentService {

  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService
  ) {
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`
  }

  //*******************************TALENT API CALLS********************************************* 
  // Add new talent.
  addTalent(talent: any): Observable<any> {
    // let talentJson:string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.generalURL}/talent/credential/${credentialID}`,
      talent, { headers: httpHeaders });
  }

  // Add new talent from excel.
  addTalentFromExcel(talent: any): Observable<any> {
    // let talentJson:string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.generalURL}/talent-excel/credential/${credentialID}`,
      talent, { headers: httpHeaders });
  }

  // Add multiple talents.
  addTalents(talents: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')

    return this._http.post<any>(`${this.generalURL}/talents/credential/${credentialID}`,
      talents, { headers: httpHeaders });
  }

  // Get all searched talents.
  getAllSearchedTalents(conditions: any, limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/talent/search/limit/${limit}/offset/${offset}`, conditions,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  // Get all talents.
  getTalents(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  // Get talent.
  getTalent(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/${talentID}`,
      { params: params, headers: httpHeaders });
  }

  // Get eligible talent.
  getEligibleTalent(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/${talentID}/eligible`,
      { params: params, headers: httpHeaders });
  }

  // Update talent.
  updateTalent(talent: any): Observable<any> {
    // let talentJson: string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent/${talent.id}/credential/${credentialID}`, talent,
      { headers: httpHeaders });
  }

  // Update talent.
  updateTalentsSalesPerson(talents: any, salesPersonID): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent/saleperson/${salesPersonID}/credential/${credentialID}`,
      talents, { headers: httpHeaders });
  }

  // Delete talent.
  deleteTalent(id: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/talent/${id}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // *******************************TALENT CALL RECORDS API CALLS********************************
  // Get all call records by talent.
  getCallRecordsByTalent(talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-call-record/talent/${talentID}`,
      { headers: httpHeaders });
  }

  // Add call record.
  addCallRecord(callRecord: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/talent-call-record/talent/${talentID}/credential/${credentialID}`,
      callRecord, { headers: httpHeaders });
  }

  getCallRecord(callRecordID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-call-record/${callRecordID}/talent/${talentID}`,
      { headers: httpHeaders, observe: "response" });
  }

  // Update talent call record.
  updateCallRecord(callRecord: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-call-record/${callRecord.id}/talent/${talentID}/credential/${credentialID}`,
      callRecord, { headers: httpHeaders });
  }

  // Delete talent callrecord.
  deleteCallRecord(callRecordID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/talent-call-record/${callRecordID}/talent/${talentID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // *******************************TALENT LIFETIME VALUE API CALLS********************************
  // Add lifetime value.
  addLifetimeValue(lifetimeValue: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/talent-lifetime-value/talent/${talentID}/credential/${credentialID}`, lifetimeValue,
      { headers: httpHeaders, observe: "response" });
  }

  // Get one talent lifetime value.
  getLifetimeValue(talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-lifetime-value/talent/${talentID}`,
      { headers: httpHeaders });
  }

  // Update talent lifetime value.
  updateLifetimeValue(lifetimeValue: any, talentID: any): Observable<any> {
    // let talentJson: string = JSON.stringify(lifetimeValue);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-lifetime-value/${lifetimeValue.id}/talent/${talentID}/credential/${credentialID}`,
      lifetimeValue, { headers: httpHeaders });
  }

  // Delete talent lifetime value.
  deleteLifetimeValue(lifetimeValueID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/talent-lifetime-value/${lifetimeValueID}/talent/${talentID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // Get all lifetime value reports.
  getLifetimeValueReports(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let loginID: string = this.localService.getJsonValue("loginID");
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-lifetime-value/login/${loginID}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  // *******************************API CALLS FOR FACULTY********************************
  // Get all batches of one talent.
  getBatchListOfOneTalent(talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/${talentID}/batch`,
      { headers: httpHeaders });
  }

  // *******************************API CALLS FOR NEXT ACTION********************************
  // Get all next actions by talent.
  getAllTalentNextActions(talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/${talentID}/talent-next-action`,
      { headers: httpHeaders, observe: "response" })
  }

  // Add next action to talent.
  addTalentNextAction(nextAction: INextAction): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/talent/${nextAction.talentID}/talent-next-action/credential/${credentialID}`,
      nextAction, { headers: httpHeaders })
  }

  // Update next action of talent.
  updateTalentNextAction(nextAction: INextAction): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put(`${this.generalURL}/talent/${nextAction.talentID}/talent-next-action/${nextAction.id}/credential/${credentialID}`,
      nextAction, { headers: httpHeaders })
  }

  // Delete next action of talent.
  deleteTalentNextAction(nextActionID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete(`${this.generalURL}/talent/${talentID}/talent-next-action/${nextActionID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  // *******************************API CALLS FOR CAREER PLAN********************************
  // Get all career plans by talent.
  getCareerPlansByTalent(talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/career-plan/talent/${talentID}`,
      { headers: httpHeaders });
  }

  // Add career plan.
  addCareerPlan(careerPlans: any[], talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/career-plan/talent/${talentID}/credential/${credentialID}`,
      careerPlans, { headers: httpHeaders });
  }

  // Update career plan.
  updateCareerPlan(careerPlan: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/career-plan/${careerPlan.id}/talent/${careerPlan.talentID}/credential/${credentialID}`,
      careerPlan, { headers: httpHeaders });
  }

  // Delete career plan.
  deleteCareerPlan(careerObjectiveID: any, talentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/career-plan/${careerObjectiveID}/talent/${talentID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // *******************************API CALLS FOR WAITING LIST********************************
  // Get waiting list by talent.
  getWaitingListByTalent(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list/talent/${talentID}`,
      { params: params, headers: httpHeaders });
  }

  // Get two waiting lists.
  getTwoWaitingLists(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-two`,
      { params: params, headers: httpHeaders });
  }

  // Add waiting list.
  addWaitingList(waitingList: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/waiting-list/credential/${credentialID}`,
      waitingList, { headers: httpHeaders });
  }

  // Update waiting list.
  updateWaitingList(waitingList: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/waiting-list/${waitingList.id}/credential/${credentialID}`,
      waitingList, { headers: httpHeaders });
  }

  // Transfer waiting list.
  transferWaitingList(updateWaitingList: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/waiting-list-transfer/credential/${credentialID}`,
      updateWaitingList, { headers: httpHeaders });
  }

  // Delete waiting list.
  deleteWaitingList(waitingListID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/waiting-list/${waitingListID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // *******************************API CALLS FOR REDIRECTED PAGES********************************
  // Get all talents by waiting list.
  getTalentsByWaitingList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent/waiting-list/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  // Get all talents by requirement.
  getAllTalentsByRequirement(requirementID: string, limit: any, offset: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/requirement/${requirementID}/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, observe: 'response' });
  }

  // Get all talents by campus drive.
  getTalentsByCampusDrive(campusDriveID: string, limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/campus-drive/${campusDriveID}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents by seminar.
  getTalentsBySeminar(seminarID: string, limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/seminar/${seminarID}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents for professional summary report.
  getTalentsForProSummaryReport(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/pro-summary-report/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents for professional summary report by technolohy count.
  getTalentsForProSummaryReportTechnologyCount(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/pro-summary-report-tech-count/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents for fresher summary report.
  getTalentsForFresherSummaryReport(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/fresher-summary-report/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents for package summary report.
  getTalentsForPackageSummaryReport(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/package-summary-report/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' })
  }

  // Get talents for excel download.
  getExcelDownloadTalents(limit: any, offset: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get<any>(`${this.generalURL}/talent/excel-download/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, observe: 'response' })
  }

  // Get searched talents for excel download.
  getSearchedExcelDownloadTalents(conditions: any, limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/talent-search/excel-download/limit/${limit}/offset/${offset}`, conditions,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  //******************************* TALENT EVENT REGISTRATION API CALLS ********************************************* 

  addTalentRegistration(registration: ITalentEventRegistration) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')

    return this._http.post(`${this.generalURL}/talent-event-registration/credential/${credentialID}`,
      registration, { headers: httpHeaders });
  }

  updateTalentRegistration(registration: ITalentEventRegistration) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')

    return this._http.put(`${this.generalURL}/talent-event-registration/${registration.id}/credential/${credentialID}`,
      registration, { headers: httpHeaders });
  }

  // Get all talents for package summary report.
  getTalentRegistrations(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get<any>(`${this.generalURL}/talent-event-registration/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get all talents for package summary report.
  getTalentRegistration(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get<any>(`${this.generalURL}/talent-event-registration`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  //************************************** TALENT ASSIGNMENT SUBMISSION ************************************************* */

  getTalentSubmissions(sessionAssignmentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get<any>(`${this.generalURL}/session-assignment/${sessionAssignmentID}/talent-submission`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  getTalentAssignmentScores(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get<any>(`${this.generalURL}/batch/${batchID}/talent-score`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  addTalentSubmission(sessionAssignmentID: string, talentSubmission: ITalentAssignmentSubmission): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.post<any>(`${this.generalURL}/session-assignment/${sessionAssignmentID}/talent-submission`,
      talentSubmission, { headers: httpHeaders, observe: 'response' });
  }

}

export interface ITalent {
  id?: string,

  // Personal Details.
  firstName: string,
  lastName: string,
  email: string,
  contact: string,
  academicYear: number,
  address?: string,
  city?: string,
  pinCode?: number,
  state?: IState,
  country?: ICountry,
  code?: string,
  talentType?: number,
  personalityType?: string,
  resume?: string,
  alternateContact?: string,
  alternateEmail?: string,
  loyaltyPoints?: string,
  lifetimeValue?: number,
  experienceInMonths?: number
  image: string,

  // Maps.
  technologies?: ITechnology[],

  // Related table IDs.
  salesPersonID?: string,
  sourceID?: string,
  referralID?: string,

  // Flags.
  isActive?: boolean,
  isMastersAbroad?: boolean,
  isSwabhavTalent: boolean,
  isExperience?: boolean,

  // Child tables.
  academics?: IAcademic[],
  experiences?: IExperience[],
  mastersAbroad?: IMastersAbroad,

  // Social media.
  facebookUrl?: string,
  instagramUrl?: string,
  githubUrl?: string,
  linkedInUrl?: string,
}

export interface ITalentDTO {
  id?: string,

  // Personal Details.
  address: string
  city: string
  pinCode: number
  state: IState
  country: ICountry
  code: string
  firstName: string
  lastName: string
  email: string
  contact: string
  academicYear: number
  talentType: number
  personalityType: number
  resume: string
  alternateContact: string
  alternateEmail: string
  loyaltyPoints: string
  lifetimeValue: number
  expectedCTC: number
  referralID: string,
  image: string,

  // Single model.
  salesPerson: ISalesperson
  talentSource: ISource
  mastersAbroad: IMastersAbroadDTO

  // Multiple Model.
  academics: IAcademicDTO[],
  experiences: IExperienceDTO[],
  faculties: IFaculty
  courses: ICourse

  // Maps.
  technologies: ITechnology[],

  // Flags.
  isActive: boolean,
  isMastersAbroad: boolean,
  isSwabhavTalent: boolean,
  isExperience: boolean

  // Social media.
  facebookUrl: string,
  instagramUrl: string,
  githubUrl: string,
  linkedInUrl: string,

  // Extra fields.
  currentCompanyName: string,
  lastDesignation: string,
  currentCompanyPackage: string,
  allExperiencesTechnologies: string[],
  firstCollegeName: string,
  totalYearsOfExperience: number,
  totalYearsOfExperienceInString: string,
  expectedCTCInString: string
}

//********************************************* TALENT EXCEL ************************************************ */

export interface ITalentExcel {
  // Personal Details.
  firstName: string
  lastName: string
  email: string
  contact: string
  academicYear: number
  isSwabhavTalent: string

  // Optional.
  countryName: string
  stateName: string
  city: string
  pinCode: number
  address: string

  // Compulsory if any one is specified
  degreeName: string
  specializationName: string
  collegeName: string
  percentage: number
  yearOfPassout: number
}

export interface ITalentDownloadExcel {

  // Personal Details.
  code: string
  firstName: string
  lastName: string
  email: string
  mobile: string
  alternateEmail: string
  alternateMobile: string
  qualification: string
  specialization: string
  CGPA: number
  collegeName: string
  academicYear: any
  yearOfPassout: number
  country: string
  state: string
  city: string
  pincode: number
  currentCompany: string
  designation: string
  experienceTechnologies: string
  package: number
  fromYear: string
  toYear: string
  totalYearOfExp: number
  talentType: any
  personalityType: string
  salesPerson: string
  source: string
  faculties: string
  courses: string
  technologies: string
  isSwabhavTalent: any
  isActive: any
  facebookUrl: string
  instagramUrl: string
  githubUrl: string
  linkedInUrl: string
  loyaltyPoints: number
  lifetimeValue: number
  expectedCTC: number
  resume: string
}

//**************************************EXPERIENCES******************************************************* */
export interface IExperience {
  id?: string,

  // Maps.
  technologies: ITechnology[],

  // Related table IDs.
  designationID: string,
  talentID: string,

  // Other fields.
  company: string,
  fromDate: string,
  toDate: string,
  package: number,
}

export interface IExperienceDTO {
  id?: string,

  // Maps.
  technologies: ITechnology[],

  // Related tables.
  designation: IDesignation,

  // Other fields.
  company: string,
  fromDate: string,
  toDate: string,
  package: number,
  talentID: string,

  // Extra fields.
  yearsOfExperience: string,
  yearsOfExperienceInNumber: number
}

//**************************************ACADEMICS******************************************************* */
export interface IAcademic {
  id?: string,

  // Related table IDs.
  degreeID: string,
  collegeID: string,
  specializationID: string,
  talentID?: string

  // Other fields.
  college?: string
  percentage: number,
  passout: number,
}

export interface IAcademicDTO {
  id?: string,

  // Realted tables.
  degree: IDegree,
  specialization: ISpecialization,
  college: string,

  // Other fields.
  percentage: number,
  passout: number,
  talentID: string
}

//**************************************MASTERS ABROAD******************************************************* */
export interface IMastersAbroad {
  id?: string

  // Related table IDs.
  talentID?: string
  enquiryID?: string
  degreeID: string

  // Maps.
  countries: ICountry[]
  universities: IUniversity[]

  // Child tables.
  scores: IScore[]

  // Other fields.
  yearOfMS: number
}

export interface IMastersAbroadDTO {
  id?: string

  // Maps.
  countries: ICountry[]
  universities: IUniversity[]

  // Related tables and IDs.
  degree: IDegree
  talentID?: string
  enquiryID?: string

  // Child tables.
  scores: IScoreDTO[]

  // Other fields.
  yearOfMS: number
}

export interface IScore {
  id?: string

  // Related table IDs.
  examinationID: string
  mastersAbroadID: string

  // Other fields.
  marksObtained: number
}

export interface IScoreDTO {
  id?: string

  // Related tables and IDs.
  examination: IExamination
  mastersAbroadID: string

  // Other fields.
  marksObtained: number
}

export interface IExamination {
  id?: string
  name: string
  totalMarks: number
}

//**************************************CALL RECORDS******************************************************* */
export interface ICallRecord {
  id?: string

  // Related table IDs.
  purposeID: string
  outcomeID: string
  talentID: string

  // Other fields.
  dateTime: string
  comment: string
  expectedCTC: number
  noticePeriod: number
  targetDate: string
}

export interface ICallRecordDTO {
  id?: string

  // Related tables.
  purpose: IPurpose
  outcome: IOutcome

  // Other fields.
  dateTime: string
  comment: string
  expectedCTC: number
  noticePeriod: number
  targetDate: string
  talentID: string
}

//**************************************LIFETIME VALUE******************************************************* */
export interface ILifetimeValue {
  id?: string
  upsell: number
  placement: number
  knowledge: number
  teaching: number
}

//**************************************NEXT ACTION******************************************************* */
export interface INextAction {
  id?: string

  // Related table IDs.
  talentID: string
  actionTypeID: string

  // Maps.
  courses: any[]
  companies: any[]
  technologies: any[]

  // Other fields.
  stipend: number
  referralCount: number
  fromDate: string
  toDate: string
  targetDate: string
  comment: string
}

export interface INextActionDTO {
  id?: string

  // Related tables.
  actionType: INextActionType

  // Maps.
  courses: any[]
  companies: any[]
  technologies: any[]

  // Other fields.
  talentID: string
  stipend: number
  referralCount: number
  fromDate: string
  toDate: string
  targetDate: string
  comment: string
}

export interface INextActionType {
  id?: string
  type: string
}

//**************************************WAITING LIST******************************************************* */
export interface IWaitingList {
  id?: string

  // Talent/Enquiry related fields.
  talentID?: string
  enquiryID?: string
  email: string
  isActive: boolean

  // Company related IDs.
  companyBranchID?: string
  companyRequirementID?: string

  // Course related IDs.
  courseID?: string
  batchID?: string

  // Other fields.
  sourceID: string
}

export interface IWaitingListDTO {
  id?: string

  // Talent/Enquiry related fields.
  talentID?: string
  enquiryID?: string
  email: string
  isActive: boolean

  // Company related models.
  companyBranch: any
  companyRequirement: any

  // Course related models.
  course: any
  batch: any

  // Other fields.
  source: ISource
}

//************************************** TALENT EVENT REGISTRATION ************************************************* */
export interface ITalentEventRegistration {
  id?: string
  talentID?: string
  eventID?: string
  talent?: ITalent
  event?: IEvent
  registrationDate: string
  hasAttended: boolean
  isTalentRegistered: boolean
}

//************************************** TALENT ASSIGNMENT SUBMISSION ************************************************* */
export interface ITalentAssignmentSubmission {
  id?: string
  talentID: string
  talent?: any
  batchTopicAssignmentID: string
  batchTopicAssignment: IBatchTopicAssignment
  isAccepted?: boolean
  isChecked?: boolean
  acceptanceDate: string
  facultyRemarks?: string
  score?: number
  solution?: string
  submittedOn?: string | Date
}

export interface ITalentAssignmentScore {
  talent: any
  talentSubmission: ITalentAssignmentSubmission[]
}

//************************************** OTHERS ********************************************** */
export interface IUniversity {
  id?: string
  universityName: string
  countryID?: string
  countryName?: string
  isVisible?: boolean
  country: ICountry
}

export interface ISalesperson {
  id?: string,
  firstName: string,
  lastName: string
}

export interface IState {
  id?: string,
  name: string,
  countryID?: string
}

export interface ICountry {
  id?: string,
  name: string
}

export interface IDegree {
  id?: string
  name: string
}

export interface ISpecialization {
  id?: string
  branchName: string
  degreeID: string
}

export interface ISource {
  id?: string
  name: string
  urlName: string
  description: string
}

export interface IFaculty {
  id: string,
  firstName: string,
  lastName: string
}

export interface IDesignation {
  id?: string,
  position: string
}

export interface ITechnology {
  id?: string,
  language: string
}

export interface IPurpose {
  id?: string,
  purpose: string
  purposeType: string
}

export interface IOutcome {
  id?: string,
  purposeID: string
  outcome: string
}

export interface ISearchSection {
  name: string
  isSelected: boolean
}

export interface ISearchFilterField {
  propertyName: string
  propertyNameText: string
  valueList: any[]
}


