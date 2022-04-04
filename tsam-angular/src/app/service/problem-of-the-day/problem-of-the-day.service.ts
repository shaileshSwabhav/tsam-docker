import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProblemOfTheDayService {

  tenantURL: string

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // Get problems of the day.
  getProblemsOfTheDay(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.tenantURL}/problem-of-the-day`, {}, 
    { params: params, headers: httpHeaders });
  }

  // Get problems of the previous days.
  getProblemsOfThePreviousDay(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.tenantURL}/problem-of-the-day-previous`, {}, 
    { params: params, headers: httpHeaders });
  }

  // Get leader board for problems of the day.
  getLeaderBoardPotd(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.tenantURL}/leader-board-potd/talent/${talentID}`, {}, 
    { params: params, headers: httpHeaders });
  }

  // Get leader board for all questions.
  getLeaderBoard(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.tenantURL}/leader-board/talent/${talentID}`, {}, 
    { params: params, headers: httpHeaders });
  }
}

export interface IProblemOfTheDayQuestion{
  id?: string
  label: string         
  level: number  
  levelName: string    
  levelClass?: string  
  problemClass?: string
  attemptedByCount: number
  solvedByCount: number
  successRatio: string
}

export interface IPerformer{
  firstName: string         
  lastName: string  
  totalScore: number    
  rank: number  
  image?: string
}

export interface ILeaderBoard{
  selfPerformer: IPerformer
  allPerformers: IPerformer[]
}


