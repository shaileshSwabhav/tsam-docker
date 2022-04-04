import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ICompanyBranch } from '../company/company.service';
import { Constant } from '../constant';
import { ICourse } from '../course/course.service';
import { ITechnologies } from '../general/general.service';
import { LocalService } from '../storage/local.service';
import { IFaculty } from '../talent/talent.service';
import { ITechnology } from '../technology/technology.service';

@Injectable({
  providedIn: 'root'
})
export class ReportService {

  private tenantURL: string
  private talentNextActionReportUrl: string
  
  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.talentNextActionReportUrl = `${this.tenantURL}/talent-next-action-report`
  }

// ================================================= TALENT NEXT ACTION REPORT =================================================


  getCallingReports(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.talentNextActionReportUrl}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

// ======================================================== FRESHER SUMMARY ========================================================

  getAllFresherSummary(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/fresher-summary-report`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getTechnologyFresherSummary(params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/fresher-summary-report/technology`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getSummaryTechnologyList(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/fresher-summary-report/technology-list`,
      { headers: httpHeaders, observe: "response" })
  }

// ======================================================== FRESHER SUMMARY ========================================================
    
  getAllPackageSummary(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/package-summary-report`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getTechnologyPackageSummary(params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/package-summary-report/technology`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

//  ============================================ FACULTY REPORT ============================================

  getAllFacultyReport(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/faculty-report`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  
  // ============================================== LOGIN REPORT ==============================================

  // Get login report.
  getLoginReports(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.tenantURL}/login-report/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get credential login report.
  getCredentialLoginReports(credentialID: string, limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.tenantURL}/credential/${credentialID}/login-report/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }


}

export interface INextActionReport {
  id: string
  loginID: string
  loginName: string
  firstName: string
  lastName: string
  talentID: string
  stipend: number
  referralCount: number
  fromDate: string
  toDate: string
  targetDate: string
  comment: string
  actionType: string
  courses: ICourse[]
  companies: ICompanyBranch[]
  technologies: ITechnologies[]
}

//  ============================================ FRESHER SUMMARY ============================================

export interface IFresherSummary {
  columnName: string
  academic: any
  outstandingCount: number
  excellentCount: number
  averageCount: number
  unrankedCount: number
}

export interface ITechnologyFresherSummary {
  technology: ITechnology
  totalCount: number
  academic: any
  columnName: string
  technologyLangugage: string
}

export interface IAcademicTechnologySummary {
  columnName: string
  academic: string
  technologySummary: ITechnologyFresherSummary[]
}

//  ============================================ PACKAGE SUMMARY ============================================

export interface IPackageSummary {
  experience: string
  lessThanThree: number
  threeToFive: number
  fiveToTen: number
  tenToFifteen: number
  greaterThanFifteen: number
}

export interface ITechnologyPackageSummary {
  technology: ITechnology
  techLanguage: string
  totalCount: number
}

export interface IExperienceTechnologySummary {
  experience: string
  technologySummary: ITechnologyPackageSummary[]
}

//  ============================================ FACULTY REPORT ============================================

export interface IFacultyReport {
  faculty: IFaculty
  week: string
  monday: IBatch[]
  tuesday: IBatch[]
  wednesday: IBatch[]
  thursday: IBatch[]
  friday: IBatch[]
  saturday: IBatch[]
  sunday: IBatch[]
  totalTrainingHours: number
  workingHours: Map<string, IWorkingHours>
}

export interface IBatch {
  id: string
  batchName: string
  batchStatus: string
  batchTimings: IBatchTiming[]
  totalDailyHours: number
}

export interface IBatchTiming {
  id: string
  batchID: string
  day: any
  fromTime: string
  toTime: string
}

export interface IWorkingHours {
  batchName: string
  totalHours: number
}

  // ============================================== LOGIN REPORT ==============================================

export interface ICredentialLoginReports {
  loginName: string
  loginTime: string
  logoutTime: string
  totalHours: string
  roleName: string
}

export interface ILoginReports {
  loginSessionID: string
  credentialID: string
  loginName: string
  roleName: string
  lastLoginTime: string
  lastLogoutTime: string
  loginTime: string
  logoutTime: string
  loginCount: number
  totalHours: string
}