import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProSummaryReportService {

  generalURL: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) { 
    this.generalURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // Get professional summary report of talents.
  getProfessionalSummaryReportList(limit: any, offset: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/professional-summary-report`, 
    {params:params, headers: httpHeaders, observe: "response" });
  }

  // Get talent counts by technology and company name.
  getProfessionalSummaryReportByTechCount(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/professional-summary-report-tech-count`, 
    {params:params, headers: httpHeaders});
  }
}

// Stores company and its respective talent counts by category.
// firstCount: 12- 24 months of experience.
// secondCount: 24- 60 months of experience.
// thirdCount: 60- 84 months of experience.
// fourthCount: above 84 months of experience.
export interface IProfessionalSummaryReport {
  company: string   
  firstCount: number 
  secondCount: number 
  thirdCount: number 
  fourthCount: number 
  firstCatergoryTechCount: ICompanyTechnologyTalent[]
  secondCatergoryTechCount: ICompanyTechnologyTalent[]
  thirdCatergoryTechCount: ICompanyTechnologyTalent[]
  fourthCatergoryTechCount: ICompanyTechnologyTalent[]
}

// Stores total talent counts for each category.
// firstCount: 12- 24 months of experience.
// secondCount: 24- 60 months of experience.
// thirdCount: 60- 84 months of experience.
// fourthCount: above 84 months of experience.
export interface IProfessionalSummaryReportCounts {
  firstCountTotal: number  
  secondCountTotal: number 
  thirdCountTotal: number  
  fourthCountTotal: number 
}

// Stores technology and its respectove talent count for each company.
export interface ICompanyTechnologyTalent {
  company: string    
  techID: string      
  techName: string    
  talentCount: number
}


