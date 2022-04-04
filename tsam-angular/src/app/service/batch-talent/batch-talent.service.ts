import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/internal/Observable';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BatchTalentService {

  private batchTalentURL: string;
  private batchUrl: string = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch`

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.batchTalentURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch-talent`
  }

  // Gets list of all active talents in the given batch.
  getBatchTalentList(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchUrl}/${batchID}/talent-list`,
      { params: params, headers: httpHeaders, observe: "response" })
  }


  // Get all batch talents by batch id.
  getBatchTalentDetails(batchID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.batchTalentURL}/batch/${batchID}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Get batches of talent.
  getBatchesForTalent(talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/talent/${talentID}/batch`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Adds one batch talent.
  addBatchTalents(batchTalents: any[], batchID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.post(`${this.batchTalentURL}/batch/${batchID}`, batchTalents,
      { headers: httpHeaders })
  }

  // Updates one batch talent.
  updateBatchTalent(batchTalent: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.batchTalentURL}/${batchTalent.id}`, batchTalent,
      { headers: httpHeaders })
  }

  // Updates suspension date of one batch talent.
  updateSuspensionDateBatchTalent(batchTalent: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch-talent-suspension-date/${batchTalent.id}`, batchTalent,
      { headers: httpHeaders })
  }

  // Updates is active of one batch talent.
  updateIsActiveBatchTalent(batchTalent: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.put(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch-talent-is-active/${batchTalent.id}`, batchTalent,
      { headers: httpHeaders })
  }

  // Get batches of talent.
  getOneBatchTalentDetails(batchID: string, talentID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/batch/${batchID}/talent/${talentID}`,
      { params: params, headers: httpHeaders })
  }

   // Get all talents for batch sessions.
   getBatchSessionTalents(batchID: string, batchSessionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    
    return this.http.get(`${this.batchUrl}/${batchID}/batch-session/${batchSessionID}/talent`,
      { params: params, headers: httpHeaders, observe: "response" })
  }
}

export interface IBatchTalent {
  id?: string
  batchID: string
  talentID: string
  dateOfJoining: string
  isActive: boolean
  suspensionDate: string
}