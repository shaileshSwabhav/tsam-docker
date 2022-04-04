import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { map } from 'rxjs/operators';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class TechnologyService {

  technologyUrl: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.technologyUrl = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/technology`
  }

  // Get all technologies.
  getAllTechnologies(limit: number, offset: number, params?:any): Observable<any> {
    params.limit = limit
    params.offset = offset
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.technologyUrl}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new technology.
  addTechnology(technology: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.technologyUrl}`, 
      technology, { headers: httpHeaders });
  }

  // Update technology.
  updateTechnology(technology: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.technologyUrl}/${technology.id}`, 
      technology, { headers: httpHeaders });
  }

  // Delete technology.
  deleteTechnology(designationID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.delete<any>(`${this.technologyUrl}/${designationID}`, 
      { headers: httpHeaders });
  }
}

//technology interface
export interface ITechnology {
  id?: string,
  language: string
  rating: number
}

export interface ISearchFilterField{
  propertyName: string
  propertyNameText: string
  valueList: any[]
}
