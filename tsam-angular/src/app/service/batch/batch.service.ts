import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams, HttpResponse } from '@angular/common/http';
import { LocalService } from '../storage/local.service';
import { IAcademics, IFeedbackOptions as IFeedbackOption, IFeedbackQuestion } from '../general/general.service';
import { ISalesperson, ITalent } from '../talent/talent.service';
import { IResource } from '../resource/resource.service';
import { ITechnology } from '../technology/technology.service';
import { ICourseSession } from '../course/course.service';
import { IProgrammingAssignment } from '../programming-assignment/programming-assignment.service';
import { IModule, IModuleTiming } from '../module/module.service';


@Injectable({
  providedIn: 'root'
})
export class BatchService {

  private batchURL: string;
  private programmingProjectURL: string;

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.batchURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`
    this.programmingProjectURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-project`
  }

  // CRUD on Batch
  getBatches(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  getBatchList(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch-list`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  getBatchesForFaculty(facultyID: string, limit: number, offset: number): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/faculty/${facultyID}/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchesForSalesPerson(salespersonID: string, limit: number, offset: number): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/salesperson/${salespersonID}/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, observe: "response" })
  }

  // getTalentsForBatch gets all the talents for a particular batch
  getTalentsForBatch(batchID: string, limit: number, offset: number): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/limit/${limit}/offset/${offset}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatch(batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}`,
      { headers: httpHeaders, observe: "response" })
  }

  addBatch(batch: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}`, batch,
      { headers: httpHeaders })
  }

  // Adds selected talents to specified batch
  addTalentsToBatch(talents: any[], batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/talent/credential/${credentialID}`, talents,
      { headers: httpHeaders })
  }

  // Deletes selected talent from the batch
  deleteTalentFromBatch(talentID: string, batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/talent/${talentID}/credential/${credentialID}`,
      { headers: httpHeaders })
  }

  deleteBatch(batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}`,
      { headers: httpHeaders })
  }

  updateBatch(batch: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.batchURL}/${batch.id}`, batch,
      { headers: httpHeaders })
  }

  // Get upcoming batches.
  getUpcomingBatches(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}/upcoming`,
      { params: params, headers: httpHeaders, observe: 'response' });
  }

  // Get batch details.
  getBatchDetails(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}/details/${batchID}`,
      { params: params, headers: httpHeaders });
  }

  // ********************************BATCH SESSION CRUD********************************

  getSessionForBatch(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/old-session`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getSessionAndAssignmentForBatch(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic-assignments`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  addSessions(batchID: string, batchSession: IMappedSession[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/old-batch-session/credential/${credentialID}`, batchSession,
      { headers: httpHeaders, observe: "response" })
  }

  updateSessions(batchID: string, batchSession: IMappedSession[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.batchURL}/${batchID}/old-batch-sessions/credential/${credentialID}`, batchSession,
      { headers: httpHeaders, observe: "response" })
  }

  updateSession(batchID: string, batchSession: IMappedSession): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.batchURL}/${batchID}/old-batch-session/${batchSession.id}/credential/${credentialID}`, batchSession,
      { headers: httpHeaders, observe: "response" })
  }

  deleteSessionFromBatch(batchID: string, sessionID: string, batchSessionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/old-batch-session/${batchSessionID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchSessionList(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/sessions`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ********************************BATCH SESSION FEEDBACK CRUD********************************

  // faculty-batch-session-feedback
  getFacultyBatchSessionFeedback(batchID: string, sessionID: string, facultyID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${sessionID}/faculty/${facultyID}/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getAllFacultyBatchSessionFeedback(batchID: string, sessionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${sessionID}/faculty/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // talent-batch-session feedback
  getSpecifiedTalentBatchSessionFeedback(batchID: string, sessionID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${sessionID}/talent/${talentID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  getTalentBatchSessionFeedback(batchID: string, sessionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${sessionID}/talent/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getAllTalentBatchSessionFeedback(batchID: string, sessionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/session/${sessionID}/feedbacks`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  //GetTalentFeedbackDetails will return feedback for the specified talent
  getTalentFeedbackDetails(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
  // add
  addFacultySessionFeedbacks(batchID: string, sessionID: string, sessionFeedbacks: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/topic/${sessionID}/faculty/feedbacks`,
      sessionFeedbacks, { headers: httpHeaders, observe: "response" })
  }

  addTalentSessionFeedbacks(batchID: string, sessionID: string, sessionFeedbacks: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/batch-session/${sessionID}/talent/feedbacks`,
      sessionFeedbacks, { headers: httpHeaders, observe: "response" })
  }

  // delete
  deleteFacultyBatchSessionFeedback(batchID: string, sessionID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/topic/${sessionID}/talent/${talentID}/faculty/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  deleteTalentBatchSessionFeedback(batchID: string, sessionID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/topic/${sessionID}/faculty/talent/${talentID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  // getBatchSessionFeedbacksForFaculty(batchID: string): Observable<any> {
  //   return this._http.get(`${this.batchURL}/${batchID}/session/feedback`, 
  //   { headers: this.httpHeaders, observe: "response" })
  // }

  // updateSessionFeedback(batchID: string, sessionFeedbacks: any[]): Observable<any> {
  //   return this._http.put(`${this.batchURL}/${batchID}/session/feedback/credential/${this.localService.getJsonValue('credentialID')}`, sessionFeedbacks,
  //   { headers: this.httpHeaders, observe: "response" })
  // }

  // deleteSessionFeedback(batchID: string, sessionFeedbackID: string): Observable<any> {
  //   return this._http.delete(`${this.batchURL}/${batchID}/session/feedback/${sessionFeedbackID}/credential/${this.localService.getJsonValue('credentialID')}`, 
  //   { headers: this.httpHeaders, observe: "response" })
  // }

  // ************************************************************* BATCH FEEDBACK CRUD *************************************************************
  getAllTalentBatchFeedback(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getAllFacultyBatchFeedback(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/faculty/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getFacultyBatchFeedback(batchID: string, facultyID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/faculty/${facultyID}/feedback`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getTalentBatchFeedback(batchID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/${talentID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  getFacultyTalentBatchFeedback(batchID: string, talentID: string, facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/${talentID}/faculty/${facultyID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  getFacultyTalentBatchSessionFeedback(batchID: string, sessionID: string, talentID: string, facultyID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${sessionID}/talent/${talentID}/faculty/${facultyID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  addFacultyBatchFeedbacks(batchID: string, sessionFeedbacks: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/faculty/feedbacks/credential/${credentialID}`,
      sessionFeedbacks, { headers: httpHeaders, observe: "response" })
  }

  addTalentBatchFeedbacks(batchID: string, sessionFeedbacks: any[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/talent/feedbacks/credential/${credentialID}`,
      sessionFeedbacks, { headers: httpHeaders, observe: "response" })
  }

  deleteFacultyBatchFeedback(batchID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/talent/${talentID}/faculty/feedback/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  deleteTalentBatchFeedback(batchID: string, facultyID: string, talentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/${batchID}/faculty/${facultyID}/talent/${talentID}/feedback/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }


  // ************************************************************* AHA MOMENT CRUD *************************************************************

  getAllAhaMoments(batchID: string, batchTopicID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/topic/${batchTopicID}/aha-moment`,
      { headers: httpHeaders, observe: "response" })
  }

  addAhaMoment(batchID: string, sessionID: string, ahaMoments: IAhaMoment[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.batchURL}/${batchID}/topic/${sessionID}/aha-moment/credential/${credentialID}`,
      ahaMoments, { headers: httpHeaders, observe: "response" })
  }

  deleteAhaMoment(ahaMomentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.batchURL}/aha-moment/${ahaMomentID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // ============================================================= PROGRAMMING PROJECT =============================================================

  addProgrammingProject(batchProject: IProgrammingProject): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.programmingProjectURL}`,
      batchProject, { headers: httpHeaders, observe: "response" })
  }

  updateProgrammingProject(batchProject: IProgrammingProject): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.programmingProjectURL}/${batchProject.id}`,
      batchProject, { headers: httpHeaders, observe: "response" })
  }

  deleteProgrammingProject(batchProjectID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.programmingProjectURL}/${batchProjectID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getProgrammingProject(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.programmingProjectURL}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ============================================== BATCH PROJECT ==============================================


  addBatchProject(batchID: string, batchProject: IBatchProject): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/project`,
      batchProject, { headers: httpHeaders, observe: "response" })
  }

  updateBatchProject(batchID: string, batchProject: IBatchProject): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchURL}/${batchID}/project/${batchProject.id}`,
      batchProject, { headers: httpHeaders, observe: "response" })
  }

  deleteBatchProject(batchID: string, batchProjectID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.batchURL}/${batchID}/project/${batchProjectID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchProject(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/project`,
      { params: params, headers: httpHeaders, observe: "response" })
  }


  // ========================================== BATCH SESSION PROGRAMMING ASSINGMENT ==========================================

  addBatchSessionProgrammingAssignment(sessionAssignment: IBatchTopicAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}-topic/${sessionAssignment.batchSessionID}/programming-assignment`,
      sessionAssignment, { headers: httpHeaders, observe: "response" })
  }

  updateBatchSessionProgrammingAssignment(sessionAssignment: IBatchTopicAssignment): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchURL}-topic/${sessionAssignment.batchSessionID}/programming-assignment/${sessionAssignment.id}`,
      sessionAssignment, { headers: httpHeaders, observe: "response" })
  }

  deleteBatchSessionProgrammingAssignment(batchSessionID: string, topicAssignmentID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.batchURL}-topic/${batchSessionID}/programming-assignment/${topicAssignmentID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchSessionProgrammingAssignment(batchProjectID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}-topic/${batchProjectID}/programming-assignment`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getSesssionAssignmentList(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}/${batchID}/topic-assignment-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getResponseProgrammingAssignment(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}/${batchID}/topic/programming-assignment`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getBatchSessionWiseProgrammingAssignment(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.batchURL}/${batchID}/programming-assignment`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ========================================== BATCH SESSIONS TALENTS ==========================================

  addTalentAttendance(batchID: string, talentAttendance: IBatchSessionsTalent): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/batch-topic/${talentAttendance.batchSessionID}/talent/${talentAttendance.talentID}/attendance`,
      talentAttendance, { headers: httpHeaders, observe: "response" })
  }

  // updateTalentAttendance(batchID: string, talentAttendance: IBatchSessionsTalent): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
  //   return this.http.put(`${this.batchURL}/${batchID}/batch-topic/${talentAttendance.batchSessionID}/talent/${talentAttendance.talentID}/attendance/${talentAttendance.id}`,
  //     talentAttendance, { headers: httpHeaders, observe: "response" })
  // }

  getBatchTalentDetails(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent-details`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getBatchSessionTalents(batchID: string, batchSessionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/batch-topic/${batchSessionID}/talent`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  //return talents feedback to faculty per session 
  getBatchTalentsFeedback(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/session-feedbacks`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getTalentTopicDetails(batchID:string,talentID:string, params?:HttpParams):Observable<any>{
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/${talentID}/talent-details`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
  
  // getTalentSessionFeedbackList(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
  //   let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

  //   return this.http.get(`${this.batchURL}/${batchID}/talent/${talentID}/topic-details`,
  //     { params: params, headers: httpHeaders, observe: "response" })
  // }

  // Get all batch session talents for specific talent.
  getAllBatchSessionTalentsForTalent(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/batch-session/talent/${talentID}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // GetAverageRatingForTalent will get average rating for batch.
  GetAverageRatingForTalent(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent/${talentID}/talent-average-rating`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // ========================================== BATCH MODULES ==========================================

  getOldBatchModules(batchID: string, params?: HttpParams): Observable<HttpResponse<IBatchModule[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get<IBatchModule[]>(`${this.batchURL}/${batchID}/old-module`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getBatchTalentList(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/talent-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get one batch session talents for specific talent.
  getOneBatchSessionForOneTalent(batchID: string, batchSessionID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/batch-session/${batchSessionID}/talent/${talentID}`,
      { params: params, headers: httpHeaders })
  }

  getOneTalentOneBatchSessionFeedback(batchID: string, talentID: string, batchSessionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchURL}/${batchID}/batch-session/${batchSessionID}/talent/${talentID}/feedback`,
      { headers: httpHeaders, observe: "response" })
  }

  // ========================================== BATCH MODULES ==========================================
    
  addBatchModules(batchID: string, batchModules: IBatchModule[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/modules`, batchModules,
      { headers: httpHeaders, observe: "response" })
  }
    
  addBatchModule(batchID: string, batchModule: IBatchModule): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchURL}/${batchID}/module`, batchModule,
      { headers: httpHeaders, observe: "response" })
  }
    
  updateBatchModule(batchID: string, batchModule: IBatchModule): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchURL}/${batchID}/modules/${batchModule.id}`, batchModule,
      { headers: httpHeaders, observe: "response" })
  }
    
  deleteBatchModule(batchID: string, batchModuleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.delete(`${this.batchURL}/${batchID}/modules/${batchModuleID}`,
      { headers: httpHeaders, observe: "response" })
  }

  getBatchModules(batchID: string, params?: HttpParams): Observable<HttpResponse<IBatchModule[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get<IBatchModule[]>(`${this.batchURL}/${batchID}/modules`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  getBatchModulesWithAllFields(batchID: string, params?: HttpParams): Observable<HttpResponse<IBatchModule[]>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get<IBatchModule[]>(`${this.batchURL}/${batchID}/modules-all-fields`,
      { params: params, headers: httpHeaders, observe: "response" })
  }


}

