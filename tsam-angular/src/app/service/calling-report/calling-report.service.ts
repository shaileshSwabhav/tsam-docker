import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class CallingReportService {

  talentCallingReportUrl: string
  talentEnquiryCallingReportUrl: string
  tenantUrl: string
  httpHeaders: HttpHeaders
  // credentialID: string

  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    this.tenantUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.talentCallingReportUrl = `${this.tenantUrl}/talent-calling-report`
    this.talentEnquiryCallingReportUrl = `${this.tenantUrl}/talent-enquiry-calling-report`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
  }

  getCallingReports(limit: number, offset: number, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentCallingReportUrl}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getLoginwiseCallingReports(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentCallingReportUrl}/loginwise`,
      { params: params, headers: httpHeaders })
  }

  getDaywiseCallingReports(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentCallingReportUrl}/daywise`,
      { headers: httpHeaders })
  }

// =======================================================TALENT-ENQUIRY=======================================================

  
  getEnquiryCallingReports(limit: number, offset: number, params?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentEnquiryCallingReportUrl}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getLoginwiseEnquiryCallingReports(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentEnquiryCallingReportUrl}/loginwise`,
      { params: params, headers: httpHeaders })
  }

  getDaywiseEnquiryCallingReports(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.httpClient.get(`${this.talentEnquiryCallingReportUrl}/daywise`,
      { headers: httpHeaders })
  }

  // //get all universities
  // getAllUniversities(limit: number, offset: number, params?: HttpParams): Observable<any> {
  //   return this.httpClient.get(`${this.universityUrl}/limit/${limit}/offset/${offset}`,
  //     { params: params, headers: new HttpHeaders({ 'token': this.localService.getJsonValue('token') }), observe: "response" })
  // }
}



export interface ILoginwiseEnquiryCallingReport {
  credentialID: string
  firstName: string
  lastName: string
  totalCallingCount: number
  totalEnquiryCount: number
}

export interface IDaywiseEnquiryCallingReport {
  date: string
  totalCallingCount: number
  totalEnquiryCount: number
}

export interface ITalentEnquiryCallingReport {
  comment: string,
  dateTime: string,
  loginName: string,
  outcome: string,
  purpose: string,
  talentEnquiry: ITalentEnquiry
}

export interface ITalentEnquiry {
  contact: string,
  email: string,
  firstName: string,
  lastName: string
}