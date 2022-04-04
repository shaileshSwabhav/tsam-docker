import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';

@Injectable({
  providedIn: 'root'
})
export class OptionService {

  constructor(
    private http: HttpClient,
    private constant: Constant
  ) { }

  getAllOptions(): Observable<any> {
    return new Observable<any>((observer) => {
      this.http.get(`${this.constant.BASE_URL}/option`).
        subscribe(data => {
          observer.next(data)
        },
          error => {
            observer.error(error)
          }
        )
    })
  }




  updateOptionByID(id: string, option: option): Observable<any> {
    return new Observable<any>(observer => {
      this.http.put(`${this.constant.BASE_URL}/option/${id}`, option).subscribe
        (data => {
          observer.next(data);
        },
          error => {
            observer.error(error);
          })
    })
  }

  deleteOptionByID(id: string): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.delete(`${this.constant.BASE_URL}/option/${id}`).subscribe
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


  getOptionByID(id: string): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.get(`${this.constant.BASE_URL}/option/${id}`).subscribe
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

  getOptionByQuestion(quesId: string): Observable<any> {
    return new Observable<any>(
      observer => {
        this.http.get(`${this.constant.BASE_URL}/option/question/${quesId}`).subscribe
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


export interface option {
  id: string,
  option: string,
  status: boolean

}