export interface IBatch {
  id?: string
  batchName: string
  code: string
  startDate?: string
  estimatedEndDate?: string
  totalStudents?: number
  totalIntake: number
  batchStatus?: string
  isActive: boolean
  eligibility: any
  course: any
  salesPerson: ISalesperson
  faculty: any
  sessions: any[]
  batchTimings: IBatchTiming[]
  requirement: any
  isB2B: boolean
  brochure?: string
  logo?: string
  batchObjective: string
  totalSessionCount: number
  completedSessionCount: number
}

export interface IBatchTiming {
  id?: string
  batchID: string
  day: IDay
  fromTime: string
  toTime: string
}

export interface IDay {
  id?: string
  day: string
  order: number
  type: string
}

export interface IBatchList {
  id?: string
  batchName: string
  code: string
  courseID: string
}

export interface ITalentBatchDTO {
  id: string
  talentCode: string
  firstName: string
  lastName: string
  email: string
  contact: string
  technologies: ITechnology[]
  academic: IAcademics[]
  address: string
  pinCode: string
  city: string
  country: string[]
  state: string[]
}

export interface IMappedTalent {
  talentID: string
}

export interface IMappedSessionObject {
  batchID: string
  sessionID: string
  tenantID: string
}

export interface IMappedSession {
  id: string
  batchID: string
  courseSessionID: string
  startDate: string
  isCompleted: boolean
  order: number
  session: IBatchSession
  tenantID: string
  isFeedbackGiven: boolean
  isAttendanceGiven: boolean
  sessionAssignment: IProgrammingAssignment[]
  talentBatchSessionFeedback: IBatchSessionFeedback[]
}

