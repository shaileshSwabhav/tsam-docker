import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { UserloginService } from '../login/userlogin.service';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {

  private tenantURL: string
  private httpHeaders: HttpHeaders

  constructor(
    private http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) { 
    this.tenantURL = `${this._constant.BASE_URL}/tenant/${_constant.TENANT_ID}`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

  }

  getSalesPeopleDashboardDetails(params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    // console.log(params)
    return this.http.get(`${this.tenantURL}/dashboard/admin/sales-people`,
      { params: params, headers: httpHeaders, observe: 'response' })
  }

  getTalentDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/talent`,
      { headers: httpHeaders, observe: 'response' })
  }

  getTalentEnquiryDashboardDetails(param: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/talent-enquiry`,
      { params: param, headers: httpHeaders, observe: 'response' })
  }

  getTalentEnquirySource(param: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/talent-enquiry-source`,
      { params: param, headers: httpHeaders, observe: 'response' })
  }

  getFacultyDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/faculty`,
      { headers: httpHeaders, observe: 'response' })
  }

  getCompanyDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/company`,
      { headers: httpHeaders, observe: 'response' })
  }

  getCollegeDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/college`,
      { headers: httpHeaders, observe: 'response' })
  }

  getBatchDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/batch`,
      { headers: httpHeaders, observe: 'response' })
  }

  getTechnologyDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/technology`,
      { headers: httpHeaders, observe: 'response' })
  }

  getCourseDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin/course`,
      { headers: httpHeaders, observe: 'response' })
  }

  getAllDashboardDetails(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/admin`,
      { headers: httpHeaders, observe: 'response' })
  }

  getBatchPerformance(params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/batch/score`, 
    { params: params, headers: httpHeaders, observe: 'response' })
  }

  getSessionWiseTalentScore(talentID: string, batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/batch/${batchID}/talent/${talentID}/session/score`, 
      { headers: httpHeaders, observe: 'response' })
  }

  getBatchDetails(params: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/batch/status`, 
      { params: params, headers: httpHeaders, observe: 'response' })
  }

  getBatchTalents(batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/batch/${batchID}/status`, 
      { headers: httpHeaders, observe: 'response' })
  }

  getTalentFeedbackScore(batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/dashboard/batch/${batchID}/talent-feedbacks`, 
      { headers: httpHeaders, observe: 'response' })
  }
}


export interface IBatchPerformance {
  totalBatches: number
  outstanding: number
  good: number
  average: number
  keywordNames: IGroupWiseKeywordName[]
  outstandingTalent: ITalentFeedbackScore[]
  goodTalent: ITalentFeedbackScore[]
  averageTalent: ITalentFeedbackScore[]
}

export interface ITalentFeedbackScore {
  talentID: string
  firstName: string
  lastName: string
  personalityType: string
  talentType: number
  batchID: string
  batchName: string
  score: number
  interviewRating: number
  feedbackKeywords: ISessionGroupScore[]
  isVisibile: boolean
}

export interface IFeedbackKeyword {
  keyword: string
  keywordScore: number
}

interface ITalentFeelingDetails {
  feelingName: string
  levelNumber: number
  description: string
}

export interface IBatchDetails {
  batchID: string
  batchName: string
  batchType: boolean
  totalStudents: number
  isVisible: boolean
}

export interface IKeyword {
  id: string
  name: string
}

export interface ITalentSessionFeedbackScore {
  talentID: string
  firstName: string
  lastName: string
  sessionDate: string
  batchID: string
  batchName: string
  keywordNames: IGroupWiseKeywordName[]
  sessionFeedback: ISessionKeywordFeedback[]
  feelingDetails: ITalentFeelingDetails[]
}

interface ISessionKeywordFeedback {
  sessionName: string
  order: number
  score: number
  batchSessionID: string
  feelingName: string
  levelNumber: string
  description: string
  feedbackGroup: ISessionGroupScore[]
  // feedbackKeywords: IFeedbackKeyword[]
  // batchSessionID: string
  // courseSessionID: string
}
// SessionGroupScore
export interface ISessionGroupScore {
  groupID: string
  groupName: string
  feedbackScore: IFeedbackKeyword[]
  showGroupDetails: boolean
}

export interface IGroupWiseKeywordName {
  groupName: string
  keywords: IKeyword[]
  showGroupDetails: boolean
}

export interface IFeedback {
  talentFeedback: ISessionTalentFeedbackScore[]
  keywords: IKeyword[]
}

export interface ISessionTalentFeedbackScore {
  talentID: string
  firstName: string
  lastName: string
  personalityType: string
  talentType: number
  batchName: string
  batchID: string
  score: number
  interviewRating: number
  sessionFeedback: IFeedbackKeyword[]
}