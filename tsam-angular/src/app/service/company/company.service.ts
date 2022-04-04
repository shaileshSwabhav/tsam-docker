import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { scaleService } from 'chart.js';


@Injectable({
  providedIn: 'root'
})


export class CompanyService {

  private companyURL: string
  private domainURL: string
  private companyEnquiryURL: string
  private enquiryCallRecordURL: string
  httpHeaders: HttpHeaders
  tenantID: string
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    // this.companyURL = `${constant.BASE_URL}/tenant/${constant.TENANT_ID}/company`
    this.companyURL = `${constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/company`
    this.companyEnquiryURL = `${constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/company-enquiry`
    this.enquiryCallRecordURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/company-enquiry-call-record`
    this.domainURL = `${constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/domain`
    this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = this.localService.getJsonValue('credentialID')
  }

  // ======================================================COMPANY BRANCH CRUD==========================================================

  // CRUD on company branches

  getAllCompanyBranches(param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}-branch`,
      { params: param, headers: httpHeaders, observe: "response" });
  }

  getAllSearchedCompanyBranches(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}/search`,
      { headers: httpHeaders, observe: "response" });

  }

  getAllBranchesOfCompany(companyID: string, param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}/${companyID}/branch`,
      { params: param, headers: httpHeaders, observe: "response" });

  }

  getCompanyBranchByID(companyID: string, companyBranchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}/${companyID}/branch/${companyBranchID}`,
      { headers: httpHeaders, observe: "response" });
  }

  // Get one company branch.
  getCompanyBranch(companyID: string, companyBranchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyURL}/${companyID}/branch/${companyBranchID}`,
      { headers: httpHeaders });
  }

  getAllBranchesForSalesPerson(salesPersonID: string, param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}/branch/salesperson/${salesPersonID}`,
      { params: param, headers: httpHeaders, observe: "response" });
  }

  addCompanyBranch(companyBranch: ICompanyBranch, companyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })


    return this._http.post(`${this.companyURL}/${companyID}/branch`,
      companyBranch, { headers: httpHeaders });
  }

  addCompanyBranches(companyBranches: ICompanyBranch, companyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.post(`${this.companyURL}/${companyID}/branches`,
      companyBranches, { headers: httpHeaders });
  }

  deleteCompanyBranch(companyID: string, companyBranchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.delete(`${this.companyURL}/${companyID}/branch/${companyBranchID}`, { headers: httpHeaders });
  }

  updateCompanyBranch(companyBranch: ICompanyBranch): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.put(`${this.companyURL}/${companyBranch.companyID}/branch/${companyBranch.id}`, { headers: httpHeaders });
  }

  // ======================================================END==========================================================


  // ======================================================COMPANY ENQUIRY==========================================================

  // Get all company enquiries.
  getAllCompanyEnquiries(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.companyEnquiryURL}`,
      { params: params, headers: httpHeaders, observe: "response" });
  }

  // Add one company enquiry.
  addCompanyEnquiry(companyEnquiry: ICompanyEnquiry): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.companyEnquiryURL}`,
      companyEnquiry, { headers: httpHeaders });
  }

  // Delete one company enquiry.
  deleteCompanyEnquiry(companyEnquiryID: string): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.companyEnquiryURL}/${companyEnquiryID}`,
      { headers: httpHeaders });
  }

  // Update one company enquiry.
  updateCompanyEnquiry(companyEnquiry: ICompanyEnquiry): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.companyEnquiryURL}/${companyEnquiry.id}`, companyEnquiry,
      { headers: httpHeaders });
  }

  // Update one or more enquiry's salesperson.
  updateCompanyEnquirysSalesPerson(enquiries: any, salesPersonID): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.companyEnquiryURL}/saleperson/${salesPersonID}`,
      enquiries, { headers: httpHeaders });
  }

  // *******************************COMPANY ENQUIRY CALL RECORDS API CALLS********************************
  // Get all enquiry call records by enquiry id.
  getCallRecordsByEnquiry(enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.enquiryCallRecordURL}/enquiry/${enquiryID}`,
      { headers: httpHeaders });
  }

  // Add one enquiry call record.
  addCallRecord(callRecord: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID: string = this.localService.getJsonValue('credentialID')

    return this._http.post(`${this.enquiryCallRecordURL}/enquiry/${enquiryID}`,
      callRecord, { headers: httpHeaders, observe: "response" });
  }

  // Update one enquiry call record.
  updateCallRecord(callRecord: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.enquiryCallRecordURL}/${callRecord.id}/enquiry/${enquiryID}`,
      callRecord, { headers: httpHeaders });
  }

  // Delete one enquiry call record.
  deleteCallRecord(callRecordID: any, enquiryID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.enquiryCallRecordURL}/${callRecordID}/enquiry/${enquiryID}`,
      { headers: httpHeaders });
  }

  // =====================================================================================================================================

  // CRUD on companies

  addCompany(company: ICompany): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.post<any>(`${this.companyURL}`, company,
      { headers: httpHeaders });

  }

  updateCompany(company: ICompany): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.put(`${this.companyURL}/${company.id}`, company,
      { headers: httpHeaders });

  }


  getCompanyByID(companyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}/${companyID}`,
      { headers: httpHeaders });

  }

  getAllCompanies(param?: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.companyURL}`,
      { params: param, headers: httpHeaders, observe: "response" });
  }


  deleteCompany(companyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this._http.delete(`${this.companyURL}/${companyID}`,
      { headers: httpHeaders });

  }



  // CRUD for domain
  addDomain(domain: IDomain): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.post(this.domainURL, domain,
      { headers: httpHeaders });

  }

  updateDomain(domain: IDomain): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.put(`${this.domainURL}/${domain.id}`, domain,
      { headers: httpHeaders });

  }


  getDomainByID(domainID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.domainURL}/${domainID}`,
      { headers: httpHeaders });

  }

  getAllDomains(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.domainURL}/list`,
      { headers: httpHeaders });

  }


  deleteDomain(domainID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.delete(`${this.domainURL}/${domainID}`,
      { headers: httpHeaders });

  }
}


// model definitions

export interface ICompanyBranch {
  id?: string,
  companyBranchCode?: string,
  address: string,
  city: string,
  pinCode: string,
  country: ICountry,
  state: IState,
  companyID: string,
  companyRating: number,
  companyName?: string,
  companyCode?: string,
  website?: string,
  mainBranch: boolean,
  domains: IDomain[],
  technologies: ITechnology[],
  hrHeadName: string,
  hrHeadContact: string,
  hrHeadEmail: string,
  technologyHeadName: string,
  technologyHeadContact: string,
  technologyHeadEmail: string,
  unitHeadName: string,
  unitHeadContact: string,
  unitHeadEmail: string,
  financeHeadName: string,
  financeHeadContact: string,
  financeHeadEmail: string,
  recruitmentHeadName: string,
  recruitmentHeadContact: string,
  recruitmentHeadEmail: string,
  numberOfEmployees: number,
  salesPersonID: string,
  salesPersonName?: string,
  salesPerson?: ISalesperson
  termsAndConditions?: string,
  onePager?: string
}


export interface ICompany {
  id?: string,
  companyName: string,
  companyCode: string,
  about: string,
  logo: string,
  website: string,
  branches: ICompanyBranch[]
}



export interface ICompanyRequirement {
  id?: string,
  isActive: boolean,
  jobLocation: IAddress,
  companyBranchID: string,
  requirementCode: string,
  qualifications: any[],
  colleges: any[],
  marksCriteria: string,
  talentRating: number,
  talentRatingValue?: string,
  personalityType: string,
  minimumExperience: number,
  maximumExperience: number,
  jobRole: string,
  packageOffered: number,
  requiredBefore: string,
  requiredFrom: string,
  isUrgent?: boolean,
  vacancy: string,
  technologies: ITechnology[],
  selectedTalents: any[],
  comment: string,
  salesPersonID?: string,
  salesPersonName?: string
  salesPerson?: ISalesperson
  termsAndConditions?: string
}

//**************************************COMPANY ENQUIRY******************************************************* */

export interface ICompanyEnquiry {
  id?: string,

  // Address.
  address: string,
  city: string,
  pinCode: string,
  country: ICountry,
  state: IState,

  // Maps.
  domains: IDomain[],
  technologies: ITechnology[],

  // Related table IDs.
  salesPersonID: string,
  companyBranchID: string,

  // Other fields.
  companyName: string,
  website: string,
  email: string,
  code: string,
  hrName: string,
  hrContact: string,
  founderName: string,
  vacancy: number,
  enquiryDate: string,
  enquiryType: string,
  enquirySource: string,
  jobRole: string,
  packageOffered: number,
  subject: string,
  message: string,
  totalBranches: number,
}

export interface ICompanyEnquiryDTO {
  id?: string,

  // Address.
  address: string,
  city: string,
  pinCode: string,
  country: ICountry,
  state: IState,

  // Maps.
  domains: IDomain[],
  technologies: ITechnology[],

  // Single model.
  salesPerson: ISalesperson

  // Other fields.
  companyName: string,
  website: string,
  email: string,
  code: string,
  hrName: string,
  hrContact: string,
  founderName: string,
  totalBranches: number,
  vacancy: number,
  enquiryDate: string,
  enquiryType: string,
  enquirySource: string,
  jobRole: string,
  packageOffered: number,
  subject: string,
  message: string,
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
}

export interface ICallRecordDTO {
  id?: string

  // Related tables.
  purpose: IPurpose
  outcome: IOutcome

  // Other fields.
  dateTime: string
  comment: string
  talentID: string
}

//**************************************OTHERS******************************************************* */
export interface IAddress {
  address: string,
  city: string,
  pinCode: string,
  country: ICountry,
  state: IState
}

export interface IState {
  id: string,
  name: string,
  countryID?: string
}

export interface ITechnology {
  id: string,
  language: string
  rating: number
}

export interface IDomain {
  id: string,
  domainName: string
}

interface ISalesperson {
  id: string
  firstName: string
  lastName: string
}

export interface ICountry {
  id: string,
  name: string
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
