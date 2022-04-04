import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class EnquiryService {
  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService
  ) {
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`;
  }


  //*******************************ENQUIRY API CALLS********************************************* 
  // Add new enquiry.
  addEnquiry(enquiry: any): Observable<any> {
    // let enquiryJson:string = JSON.stringify(enquiry);
    // console.log("enquiry" + enquiryJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.generalURL}/talent-enquiry/credential/${credentialID}`, 
    enquiry, { headers: httpHeaders });
  }

  // Add new enquiry from excel.
  addEnquiryFromExcel(enquiry: any): Observable<any> {
    // let talentJson:string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.generalURL}/talent-enquiry-excel/credential/${credentialID}`, 
    enquiry, { headers: httpHeaders });
  }

  // Add multiple enquiries.
  addEnquiries(enquiries: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')

    return this._http.post<any>(`${this.generalURL}/enquiries/credential/${credentialID}`, 
      enquiries, { headers: httpHeaders });
  }

  // Get all searched enquiries.
  getAllSearchedEnquiries(conditions: any, limit: number, offset: number, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/talent-enquiry/search/limit/${limit}/offset/${offset}`, conditions, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Get all enquiries.
  getEnquiries(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-enquiry/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Get all enquiries by waiting list.
  getEnquiriesByWaitingList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-enquiry/waiting-list/limit/${limit}/offset/${offset}`, 
    {params:params, headers: httpHeaders, observe: "response" });
  }

  // Update enquiry.
  updateEnquiry(enquiry: any): Observable<any> {
    // let enquiryJson:string = JSON.stringify(enquiry);
    // console.log("enquiry" + enquiryJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-enquiry/${enquiry.id}/credential/${credentialID}`, 
    enquiry, { headers: httpHeaders });
  }

  // Update enquiries' salesperson.
  updateEnquiriesSalesPerson(enquiries: any, salesPersonID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-enquiry/saleperson/${salesPersonID}/credential/${credentialID}`, 
    enquiries, { headers: httpHeaders });
  }

  // Update enquiry.
  convertEnquiryToTalent(enquiryID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-enquiry/${enquiryID}/convert-to-talent/credential/${credentialID}`, {},
      { headers: httpHeaders });
  }

  // Delete enquiry.
  deleteEnquiry(id: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/talent-enquiry/${id}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  // *******************************ENQUIRY CALL RECORDS API CALLS********************************
  // Get all call records by enquiry.
  getCallRecordsByEnquiry(enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-enquiry-call-record/talent-enquiry/${enquiryID}`, 
    { headers: httpHeaders });
  }

  // Add call record.
  addCallRecord(callRecord: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/talent-enquiry-call-record/talent-enquiry/${enquiryID}/credential/${credentialID}`, callRecord, 
    { headers: httpHeaders, observe: "response" });
  }

  // Get one call record by id.
  getCallRecord(callRecordID: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/talent-enquiry-call-record/${callRecordID}/talent-enquiry/${enquiryID}`, 
    { headers: httpHeaders, observe: "response" });
  }

  // Update enquiry call record.
  updateCallRecord(callRecord: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/talent-enquiry-call-record/${callRecord.id}/talent-enquiry/${enquiryID}/credential/${credentialID}`, 
    callRecord, { headers: httpHeaders });
  }

  // Delete enquiry callrecord.
  deleteCallRecord(callRecordID: any, enquiryID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/talent-enquiry-call-record/${callRecordID}/talent-enquiry/${enquiryID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  // *******************************API CALLS FOR WAITING LIST********************************
  // Get waiting list by enquiry.
  getWaitingListByEnquiry(enquieyID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list/enquiry/${enquieyID}`, 
    { headers: httpHeaders });
  }

  // Add waiting list.
  addWaitingList(waitingList: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/waiting-list/credential/${credentialID}`, 
    waitingList, { headers: httpHeaders});
  }

  // Update waiting list.
  updateWaitingList(waitingList: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/waiting-list/${waitingList.id}/credential/${credentialID}`, 
    waitingList, { headers: httpHeaders });
  }

  // Delete waiting list.
  deleteWaitingList(waitingListID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/waiting-list/${waitingListID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }
}

//**************************************ENQUIRY******************************************************* */
export interface IEnquiry {
  id?: string,

	// Personal Details.
  address: string,
  city: string,
  pinCode: number,
  state: IState,
  country: ICountry,
  code: string,
  firstName: string,
  lastName: string,
  email: string,
  contact: string,
  alternateContact: string,
  alternateEmail: string,
  enquiryDate: string,
  enquiryType: string,
  talentID: string
  additionalDetails: string,
  academicYear: number,
  resume: string,

	// Maps.
  technologies: ITechnology[],
  courses: ICourse[],

	// Related table IDs.
  salesPersonID: string,
  sourceID: string,

	// Flags.
  isMastersAbroad: boolean,
  isExperience: boolean,

	// Child tables.
  academics: IAcademic[],
  experiences: IExperience[],
  mastersAbroad: IMastersAbroad,

	// Social media.
  facebookUrl: string,
  instagramUrl: string,
  githubUrl: string,
  linkedInUrl: string,
}

export interface IEnquiryDTO {
  id?: string,
  
	// Personal Details.
  address: string,
  city: string,
  pinCode: number,
  state: IState,
  country: ICountry,
  code: string,
  firstName: string,
  lastName: string,
  email: string,
  contact: string,
  alternateContact: string,
  alternateEmail: string,
  enquiryDate: string,
  enquiryType: string,
  talentID: string
  additionalDetails: string,
  academicYear: number,
  resume: string,
  expectedCTC: number

	// Single model.
  salesPerson: ISalesperson
  talentSource: ISource
  mastersAbroad: IMastersAbroadDTO

	// Multiple Model.
  academics: IAcademicDTO[],
  experiences: IExperienceDTO[],

	// Maps.
  technologies: ITechnology[],
  courses: ICourse[],

	// Flags.
  isMastersAbroad: boolean,
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

//**************************************EXPERIENCES******************************************************* */
export interface IExperience {
  id?: string,

	// Maps.
  technologies: ITechnology[],

	// Related table IDs.
  designationID: string,
  enquiryID: string,

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
  enquiryID: string,

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
  enquiryID: string

	// Other fields.
  college: string
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
  enquiryID: string
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
  enquiryID: string

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
  enquiryID: string
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

//**************************************OTHERS******************************************************* */
export interface ICollegeBranch {
  id?: string
  branchName: string
  code: string
}

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
  branchName: string,
  degreeID: string
}

export interface ISource {
  id?: string,
  name: string
  urlName: string
  description: string
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

export interface ICourse {
  id?: string,
  code: string
  name: string
}

export interface ISearchSection {
  name: string
  isSelected: boolean
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}

export interface IEnquiryExcel {

	// Personal Details.
  firstName: string
  lastName: string
  email: string
  contact: string
  academicYear: number

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