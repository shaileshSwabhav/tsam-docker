import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { UserloginService } from '../login/userlogin.service';


@Injectable({
  providedIn: 'root'
})

export class CollegeService {

  httpHeaders: HttpHeaders;
  tenantURL: string;
  private collegeURL: string
  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') });
    this.tenantURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`;
  }

  //***************************************** College Branch API calls********************************************************
  // Get all college branches
  getAllCollegeBranches(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.tenantURL}/college/branch/limit/${limit}/offset/${offset}`, 
    { params: params, headers: httpHeaders, observe: "response" });
  }

  // getAllBranchesOfCollege(collegeID: string, params?: HttpParams): Observable<any> {
  //   return this._http.get(`${this.generalURL}/college/${collegeID}/branch`, 
  // { params: params, headers: httpHeaders, observe: "response" });
  // }

  // getAllBranchesOfSalesPerson(limit: number, offset: number, params?: HttpParams): Observable<any> {
  //   let salesPersonID: string = this.localService.getJsonValue("loginID");
  //   if (!salesPersonID) {
  //     return
  //   }
  //   return this._http.get(`${this.generalURL}/college/branch/sales-person/${salesPersonID}/limit/${limit}/offset/${offset}`,
  //     { params: params, headers: httpHeaders, observe: "response" });
  // }

  getCollegeBranchByID(collegeBranchID: string, collegeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.tenantURL}/college/${collegeID}/branch/${collegeBranchID}`, 
    { headers: httpHeaders });
  }

  addCollegeBranch(collegeBranch: any): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call add college branch api
    return this._http.post<any>(`${this.tenantURL}/college/${collegeBranch.collegeID}/branch/credential/${credentialID}`, 
    collegeBranch, { headers: httpHeaders });
  }

  deleteCollegeBranch(collegeID: string, collegeBranchID: string): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call delete talent api
    return this._http.delete<any>(`${this.tenantURL}/college/${collegeID}/branch/${collegeBranchID}/credential/${credentialID}`, 
    { headers: httpHeaders });

  }

  updateCollegeBranch(collegeBranch: ICollegeBranch): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call add college branch api
    return this._http.put<any>(`${this.tenantURL}/college/${collegeBranch.collegeID}/branch/${collegeBranch.id}/credential/${credentialID}`, 
    collegeBranch, { headers: httpHeaders });

  }

  //***************************************** College API calls********************************************************

  // Add new college
  addCollege(college: any): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call add college api
    return this._http.post<any>(`${this.tenantURL}/college/credential/${credentialID}`, 
    college, { headers: httpHeaders });
  }

  // Add new college
  addMultipleColleges(colleges: ICollege[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    return this._http.post<any>(`${this.tenantURL}/colleges/credential/${credentialID}`, 
    colleges, { headers: httpHeaders });
  }

  updateCollege(college: ICollege): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    let talentJSON = JSON.stringify(college)
    console.log("talent json")
    console.log(talentJSON)

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call update talent api
    return this._http.put<any>(`${this.tenantURL}/college/${college.id}/credential/${credentialID}`, 
    college, { headers: httpHeaders });
  }


  getCollegeByID(collegeID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.get(`${this.tenantURL}/college/${collegeID}`, { headers: httpHeaders });
  }

  getAllColleges(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    
    return this._http.get(`${this.tenantURL}/college/limit/${limit}/offset/${offset}`, 
    { params: params, headers: httpHeaders, observe: "response" });
  }


  deleteCollege(collegeID: string): Observable<string> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //get credential id from local storage
    let credentialID: string = this.localService.getJsonValue('credentialID')

    //call delete talent api
    return this._http.delete<any>(`${this.tenantURL}/college/${collegeID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }
}

export interface ICollegeBranch {
  id?: string
  branchName: string
  code?: string
  address?: IAddress
  collegeID?: string
  tpoName?: string
  tpoContact?: string
  tpoAlternateContact?: string
  tpoEmail?: string
  collegeRating?: number
  email?: string
  salesPerson?: ISalesperson
  allIndiaRanking?: number
  university?: IUniversity
  // countryName?: string
  // stateName?: string
  // universityName?: string
}

export interface ICollege {
  id?: string,
  collegeName: string,
  code?: string,
  chairmanName?: string,
  chairmanContact?: string,
  collegeBranches: ICollegeBranch[],
}

export interface ICountry {
  id: string,
  name: string
}
export interface IUniversity {
  id: string,
  universityName: string
}

export interface IState {
  id: string,
  name: string,
  countryID?: string
}

export interface IAddress {
  address: string,
  city: string,
  pinCode: string,
  country: ICountry,
  state: IState,
}

export interface ISalesperson {
  id: string,
  firstName: string,
  lastName: string
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}



