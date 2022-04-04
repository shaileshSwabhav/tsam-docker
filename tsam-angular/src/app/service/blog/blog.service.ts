import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { IBlogReactionDTO } from '../blog-reaction/blog-reaction.service';
import { IBlogTopic } from '../blog-topic/blog-topic.service';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  blogURL: string;

  constructor(
    private _http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.blogURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog`
  }

  // Get all blogs.
  getAllBlogs(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogURL}/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders, observe: "response" })
  }

  // Get blog by id.
  getBlogByID(blogID: string, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.blogURL}/${blogID}`, 
      {params: params, headers: httpHeaders})
  }

  // Get all latest blog snippets.
  getAllLatestBlogSnippet(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-snippet-latest/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders, observe: "response" })
  }

  // Get all trending blog snippets.
  getAllTrendingBlogSnippet(limit: any, offset: any, params?:HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-snippet-trending/limit/${limit}/offset/${offset}`, 
      {params: params, headers: httpHeaders, observe: "response" })
  }

  // Add new blog.
  addBlog(blog: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.post<any>(`${this.blogURL}/credential/${credentialID}`, 
    blog, { headers: httpHeaders })
  }

  // Update blog.
  updateBlog(blog: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.blogURL}/${blog.id}/credential/${credentialID}`, 
    blog, { headers: httpHeaders })
  }

  // Update blog flags.
  updateBlogFlags(blog: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.put<any>(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/blog-flags-update/credential/${credentialID}`, 
    blog, { headers: httpHeaders })
  }

  // Delete blog.
  deleteBlog(blogID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID:string = this.localService.getJsonValue('credentialID')
    return this._http.delete<any>(`${this.blogURL}/${blogID}/credential/${credentialID}`, 
      { headers: httpHeaders })
  }
}

// Blog.
export interface IBlog {
  id?: string,

	// Related table IDs.
  authorID: string,

	// Other fields.
  content: string,      
  title: string,        
  description: string,  
  timeToRead: number,   
  publishedDate: string,
  bannerImage: string

	// Flags.
  isVerified: boolean
  isPublished: boolean

	// Maps.
  blogTopics: IBlogTopic[]
}

// Blog DTO.
export interface IBlogDTO {
  id?: string,

	// Single model.
  author: IAuthor,

	// Other fields.
  content: string,      
  title: string,        
  description: string,  
  timeToRead: number,   
  publishedDate: string,
  bannerImage: string
  blogViewCount: number

	// Flags.
  isVerified: boolean
  isPublished: boolean
  isLoggedInClap: boolean

	// Maps.
  blogTopics: IBlogTopic[]

  // Mutiple field.
  reactions: IBlogReactionDTO[]
  clapReactions: IBlogReactionDTO[]
  slapReactions: IBlogReactionDTO[]
}

// Blog Snippet DTO.
export interface IBlogSnippetDTO {
  id?: string,

	// Single model.
  author: IAuthor,

	// Other fields.
  title: string,        
  description: string,  
  timeToRead: number,   
  publishedDate: string,
  bannerImage: string
  clapCount: number
  slapCount: number

	// Maps.
  blogTopics: IBlogTopic[]
}

// Others.
export interface IAuthor {
  id?: string       
  firstName: string
  lastName: string 
  role: any    
  image: string
}

