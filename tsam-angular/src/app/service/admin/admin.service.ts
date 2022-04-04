import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';
import { IBatch, IBatchSession, IBatchTopic } from '../batch/batch.service';
import { ICountry, IProject, IState } from '../general/general.service';
import { UserloginService } from '../login/userlogin.service';

@Injectable({
  providedIn: 'root'
})
export class AdminService {
  roleURL: string
  tenantURL: string
  httpHeaders: HttpHeaders
  credentialID: string

  constructor(
    private http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
    private _http: HttpClient,
    private userLoginService: UserloginService
  ) {
    this.roleURL = `${this._constant.BASE_URL}/tenant/${this._constant.TENANT_ID}/role`
    this.tenantURL = `${this._constant.BASE_URL}/tenant/${this._constant.TENANT_ID}`

    // this.credentialID = localService.getJsonValue('credentialID')
    // this.httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  }



  getAllRole(): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.roleURL}`, { headers: httpHeaders, observe: 'response' })

    // return new Observable<any>((observer) => {
    //   this.http.get(`${this.roleURL}`,
    //     { headers: httpHeaders, observe: 'response' }).
    //     subscribe(data => {
    //       observer.next(data)
    //     }, (error) => {
    //       observer.error(error)
    //     })
    // })
  }


  postRole(Role: RoleType): Observable<RoleType> {
    let data: string = JSON.stringify(Role);
    console.log(data)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post<RoleType>(`${this._constant.BASE_URL}/role`,
      data, { headers: httpHeaders });
  }


  getRolebyId(roleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get<RoleType>(`${this.roleURL}/${roleID}`, { headers: httpHeaders, observe: 'response' });

    // return new Observable<any>((observer) => {
    //   this.http.get(`${this.roleURL}/${roleID}`,
    //     { headers: httpHeaders, observe: 'response' }).
    //     subscribe((data: any) => {
    //       observer.next(data)
    //     },
    //       (error) => {
    //         observer.error(error)
    //       })
    // })
  }

  getRole(roleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    //call get role api
    return this._http.get(`${this._constant.BASE_URL}/tenant/${this._constant.TENANT_ID}/role/${roleID}`,
      { headers: httpHeaders, observe: "response" });
  }

  deleteRolebyId(roleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.delete(`${this._constant.BASE_URL}/role/${roleID}`,
      { headers: httpHeaders });
  }

  updateRolebyId(role: any, roleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this._http.put(`${this._constant.BASE_URL}/role/${roleID}`, role,
      { headers: httpHeaders });
  }

  getSearchedRole(constraints, limit: number, offset: number): Observable<any> {
    console.log(limit, offset)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })


    return this._http.post(`${this._constant.BASE_URL}/role/search/${limit}/${offset}`, constraints,
      { headers: httpHeaders, observe: 'response' });
  }

  // ============================================== Timesheets ==============================================
  // Add Timesheet
  addTimesheet(timesheet: ITimesheet): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.tenantURL}/new-timesheet/credential/${credentialID}`, timesheet,
      { headers: httpHeaders })
  }

  // Add Timesheet
  addMultipleTimesheets(timesheets: ITimesheet[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.tenantURL}/new-timesheets/credential/${credentialID}`, timesheets,
      { headers: httpHeaders })
  }

  // update timesheet
  updateTimesheet(timesheet: ITimesheet): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.tenantURL}/new-timesheet/${timesheet.id}/credential/${credentialID}`,
      timesheet, { headers: httpHeaders })
  }

  // update timesheet activity
  updateTimesheetActivity(timesheetActivity: ITimesheetActivity): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.tenantURL}/timesheet-activity/${timesheetActivity.id}/credential/${credentialID}`,
      timesheetActivity, { headers: httpHeaders })
  }


  // deletes timesheets.
  deleteTimesheet(timesheetID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.tenantURL}/new-timesheet/${timesheetID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  // Get specific timesheets.
  getTimesheets(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/new-timesheet/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ============================================== Events ==============================================
  addEvent(event: IEvent): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.tenantURL}/event/credential/${credentialID}`, event,
      { headers: httpHeaders })
  }

  updateEvent(event: IEvent): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.tenantURL}/event/${event.id}/credential/${credentialID}`, event,
      { headers: httpHeaders })
  }

  deleteEvent(eventID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.tenantURL}/event/${eventID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  getEvents(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/event/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getAllEvent(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/event`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getEvent(eventID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.tenantURL}/event/${eventID}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get count of supervisors.
  getSupervisorCount(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.tenantURL}/supervisor-count`,
      { params: params, headers: httpHeaders, observe: "response" });
  }
}

export interface ITimesheet {
  id?: string
  tenantID?: string
  date?: string
  departmentID?: string
  credentialID?: string
  isOnLeave?: boolean
  activities: ITimesheetActivity[]
}

export interface ITimesheetActivity {
  id?: string
  timesheetID?: string
  project?: IProject
  projectID?: string
  activity?: string
  subProject?: IProject
  subProjectID?: string
  hoursNeeded?: number
  batch?: IBatch
  batchID?: string
  // batchSession?: IBatchSession
  // batchSessionID?: string
  batchTopicID?: string
  batchTopic?: IBatchTopic
  isBillable?: boolean
  isCompleted?: boolean
  workDone?: string
  nextEstimatedDate?: string
}

export interface IRole {
  id?: string
  roleName: string
}

export interface RoleType {
  id: string,
  roleName: string,
  level: number
}

export interface IEvent {
  id: string
  title: string
  description: string
  entryFee: number
  isOnline: boolean
  fromDate: string
  toDate: string
  fromTime: string
  toTime: string
  totalHours: number
  lastRegistrationDate: string
  eventStatus: string
  eventMeetingLink: string
  isActive: boolean
  address?: string
  city?: string
  pinCode?: number
  country?: ICountry
  state?: IState
  eventImage?: string
  totalRegistrations?: number
}