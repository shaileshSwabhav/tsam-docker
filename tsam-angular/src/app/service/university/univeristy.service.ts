import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class UniveristyService {

  universityUrl: string
  tenantUrl: string

  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.tenantUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
    this.universityUrl = `${this.tenantUrl}/university`
    
  }

  //get all universities
  getAllUniversities(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.universityUrl}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  addUniversity(university: IUniversity): Observable<IUniversity> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.post<IUniversity>(`${this.universityUrl}`, university,
      { headers: httpHeaders })
  }

  addMultipleUniversities(universities: IUniversity[]): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.post<any>(`${this.tenantUrl}/universities`, universities,
      { headers: httpHeaders })
  }

  getUniversitiesByCountryList(countryID: string): Observable<IUniversity[]> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get<IUniversity[]>(`${this.tenantUrl}/country/${countryID}/university`,
      { headers: httpHeaders })
  }

  deleteUniversity(universityID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.delete(`${this.universityUrl}/${universityID}`,
      { headers: httpHeaders })
  }

  updateUniversity(university: IUniversity): Observable<IUniversity> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.put<IUniversity>(`${this.universityUrl}/${university.id}`, university,
      { headers: httpHeaders })
  }

}

// University interface.
export interface IUniversity {
  id: string
  universityName: string
  countryID?: string
  countryName?: string
  country: ICountry
}

export interface ICountry {
  id?: string
  name: string
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}




