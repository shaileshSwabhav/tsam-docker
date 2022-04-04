import { HttpHeaders, HttpClient, HttpParams, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import * as Q from 'q';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class ProgrammingQuestionService {

  programmingURL: string
  programmingTypeURL: string
  solutionURL: string
  generalURL: string
  httpHeaders: HttpHeaders
  credentialID: string
  testCaseURL: string
  tenantURL: string

  constructor(
    private http: HttpClient,
    private constant: Constant,
    private localService: LocalService,
  ) {
    this.programmingURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question`
    this.programmingTypeURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-type`
    this.solutionURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-solution`
    this.testCaseURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-test-case`
    this.tenantURL = `${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}`
  }

  getTopicProgrammingQuestion(params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.tenantURL}/topics/programming-questions`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Returns all the programming questions for practice in problem of the day format.
  getProgrammingQuestionsPractice(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/programming-question-practice/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Returns all the programming questions with limit and offset.
  getProgrammingQuestions(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.programmingURL}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Returns one programming question by id.
  getProgrammingQuestion(questionID: string, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.programmingURL}/${questionID}`,
      { params: params, headers: httpHeaders })
  }

  // Adds programming questions.
  addProgrammingQuestion(programmingQuestion: IProgrammingQuestion): Observable<HttpResponse<string>> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post<string>(`${this.programmingURL}/credential/${credentialID}`,
      programmingQuestion, { headers: httpHeaders, observe: "response" })
  }

  // Updates programming question and its options.
  updateProgrammingQuestion(programmingQuestion: IProgrammingQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.programmingURL}/${programmingQuestion.id}/credential/${credentialID}`,
      programmingQuestion, { headers: httpHeaders, observe: "response" })
  }

  // Update isActive of question.
  updateProgrammingQuestionIsActive(question: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')
    return this.http.put<any>(`${this.programmingURL}/is-active/${question.id}/credential/${credentialID}`,
      question, { headers: httpHeaders });
  }

  // Deletes programming questions.
  deleteProgrammingQuestion(programmingQuestionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.programmingURL}/${programmingQuestionID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // ========================================== PROGRAMMING QUESTION TYPE ==========================================  

  // Returns all the programming questions with limit and offset.
  getProgrammingQuestionTypes(limit: number, offset: number, params?: HttpParams): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })

    return this.http.get(`${this.programmingTypeURL}/limit/${limit}/offset/${offset}`,
      { params: params, headers: httpHeaders, observe: "response" })
  }

  // Adds programming questions.
  addProgrammingQuestionType(programmingQuestion: IProgrammingQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.post(`${this.programmingTypeURL}/credential/${credentialID}`,
      programmingQuestion, { headers: httpHeaders, observe: "response" })
  }

  // Updates programming question and its options.
  updateProgrammingQuestionType(programmingQuestion: IProgrammingQuestion): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.put(`${this.programmingTypeURL}/${programmingQuestion.id}/credential/${credentialID}`,
      programmingQuestion, { headers: httpHeaders, observe: "response" })
  }

  // Deletes programming questions.
  deleteProgrammingQuestionType(programmingQuestionID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID = this.localService.getJsonValue('credentialID')

    return this.http.delete(`${this.programmingTypeURL}/${programmingQuestionID}/credential/${credentialID}`,
      { headers: httpHeaders, observe: "response" })
  }

  // *******************************PROGRAMMING QUESTION SOLUTION API CALLS********************************

  // Get all programming question solutions by programming question.
  getSolutionsByQuestions(questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.solutionURL}/programming-question/${questionID}`,
      { headers: httpHeaders });
  }

  // Add programming question solution.
  addSolution(solution: any, questionID: any): Observable<any> {
    console.log(questionID)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.post(`${this.solutionURL}/programming-question/${questionID}/credential/${credentialID}`,
      solution, { headers: httpHeaders });
  }

  // get programming question solution by id
  getSolution(solutionID: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.solutionURL}/${solutionID}/programming-question/${questionID}`,
      { headers: httpHeaders });
  }

  // Update programming question solution.
  updateSolution(solution: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.put<any>(`${this.solutionURL}/${solution.id}/programming-question/${questionID}/credential/${credentialID}`,
      solution, { headers: httpHeaders });
  }

  // Delete programming question solution.
  deleteSolution(solutionID: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.delete<any>(`${this.solutionURL}/${solutionID}/programming-question/${questionID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }

  // ========================================== PROGRAMMING QUESTION TEST CASE ==========================================  

  // Get all programming question test cases by programming question.
  getTestCasesByQuestions(questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.testCaseURL}/programming-question/${questionID}`,
      { headers: httpHeaders });
  }

  // Get programming question test case by id
  getTestCase(testCaseID: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.http.get(`${this.testCaseURL}/${testCaseID}/programming-question/${questionID}`,
      { headers: httpHeaders });
  }

  // Add programming question test case.
  addTestCase(testCase: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.post(`${this.testCaseURL}/programming-question/${questionID}/credential/${credentialID}`,
      testCase, { headers: httpHeaders });
  }

  // Update programming question test case.
  updateTestCase(testCase: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.put<any>(`${this.testCaseURL}/${testCase.id}/programming-question/${questionID}/credential/${credentialID}`,
      testCase, { headers: httpHeaders });
  }

  // Delete programming question test case.
  deleteTestCase(testCaseID: any, questionID: any): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    let credentialID: string = this.localService.getJsonValue('credentialID')
    return this.http.delete<any>(`${this.testCaseURL}/${testCaseID}/credential/${credentialID}`,
      { headers: httpHeaders });
  }
}

export interface IProgrammingQuestion {
  id?: string
  label: string
  question: string
  inputFormat: string
  example: string
  constraints: string
  comment: string
  hasOptions: boolean
  isActive: boolean
  level: number
  levelName: string
  levelClass: string
  score: number
  timeRequired: number
  programmingQuestionTypes: IProgrammingQuestionType[]
  options: IProgrammingQuestionOption[]
  testCases: IProgrammingQuestionTestCase[]
  isAnswered: boolean
  answer: string
  programmingQuestionOptionID: string
  attemptedByCount: number
  solvedByCount: number
  successRatio: string
  solutionCount: number
  solutions: IProgrammingQuestionSolutionDTO[]
  programmingLanguage: IProgrammingLanguage
  solutonIsViewed: boolean
  hasAnyTalentAnswered: boolean
  isLanguageSpecific: boolean
  programmingLanguages: IProgrammingLanguage[]

  // flags
  isMarked: boolean
}

export interface IProgrammingQuestionOption {
  id?: string
  programmingQuestionID: string
  option: string
  isCorrect: boolean
  isActive: boolean
  order: number
  optionClass: string
}

export interface IProgrammingQuestionType {
  id?: string
  programmingType: string
}

export interface IProgrammingQuestionTalentAnswerQuestionDTO {
  id?: string
  label: string
  level: number
  score: number
  levelName: string
  levelClass: string
}

export interface IProgrammingQuestionListDTO {
  id?: string
  label: string
  level: number
}

export interface IProgrammingQuestionSolution {
  id?: string
  solution: string
  programmingQuestionID: string
  programmingLanguageID: string
}

export interface IProgrammingQuestionSolutionDTO {
  id?: string
  solution: string
  programmingQuestionID: string
  programmingLanguage: IProgrammingLanguage
}

export interface IProgrammingQuestionTestCase {
  id?: string
  programmingQuestionID: string
  input: string
  output: string
  explanation: string
  isActive: boolean
  isHidden: boolean
}

export interface IProgrammingQuestionIsActive {
  id?: string
  isActive: boolean
}

export interface IProgrammingLanguage {
  id?: string,
  name: string
  rating: number
}


