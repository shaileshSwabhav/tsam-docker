import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../../constant';
import { LocalService } from '../../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class CollegeCampusService {

  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
  ) {
    // this.httpHeaderss = new HttpHeaders({ 'token': this.localService.getJsonValue('token') });
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`;
  }

  //*******************************CAMPUS DRIVE API CALLS********************************************* 

  // Add one campus drive.
  addCampusDrive(campusDrive: any): Observable<any> {
    // let talentJson:string = JSON.stringify(campusDrive);
    // console.log("talent" + talentJson)
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/campus-drive/credential/${credentialID}`, 
      campusDrive, { headers: httpHeaders });
  }

  // Get all campus drives.
  getCampusDrives(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/campus-drive/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Get one campus drive by code.
  getCampusDriveByCode(code: string): Observable<any> {
    return this._http.get(`${this.generalURL}/campus-drive/code/${code}`);
  }

  // Update campus drive.
  updateCampusDrive(campusDrive: any): Observable<any> {
    // let talentJson: string = JSON.stringify(talent);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/campus-drive/${campusDrive.id}/credential/${credentialID}`, 
    campusDrive, { headers: httpHeaders });
  }

  // Delete campus drive.
  deleteCampusDrive(id: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/campus-drive/${id}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

// *******************************CANDIDATE API CALLS********************************
  // Get all candidates by campus drive.
  getCandidatesByCampusDrive(campusDriveID: any, limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.generalURL}/candidate/campus-drive/${campusDriveID}/limit/${limit}/offset/${offset}`, 
    { params: params, headers: httpHeaders, observe: "response"  });
  }

  // Add candidate.
  addCandidate(candidate: any, campusDriveID: any): Observable<any> {
    // let talentJson:string = JSON.stringify(candidate);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.post(`${this.generalURL}/candidate/campus-drive/${campusDriveID}/credential/${credentialID}`, 
    candidate, { headers: httpHeaders});
  }

  // Update candidate.
  updateCandidate(candidate: any, campusDriveID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/candidate/${candidate.campusTalentRegistrationID}/campus-drive/${campusDriveID}/credential/${credentialID}`, 
    candidate, { headers: httpHeaders });
  }

  // Delete candidate.
  deleteCandidate(candidateID: any, campusDriveID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.generalURL}/candidate/${candidateID}/campus-drive/${campusDriveID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }

  // Update multiple fields of candidates.
  updateMultipleCandidate(updateMultipleCandidate: IUpdateMultipleCandidate) {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })  
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.generalURL}/candidate/campus-drive/${updateMultipleCandidate.campusDriveID}/credential/${credentialID}`, 
    updateMultipleCandidate, { headers: httpHeaders });
  }
}

export interface IUpdateMultipleCandidate {
  isTestLinkSent: boolean
  hasAttempted: boolean
  result: string
  campusDriveID: string
  campusTalentRegistrationIDs: string[]
}

export interface ICampusDrive {
  id?: string

	// Maps.
  salesPeople: ISalesperson[]        
  faculties: IFaculty[]         
  developers: IDeveloper[]       
  collegeBranches: ICollegeBranch[]  

	// Flags.
  cancelled: boolean

  // Other fields.
  campusName: string             
  description: string          
  location: string                
  code: string                    
  totalRequirements: number       
  campusDate: string              
  studentRegistrationLink: string 
}

export interface ICampusDriveDTO {
  id?: string

	// Maps.
  salesPeople: ISalesperson[]        
  faculties: IFaculty[]         
  developers: IDeveloper[]       
  collegeBranches: ICollegeBranch[]  

	// Flags.
  cancelled: boolean

  // Other fields.
  campusName: string             
  description: string          
  location: string                
  code: string                    
  totalRegisteredCandidates: number
  totalAppearedCandidates: number 
  totalRequirements: number       
  campusDate: string              
  studentRegistrationLink: string 
}

//**************************************OTHERS******************************************************* */

export interface ISalesperson {
  id?: string,
  firstName: string,
  lastName: string
}

export interface IFaculty{
  id: string,
  firstName: string,
  lastName: string
}

export interface IDeveloper{
  id: string,
  firstName: string,
  lastName: string
}

export interface ICollegeBranch {
  id?: string
  branchName: string
  code: string
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}


