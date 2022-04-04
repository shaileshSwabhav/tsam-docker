import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';

@Injectable({
  providedIn: 'root'
})
export class QuestionService {

  constructor(
    private http: HttpClient,
    private constant: Constant
  ) { }


  getAllQuestions(limit: any, offset: any): Observable<any> {
    return new Observable<any>((observer) => {
      this.http.get(`${this.constant.BASE_URL}/question/${limit}/${offset}`,
        { observe: 'response' }).
        subscribe(data => {
          observer.next(data)
        },
          error => {
            observer.error(error)
          }
        )
    })
  }


  addQuestion(ques: question): Observable<any> {
    return new Observable<any>(observer => {
      this.http.post(`${this.constant.BASE_URL}/question`, ques).
        subscribe(data => {
          observer.next(data)
        },
          error => {
            observer.error(error)
          })
    })
  }

  updateQuestionByID(id: string, ques: question): Observable<any> {
    return new Observable<any>(observer => {
      this.http.put(`${this.constant.BASE_URL}/question/${id}`, ques).subscribe
        (data => {
          observer.next(data);
        },
          error => {
            observer.error(error);
          })
    })
  }

  deleteQuestionByID(id: string): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.delete(`${this.constant.BASE_URL}/question/${id}`).subscribe
          (
            data => {
              observer.next(data)
            },
            error => {
              observer.error(error)
            }
          )
      }
    )
  }


  getQuestionByID(id: string): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.get(`${this.constant.BASE_URL}/question/${id}`).subscribe
          (
            data => {
              observer.next(data)
            },
            error => {
              observer.error(error)
            }
          )
      }
    )
  }

  searchQuestion(conditions, limit: any, offset: any): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.post(`${this.constant.BASE_URL}/question/search/${limit}/${offset}`,
          conditions, { observe: "response" }).subscribe
          (
            data => {
              observer.next(data)
            },
            error => {
              observer.error(error)
            }
          )
      }
    )
  }

}


export interface question {
  id: string,
  question: string,
  subject: string,
  difficulty: string,
  options: [
    {
      id: string,
      option: string,
      status: boolean
    }
  ]

}