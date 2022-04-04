import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BlogViewerService {

  
  blogViewURL: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.blogViewURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-viewer`
  }

  // Get all blog views.
  getAllBlogViews(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogViewURL}/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders});
  }

  // Add new blog view.
  addBlogView(blogView: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.blogViewURL}/credential/${credentialID}`, 
    blogView, { headers: httpHeaders });
  }

  // Update blog view.
  updateBlogView(blogView: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.blogViewURL}/${blogView.id}/credential/${credentialID}`, 
    blogView, { headers: httpHeaders });
  }

  // Delete blog view.
  deleteBlogView(blogViewID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.blogViewURL}/${blogViewID}/credential/${credentialID}`, 
      { headers: httpHeaders });
  }
}

// Blog view.
export interface IBlogView {
  id?: string,

	// Related table IDs.
  viewerID: string
  blogID: string
}