export interface IBatchSession {
  id: string
  name: string
  order: number
  studentOutput: string
  hours: string
  sessionID: string
  subSessions: IBatchSession[]
  courseID: string
  isChecked: boolean
  viewSubSessionClicked: boolean
  cardColumn: string
  materialIcon: string
}

export interface IBatchTopic {
  id: string
  batch?: any
  moduleTopic?: any
  order: number
  totalTime: string
  isCompleted: boolean
  completedDate: string
}

// *********************************************************** FEEDBACK ***********************************************************

// -> batch-feedback
export interface IBatchFeedback {
  id?: string
  batchID: string
  talentID: string
  facultyID: string
  questionID: string
  optionID: string
  talent: ITalent
  answer: string
  question: IFeedbackQuestion
  option: IFeedbackOption
}

// batch-feedback
export interface IBatchFeedbackDTO {
  talent: ITalent
  batchFeedbacks: IBatchFeedback[]
  averageScore: number
  showFeedback: boolean
}

// talent batch feedback
export interface ITalentBatchFeedback {
  faculty: any
  batchFeedbacks?: IBatchFeedback[]
  feedbacks?: IBatchFeedbackDTO[]
  averageScore: number
  showFeedback: boolean
}

// faculty talent batch feedback
export interface IFacultyBatchFeedback {
  talent: ITalent
  batchFeedbacks: IBatchFeedback[]
  faculty: any
  averageScore: number
  showFeedback: boolean
}


