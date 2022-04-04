import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ConceptDashboardService {

  private tenantURL: string
  
  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  // Get complex concepts.
	getComplexConcepts(params?: HttpParams): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
		return this.httpClient.get(`${this.tenantURL}/complex-concepts`, 
		{ params: params, headers: httpHeaders, observe: "response" })
	}

}
