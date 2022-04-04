import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class WaitingListReportService {

  generalURL: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) { 
    this.generalURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // Get waiting list report by company branch.
  getWaitingListCompanyBranchReportList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-report-company-branch/limit/${limit}/offset/${offset}`, 
    {params:params, headers: httpHeaders, observe: "response" });
  }

  // Get waiting list report by course.
  getWaitingListCourseReportList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-report-course/limit/${limit}/offset/${offset}`, 
    {params:params, headers: httpHeaders, observe: "response" });
  }

   // Get waiting list report by company requirement.
   getWaitingListRequirementReportList(companyBranchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-report-requirement/company-branch/${companyBranchID}`, 
    {params:params, headers: httpHeaders});
  }

  // Get waiting list report by batch.
  getWaitingListBatchReportList(courseID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-report-batch/course/${courseID}`, 
    {params:params, headers: httpHeaders});
  }

  // Get waiting list report by technology.
  getWaitingListTechnologyReportList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/waiting-list-report-technology/limit/${limit}/offset/${offset}`, 
    {params:params, headers: httpHeaders, observe: "response" });
  }
}

export interface IWaitingListCompanyBranchDTO {
  companyBranch: any 
  talentCount: number    
  enquiryCount: number   
}

export interface IWaitingListCourseDTO {
  course: any 
  talentCount: number    
  enquiryCount: number   
}

export interface IWaitingListRequirementDTO {
  requirement: any 
  talentCount: number    
  enquiryCount: number   
}

export interface IWaitingListBatchDTO {
  batch: any 
  talentCount: number    
  enquiryCount: number   
}

export interface IWaitingListTechnologyDTO {
  technology: any 
  talentCount: number    
  enquiryCount: number   
}

// export interface ITalentWaitingList {
// 	// Talent related details.
// 	talentID: string       
// 	firstName: string             
// 	lastName: string              
// 	contact: string               
// 	academicYear: number   
// 	college: string               
// 	salesPersonFirstName: string 
// 	salesPersonLastName: string        

// 	// Waiting list related details.
// 	companyBranch: string       
// 	companyRequirement: string 
//   companyRequirementCode: string 
// 	course: string              
// 	batch: string               
// }

// export interface IEnquiryWaitingList {
// 	// Enquiry related details.
// 	enquiryID: string       
// 	firstName: string             
// 	lastName: string              
// 	contact: string               
// 	academicYear: number   
// 	college: string               
// 	salesPersonID: string         

// 	// Waiting list related details.
// 	companyBranchID: string       
// 	companyRequirementID: string  
// 	courseID: string              
// 	batchID: string               
// }
