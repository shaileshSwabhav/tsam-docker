import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class FacultyDashboardService {

	private dashboardURL: string

	constructor(
		private constant: Constant,
		private localService: LocalService,
		private http: HttpClient
	) { 
		this.dashboardURL = `${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/dashboard/faculty`
	}

	getFacultyBatchDetails(facultyID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

		return this.http.get(`${this.dashboardURL}/${facultyID}/faculty-batch-details`, 
		{ params: params, headers: httpHeaders, observe: "response" });

	}

	getOngoingBatchDetails(facultyID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

		return this.http.get(`${this.dashboardURL}/${facultyID}/ongoing-batch-details`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}

	getBarchartData(facultyID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

		return this.http.get(`${this.dashboardURL}/${facultyID}/barchart`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}
	
	getPiechartData(facultyID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

		return this.http.get(`${this.dashboardURL}/${facultyID}/piechart`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}

	getTaskList(params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		let credentialID = this.localService.getJsonValue('credentialID')

		return this.http.get(`${this.dashboardURL}/credential/${credentialID}/task-list`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}

	// Get avaerage weekly rating for talent to faculty.
	getWeeklyAverageRating(batchID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		let credentialID = this.localService.getJsonValue('credentialID')
		return this.http.get(`${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/batch/${batchID}/faculty-weekly-rating`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}

	// Get avaerage weekly rating for talent to faculty for all weeks.
	getAllWeeksAverageRating(batchID: string, params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		let credentialID = this.localService.getJsonValue('credentialID')
		return this.http.get(`${this.constant.BASE_URL}/tenant/${this.localService.getJsonValue("tenantID")}/batch/${batchID}/weekly-feedback-rating`, 
		{ params: params, headers: httpHeaders, observe: "response" });
	}
}

export interface IFacultyBatchDetails {
	ongoingBatches: number
	upcomingBatches: number
	finishedBatches: number
	completedTrainingHrs: string
	totalStudents: number
}

export interface IOngoingBatchDetails {
	batchID: string
	courseName: string
	batchName: string
	totalSession: number
	pendingSession: number
	totalStudents: number
}

export interface IBarchart {
	totalStudents: number
	fresher: number
	professional: number
	studentsPlaced: number
}

export interface IPiechart {
	projectName: string
	totalCount: number
	hours: number
}