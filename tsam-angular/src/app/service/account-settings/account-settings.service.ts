import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class AccountSettingsService {

  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService,
  ) { 
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`;
  }

  // Verify password.
  verifyPassword(passwordChange: any): Observable<any> {
    // let talentJson:string = JSON.stringify(campusDrive);
    // console.log("talent" + talentJson)
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.post<any>(`${this.generalURL}/login/verify-password`, 
    passwordChange, { headers: httpHeaders });
  }

  // Change password.
  changePassword(passwordChange: any): Observable<any> {
    // let talentJson:string = JSON.stringify(campusDrive);
    // console.log("talent" + talentJson)
    let credentialID: string = this.localService.getJsonValue('credentialID')
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this._http.put<any>(`${this.generalURL}/login/change-password/credential/${credentialID}`, 
    passwordChange, { headers: httpHeaders });
  }
}

// Interface for password change.
export interface IPassswordChange {
  email: string      
  roleID: string      
  password: string    
}
