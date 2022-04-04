import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IBlogReactionDTO } from '../blog-reaction/blog-reaction.service';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BlogReplyService {

  blogReplyURL: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.blogReplyURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-reply`
  }

  // Get all blog replies.
  getAllBlogRelpies(blogID: string, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogReplyURL}/blog/${blogID}`, 
      {params: params, headers: httpHeaders});
  }

  // Add new blog reply.
  addBlogReply(blogReply: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.blogReplyURL}/credential/${credentialID}`, 
    blogReply, { headers: httpHeaders });
  }

  // Update blog reply.
  updateBlogReply(blogReply: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.blogReplyURL}/${blogReply.id}/credential/${credentialID}`, 
    blogReply, { headers: httpHeaders });
  }

  // Delete blog reply.
  deleteBlogReply(blogReplyID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.blogReplyURL}/${blogReplyID}/credential/${credentialID}`, 
      { headers: httpHeaders });
  }
}

// Blog reply.
export interface IBlogReply {
  id?: string,

	// Related table IDs.
  replierID: string
  blogID: string
  replyID: string

	// Other fields.
  reply: string

	// Flags.
  isVerified: boolean
}

// Blog reply DTO.
export interface IBlogReplyDTO {
  id?: string,

	// Related table IDs.
  blogID: string
  replyID: string

  // Singe models.
  replier: IAuthor

	// Mutiple field.
  replies: IBlogReplyDTO[]
  reactions: IBlogReactionDTO[]
  clapReactions: IBlogReactionDTO[]
  slapReactions: IBlogReactionDTO[]

	// Other fields.
  reply: string

	// Flags.
  isVerified: boolean
  isRepliesVisible: boolean
  isReplyToReplyFormVisible: boolean
  isReplyUpdateFormVisible: boolean
  isBlogReply: boolean
  isLoggedInClap: boolean
}

// Others.
export interface IAuthor {
  id?: string       
  firstName: string
  lastName: string 
  role: any    
  image: string
}

