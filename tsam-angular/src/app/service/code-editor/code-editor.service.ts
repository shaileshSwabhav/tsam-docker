import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class CodeEditorService {

  // **************************** Paiza Compiler *****************************************
  sendCodeURLPaiza: string
  receiveOutputURLPaiza: string
  statusURLPaiza: string

  // **************************** Codex Compiler *****************************************
  sendAndReceiveURLCodex: string
  tempAccessToken: string

  // **************************** Sphere Compiler *****************************************
  sendCodeURLSphere: string
  outputURLSphere: string
  resultURLSphere: string
  addJudgeURLSphere: string
  addProblemURLSphere: string
  addTestCaseURLSphere: string
  addSubmissionURLSphere: string

  // **************************** Remote Code Compiler *****************************************
  sendCodeURLRemoteCode: string

  constructor(
    private http: HttpClient,
    private localService: LocalService,
    private constant: Constant,
  ) { 

  // **************************** Paiza Compiler *****************************************
    this.sendCodeURLPaiza = "http://api.paiza.io:80/runners/create"
    this.receiveOutputURLPaiza = "http://api.paiza.io:80/runners/get_details"
    this.statusURLPaiza = "http://api.paiza.io:80/runners/get_status"


  // **************************** Codex Compiler *****************************************
    this.sendAndReceiveURLCodex = "https://codexweb.netlify.app/.netlify/functions/enforceCode"

  // **************************** Sphere Compiler *****************************************
    this.tempAccessToken = "664b16978717b76aabe8bdcc95c53759"
    this.sendCodeURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/compiler-send-code-sphere`
    this.outputURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/compiler-send-code-sphere-output`
    this.resultURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/compiler-send-code-sphere-result`
    this.addJudgeURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/add-judge-sphere`
    this.addProblemURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/add-problem-sphere`
    this.addTestCaseURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/add-test-case-sphere`
    this.addSubmissionURLSphere = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/add-submission-sphere`

  // **************************** Remote Code Compiler *****************************************
    this.sendCodeURLRemoteCode = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/compiler-send-code-remote-code`
  }

  // **************************** Paiza Compiler *****************************************

  // Send the code to compiler API.
  sendCode(params?: HttpParams): Observable<any> {
    return this.http.post(this.sendCodeURLPaiza,
      {}, { params: params})
  }

  // Gets the status of code compilation from compiler API.
  getStatus(params?: HttpParams): Observable<any> {
    return this.http.get(this.statusURLPaiza,
      { params: params})
  }

  // Receives the output by id.
  receiveOutput(params?: HttpParams): Observable<any> {
    return this.http.get(this.receiveOutputURLPaiza,
      { params: params })
  }

  // **************************** Codex Compiler *****************************************

  // Send the code and receive output from compiler API.
  sendAndReceive(data: any): Observable<any> {
    return this.http.post(this.sendCodeURLPaiza, data)
  }

  // **************************** Sphere Compiler *****************************************
  
  // Send the code to compiler API.
  sendCodeSphere(sourceCodeAndInput: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.sendCodeURLSphere, sourceCodeAndInput, 
      { params: params, headers: httpHeaders})
  }

  // Gets the output of code compilation from compiler API.
  getOutputShpere(submissionID: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.outputURLSphere}/${submissionID}`,
      { params: params, headers: httpHeaders})
  }

  // Gets the result of code compilation from compiler API.
  getResultShpere(resultURL: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(`${this.resultURLSphere}`, {url: resultURL}, 
      { params: params, headers: httpHeaders, responseType:"text"})
  }

  // Add judge.
  addJudgeSphere(judge: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.addJudgeURLSphere, judge, 
      { params: params, headers: httpHeaders})
  }

  // Add problem.
  addProblemSphere(problem: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.addProblemURLSphere, problem, 
      { params: params, headers: httpHeaders})
  }

  // Add test case.
  addTestCaseSphere(testCase: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.addTestCaseURLSphere, testCase, 
      { params: params, headers: httpHeaders})
  }

  // Add submission.
  addSubmissionSphere(submission: any, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.addSubmissionURLSphere, submission, 
      { params: params, headers: httpHeaders})
  }

  // **************************** Remote Code Compiler *****************************************
  
  // Send the code to compiler API.
  sendCodeRemoteCode(sourceCodeAndInput: any, params?: HttpParams): Observable<any> {
    console.log(sourceCodeAndInput)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.post(this.sendCodeURLRemoteCode, sourceCodeAndInput, 
      { params: params, headers: httpHeaders})
  }
}
