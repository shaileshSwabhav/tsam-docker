import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { UserloginService } from '../login/userlogin.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class DepartmentService {

  departmentUrl: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
    private userLoginService: UserloginService
  ) {
    this.departmentUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/department`
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    // this.credentialID = localService.getJsonValue('credentialID')
  }

  // Get all departments.
  getAllDepartments(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.departmentUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new department.
  addDepartment(department: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.departmentUrl}`, 
      department, { headers: httpHeaders });
  }

  // Update department.
  updateDepartment(department: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.departmentUrl}/${department.id}`, 
      department, { headers: httpHeaders });
  }

  // Delete department.
  deleteDepartment(departmentID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.departmentUrl}/${departmentID}`, 
      { headers: httpHeaders });
  }
}

export interface IDepartment {
  id?: string
  name: string
  roleID: string
}

export interface IDepartmentDTO {
  id?: string
  name: string
  roleID: string
  role: IRole
}

export interface ITargetCommunityFunction {
  id?: string
  functionName: string
  departmentID: string
}

export interface ITargetCommunityFunctionDTO {
  id?: string
  functionName: string
  departmentID: string
  department?: IDepartmentDTO
}

export interface IRole {
  id?: string
  roleName: string
  level: number
  isEmployee: boolean
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}