// -> batch-session feedback
export interface IBatchSessionFeedback {
  id?: string
  batchID: string
  batchSessionID: string
  talentID: string
  facultyID: string
  talent: ITalent
  questionID: string
  optionID?: string
  answer: string
  question: IFeedbackQuestion
  option: IFeedbackOption
  batchSession?: IMappedSession
}

export interface ITalentBatchSessionFeebackDTO {
  talent: ITalent
  faculty: any
  sessionFeedbacks: IBatchSessionFeedback[]
  showFeedback: boolean
  averageScore: number
}

// faculty batch session feedback
export interface IFacultyBatchSessionFeedback {
  talent: ITalent
  sessionFeedbacks: IBatchSessionFeedback[]
  faculty: any
  showFeedback: boolean
  averageScore: number
}

// talent batch session feedback
export interface ITalentBatchSessionFeedback {
  faculty: any
  sessionFeedbacks?: IBatchSessionFeedback[]
  feedbacks?: ITalentBatchSessionFeebackDTO[]
  showFeedback: boolean
  averageScore: number
}

// *********************************************************** AHA MOMENTS ***********************************************************

// ahaMoments
export interface IAhaMoment {
  id?: string
  batch: IBatch
  session: IMappedSession
  faculty: any
  talent: ITalent
  feeling: any
  feelingLevel: any
  ahaMomentResponse: IAhaMomentResponse[]
}

