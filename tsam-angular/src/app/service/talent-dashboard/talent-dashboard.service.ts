import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class TalentDashboardService {

  private tenantURL: string
  
  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  //  ============================================ TALENT REPORT ============================================

  getTalentReport(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/talent-report/${talentID}`,
      { params: params, headers: httpHeaders})
  }

  // Get faculty feedback for talent dashboard for current week and previous week.
  getFacultyFeedbackForTalentWeekWiseDashboard(talentID: string, batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/talent-dashboard/faculty-feedback-week-wise/talent/${talentID}/batch/${batchID}`,
      { params: params, headers: httpHeaders})
  }

  // Get faculty feedback for talent dashboard.
  getFacultyFeedbackForTalentDashboard(talentID: string, batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/talent-dashboard/faculty-feedback/talent/${talentID}/batch/${batchID}`,
      { params: params, headers: httpHeaders})
  }

  // Get faculty feedback for leader board.
  getFacultyFeedbackForTalentLeaderBoard(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/talent-dashboard/faculty-feedback-leader-board/batch/${batchID}`,
      { params: params, headers: httpHeaders})
  }

  // Get batch sessions along with batch talent details.
  getBatchSessionWithBatchTalentDetails(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.tenantURL}/talent-dashboard/batch/${batchID}/session/talent/${talentID}`,
      { params: params, headers: httpHeaders})
  }

  // Get talent feedback rating for all talents.
	getTalentFeedbackRatingForAllTalents(batchID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		return this.httpClient.get(`${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/batch/${batchID}/talent-feedback-rating`, 
		{ params: params, headers: httpHeaders, observe: "response" })
	}

  // Get talent weekly rating for all talents.
	getTalentWeeklyRatingForAllTalents(batchID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		return this.httpClient.get(`${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/batch/${batchID}/talent-weekly-rating`, 
		{ params: params, headers: httpHeaders, observe: "response" })
	}

	// GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
	getTalentConceptRatingWithBatchTopicAssignment(batchID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		return this.httpClient.get(`${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/batch/${batchID}/talent-concept-rating-with-assignment`, 
		{ params: params, headers: httpHeaders, observe: "response" })
	}

}
//  ============================================ TALENT REPORT ============================================

export interface ITalentReport {
  faculty: ITalent
  week: string
  monday: IBatch[]
  tuesday: IBatch[]
  wednesday: IBatch[]
  thursday: IBatch[]
  friday: IBatch[]
  saturday: IBatch[]
  sunday: IBatch[]
  totalTrainingHours: number
  workingHours: Map<string, IWorkingHours>
}

export interface ITalent{
  id: string,
  firstName: string,
  lastName: string
}

export interface IBatch {
  id: string
  batchName: string
  batchStatus: string
  batchTimings: IBatchTiming[]
  totalDailyHours: number
}

export interface IBatchTiming {
  id: string
  batchID: string
  day: any
  fromTime: string
  toTime: string
}

export interface IWorkingHours {
  batchName: string
  totalHours: number
}


