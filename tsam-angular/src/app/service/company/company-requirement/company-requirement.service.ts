import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../../constant';
import { IDesignation, ISalaryTrend } from '../../general/general.service';
import { LocalService } from '../../storage/local.service';
import { ITechnology } from '../../technology/technology.service';

@Injectable({
  providedIn: 'root'
})
export class CompanyRequirementService {

  companyRequirementURL: string

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService
  ) {
    this.companyRequirementURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}/company-requirement`
   }

  // Add new company requirement.
  addCompanyRequirement(requirement: any): Observable<any> {
    // let talentJson:string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.companyRequirementURL}`, 
    requirement, { headers: httpHeaders });
  }

  // Get all requirements.
  getCompanyRequirements(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyRequirementURL}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Update requirement.
  updateCompanyRequirement(requirement: any): Observable<any> {
    // let talentJson: string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.companyRequirementURL}/${requirement.id}`, requirement, 
    { headers: httpHeaders });
  }

  // Allocate salesperson to requirement.
  allocateSalesPersonToCompanyRequirement(requirementIDsToBeUpdated: any[], salesPersonID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.companyRequirementURL}/salesperson/${salesPersonID}`, requirementIDsToBeUpdated, 
    { headers: httpHeaders });
  }

  // Close requirement.
  closeCompanyRequirement(requirementID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.companyRequirementURL}/close/${requirementID}`, 
    { headers: httpHeaders });
  }

  // Delete requirement.
  deleteCompanyRequirement(id: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.companyRequirementURL}/${id}`, 
    { headers: httpHeaders });
  }

  // Adds talents to requirements.
  addTalentsToRequirement(id: any, requirementTalents: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.companyRequirementURL}/${id}/talent`,
    requirementTalents, { headers: httpHeaders });
  }

  // Get one requirement.
  getCompanyRequirement(companyRequirementID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyRequirementURL}/${companyRequirementID}`, 
    {headers: httpHeaders});
  }

  // Get one requirement's company deatils.
  getCompanyDetails(companyRequirementID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyRequirementURL}/${companyRequirementID}/company-details`, 
    {headers: httpHeaders});
  }

  // Get my opputunities.
  getMyOppurtunities(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyRequirementURL}/my-opportunities`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }
}

export interface IRequirement {
  id?: string

  // Address.
  jobLocation: IJobLocation

  // Multiple.
  qualifications: IDegree[]
  universities: IUniversity[]
  technologies: ITechnology[]
  selectedTalents: ITalent[]
  designation: IDesignation

  // Single objects.
  salesPersonID: string
  companyID: string
  designationID: string

  // Other fields.
  isActive: boolean
  code: string
  talentRating: string
  personalityType: string
  minimumExperience: number
  maximumExperience: number
  jobRole: string
  jobDescription: string
  jobRequirement: string
  jobType: string
  // packageOffered: number
  requiredBefore: string
  requiredFrom: string
  vacancy: number
  comment: string
  termsAndConditions: string

  // Criteria
  criteria1: string
  criteria2: string
  criteria3: string
  criteria4: string
  criteria5: string
  criteria6: string
  criteria7: string
  criteria8: string
  criteria9: string
  criteria10: string

  // Rating
  rating: number
  
  // Salary trend
  salaryTrend: ISalaryTrend
}

export interface IRequirementDTO {
  id?: string

  // Address.
  jobLocation: IJobLocation

  // Multiple.
  qualifications: IDegree[]
  universities: IUniversity[]
  technologies: ITechnology[]
  selectedTalents: ITalent[]

  // Single objects.
  salesPerson: ISalesperson
  company: ICompanyBranch

  // Other fields.
  isActive: boolean
  code: string
  talentRating: string
  personalityType: string
  minimumExperience: number
  maximumExperience: number
  jobRole: string
  jobDescription: string
  jobRequirement: string
  jobType: string
  packageOffered: number
  packageOfferedInstring: string
  requiredBefore: string
  requiredFrom: string
  vacancy: number
  comment: string
  termsAndConditions: string
  totalApplicants: number
  isUrgent: boolean
  talentRatingValue: string
  criteria1: string
  criteria2: string
  criteria3: string
  criteria4: string
  criteria5: string
  criteria6: string
  criteria7: string
  criteria8: string
  criteria9: string
  criteria10: string
  rating: number

  // Salary trend
  salaryTrend: ISalaryTrend
}

export interface ISearchTalentParams {
  requirementID?: string
  requirementSearch?: string
  qualifications: any[],
  technologies: any[],
  talentType: number,
  isExperience: boolean
  maximumExperience: number,
  minimumExperience: number,
  personalityType: string
  // marksCriteria : number,
  // totalExperience: number,
}

export interface IJobLocation{
  address: string
  city: string
  pinCode: number
  state: IState
  country: ICountry
}

export interface IState {
  id?: string
  name: string
  countryID?: string
}

export interface ICountry {
  id?: string
  name: string
}

export interface IDegree {
  id?: string
  name: string
}

export interface IUniversity {
  id?: string
  universityName: string
  countryID?: string
  countryName?: string
  isVisible?: boolean
  country: ICountry
}

export interface ITalent {
  id?: string
}

export interface ISalesperson {
  id?: string,
  firstName: string,
  lastName: string
}

export interface ICompanyBranch {
  id?: string,
  branchName: string,
  companyID?: string
}