export interface IAhaMomentResponse {
  id?: string
  // batch: IBatch
  // session: IMappedSession
  // faculty: any
  // talent: ITalent
  // feeling: any
  // feelingLevel: any
  ahaMomentID: string
  questionID: string
  question: IFeedbackQuestion
  response: string
}

// ****************************************************** BATCH PROJECT ******************************************************

export interface IProgrammingProject {
  id?: string
  projectName: string
  description: string
  projectType: string
  code: string
  isActive: boolean
  complexityLevel: number
  requiredHours: number
  sampleUrl?: string
  resourceType: string
  document?: string
  technologies: ITechnology[]
  resources: IResource[]
}

export interface IBatchProject {
  id?: string
  batchID?: string
  batch: IBatch
  programmingProjectID?: string
  programmingProject?: IProgrammingProject
  dueDate: string
  assignedDate: string
}

// ****************************************************** BATCH SESSION PROGRAMMING ASSIGNMENT ******************************************************

export interface IBatchTopicAssignment {
  id?: string
  batchID: string
  batchSessionID: string
  programmingQuestionID: string
  programmingQuestion: IProgrammingAssignment
  order: number
  dueDate: string
  isMarked?: boolean
  showDetails?: boolean
}

export interface IBatchSessionAssignment {
  id?: string
  batchID: string
  courseSessionID: string
  session: ICourseSession
  sessionProgrammingAssignment: IBatchTopicAssignment[]
  showDetails?: boolean
  totalAssignment: number
}

// ********************************************* BATCH SESSIONS TALENT ********************************************
export interface IBatchSessionsTalent {
  id?: string
  batchID?: string
  batchSessionID?: string
  talentID?: string
  batch?: IBatch
  batchSession?: IMappedSession
  talent?: ITalentBatchDTO
  isPresent: boolean
  attendanceDate: string
  averageRating?: number
}

// ********************************************* BATCH MODULES ********************************************
export interface IBatchModule {
  id?: string
  module?: IModule
  moduleID?: string
  facultyID?: string
  faculty?: any
  order: number
  isCompleted: boolean
  startDate: string
  estimatedEndDate?: string
  moduleTimings?: IModuleTiming[]
    
  // flags
  showDetails?: boolean
}

// // ********************************************* TALENT ASSIGNMENT SCORE ********************************************

// export interface ITalentAssignmentSubmission {
//   id: string
//   batchSessionProgrammingAssignment: IBatchSessionProgrammingAssignment
//   isAccepted: boolean
//   isChecked: boolean
//   acceptanceDate: string
//   facultyRemarks: string
//   score: number
//   solution: string
//   githubURL: string
//   talentAssignmentSubmissionUpload: ITalentAssignmentSubmissionUpload[]
// }

// export interface ITalentAssignmentSubmissionUpload {
//   id: string
//   talentAssignmentSubmissionID: string
//   imageURL: string
//   description: string
// }

// export interface ITalentAssignmentScore {
//   talent: ITalent
//   talentSubmission: ITalentAssignmentSubmission[]
// }

