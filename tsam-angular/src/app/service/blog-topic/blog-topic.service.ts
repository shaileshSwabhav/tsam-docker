import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BlogTopicService {

  blogTopicURL: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.blogTopicURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-topic`
  }

  // Get all blog topics.
  getAllBlogTopics(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogTopicURL}/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders, observe: "response" });
  }

  // Add new blog topic.
  addBlogTopic(blogTopic: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.blogTopicURL}/credential/${credentialID}`, 
    blogTopic, { headers: httpHeaders });
  }

  // Update blog topic.
  updateBlogTopic(blogTopic: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.blogTopicURL}/${blogTopic.id}/credential/${credentialID}`, 
    blogTopic, { headers: httpHeaders });
  }

  // Delete blog topic.
  deleteBlogTopic(blogTopicID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.blogTopicURL}/${blogTopicID}/credential/${credentialID}`, 
      { headers: httpHeaders });
  }
}

export interface IBlogTopic {
  id?: string,
  name: string
}
