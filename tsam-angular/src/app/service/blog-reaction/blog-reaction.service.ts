import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BlogReactionService {

  blogReactionURL: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.blogReactionURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-reaction`
  }

  // Get all blog reactions.
  getAllBlogReactions(params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogReactionURL}`, 
      {params: params, headers: httpHeaders});
  }

  // Add new blog reaction.
  addBlogReaction(blogReaction: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.blogReactionURL}/credential/${credentialID}`, 
    blogReaction, { headers: httpHeaders });
  }

  // Update blog reaction.
  updateBlogReaction(blogReaction: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.blogReactionURL}/${blogReaction.id}/credential/${credentialID}`, 
    blogReaction, { headers: httpHeaders });
  }

  // Delete blog reaction.
  deleteBlogReaction(blogReactionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.blogReactionURL}/${blogReactionID}/credential/${credentialID}`, 
      { headers: httpHeaders });
  }
}

// Blog reaction.
export interface IBlogReaction {
  id?: string,

	// Related table IDs.
  reactorID: string
  blogID?: string
  replyID?: string

	// Flags.
  isClap: boolean
}

// Blog reaction DTO.
export interface IBlogReactionDTO {
  id?: string,

	// Related table IDs.
  blogID: string
  replyID: string

  // Singe models.
  reactor: IAuthor

	// Flags.
  isClap: boolean
}

// Others.
export interface IAuthor {
  id?: string       
  firstName: string
  lastName: string 
  role: any    
  image: string
}
