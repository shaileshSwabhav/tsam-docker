import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { IDepartmentDTO } from '../department/department.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class TargetCommunityService {

  targetCommunityUrl: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.targetCommunityUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/target-community`
  }

  // Get all target communities.
  getAllTargetCommunities(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.targetCommunityUrl}/limit/${limit}/offset/${offset}`, 
    {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new target community.
  addTargetCommunity(targetCommunity: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.targetCommunityUrl}/credential/${credentialID}`, targetCommunity, 
    { headers: httpHeaders });
  }

  // Update target community.
  updateTargetCommunity(targetCommunity: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.targetCommunityUrl}/${targetCommunity.id}/credential/${credentialID}`, 
    targetCommunity, { headers: httpHeaders });
  }

  // Update target community for updating is target achieved field.
  updateTargetCommunityIsTargetAchieved(targetCommunity: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/target-community-achieved/credential/${credentialID}`, 
    targetCommunity, { headers: httpHeaders });
  }

  // Delete target community.
  deleteTargetCommunity(targetCommunityID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.targetCommunityUrl}/${targetCommunityID}/credential/${credentialID}`, 
    { headers: httpHeaders });
  }
}

export interface ITargetCommunity {
  id?: string

	// Maps.
  colleges: ICollege[]
  companies: ICompany[]
  courses: ICourse[]

	// Related table IDs.
  departmentID: string
  // credentialID: string
  functionID: string
  facultyID: string
  salesPersonID: string

	// Other fields.
  targetType: string
  studentType: string
  numberOfBatches: number
  targetStartDate: string
  targetEndDate: string
  isTargetAchieved: boolean
  targetStudentCount: number
  hours: number
  fees: number
  rating: number
  requiredTalentRating: number
  talentType: number        
  minExperienceYears: number
  maxExperienceYears: number
  salary: number  
  salaryInString: string            
  upSell: number   
  crossSell: number       
  referral: number      
  action: string         
}

export interface ITargetCommunityDTO {
  id?: string

	// Maps.
  colleges: ICollege[]
  companies: ICompany[]
  courses: ICourse[]

	// Related tables.
  department: IDepartmentDTO
  departmentID: string
  // credential: ICredential
  // credentialID: string
  function: ITargetCommunityFunction
  functionID: string
  faculty: IFaculty
  facultyID: string
  salesPerson: ISalesperson

	// Other fields.
  targetType: string
  studentType: string
  numberOfBatches: number
  targetStartDate: string
  targetEndDate: string
  isTargetAchieved: boolean
  targetStudentCount: number
  hours: number
  fees: number
  feesInString: string
  rating: number
  requiredTalentRating: number
  talentType: number        
  minExperienceYears: number
  maxExperienceYears: number
  salary: number   
  salaryInString: string            
  upSell: number   
  crossSell: number       
  referral: number      
  action: string  
}

export interface ICredential{
  id?: string      
  firstName: string
  lastName: string 
  roleID: string   
}

export interface ITargetCommunityFunction{
  id?: string
  functionName: string
  departmentID: string
}

export interface ICollege{
  id?: string        
  branchName: string
  code: string      
}

export interface ICompany{
  id?: string        
  branchName: string
  companyID: string
}

export interface ICourse{
  id?: string
  code: string
  name: string
  courseType: string
}

export interface IFaculty{
  id?: string       
  firstName: string
  lastName: string 
}

export interface ISalesperson {
  id?: string,
  firstName: string,
  lastName: string
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


