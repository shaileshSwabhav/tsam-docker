import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Constant } from '../constant';
import { Observable } from 'rxjs';
import { LocalService } from '../storage/local.service';

@Injectable({
	providedIn: 'root'
})
export class UserloginService {

	loginUrl: string;
	menuExists = false
	tenantID: string
	httpHeaders: HttpHeaders

	TOKEN: string
	TENANT_ID: string
	ROLE_ID: string
	DEPARTMENT_ID: string
	CREDENTIAL_ID: string
	LOGIN_ID: string
	FIRST_NAME: string
	LAST_NAME: string
	EMAIL: string

	//tenantID used for testing purpose, should not be hardcoded
	constructor(
		private http: HttpClient,
		private constant: Constant,
		private localService: LocalService,
	) {
		this.tenantID = this.constant.TENANT_ID
		this.httpHeaders = new HttpHeaders({ 'Content-Type': 'application/json' });

		// login info
		this.TOKEN = this.localService.getJsonValue("token")
		this.TENANT_ID = this.localService.getJsonValue("tenantID")
		this.ROLE_ID = this.localService.getJsonValue("roleID")
		this.DEPARTMENT_ID = this.localService.getJsonValue("departmentID")
		this.CREDENTIAL_ID = this.localService.getJsonValue("credentialID")
		this.LOGIN_ID = this.localService.getJsonValue("loginID")
		this.FIRST_NAME = this.localService.getJsonValue("firstName")
		this.LAST_NAME = this.localService.getJsonValue("lastName")
		this.EMAIL = this.localService.getJsonValue("email")

	}

	// userLogin(login: any) {
	//       return new Promise((resolve, reject) => {
	//             let status = this.dummydata.loginUser(login.username, login.password);
	//             if (status != undefined) {
	//                   this.util.error(this.dummydata.setNavigationByRole(status));
	//                   resolve(status);
	//                   return
	//             }
	//             reject("Invalid Login")
	//       })
	// }

	userLogin(loginDetails: any): Observable<ITokenDTO> {
		let paramJSON = JSON.stringify(loginDetails)
		// console.log("Main Service -> ", paramJSON)
		return this.http.post<ITokenDTO>(`${this.constant.BASE_URL}/tenant/${this.tenantID}/login`, paramJSON, { headers: this.httpHeaders });
	}

	userLogout(tenantID: any, loginID: any, token: string, params?: HttpParams): Observable<any> {
		// console.log(tenantID, loginID);
		let httpHeaders = new HttpHeaders({ 'token': token })
		return this.http.post<any>(`${this.constant.BASE_URL}/tenant/${tenantID}/logout/credential/${loginID}`, {},
			{ headers: httpHeaders, params: params })
	}

	validateSession(): Observable<any> {
		let httpHeaders = new HttpHeaders({ 'token': this.TOKEN })
		return this.http.get<any>(`${this.constant.BASE_URL}/tenant/${this.tenantID}/session`, { headers: httpHeaders })
	}

	getRole(id: string): Observable<any> {
		//add token to header from local storage
		let httpHeaders = new HttpHeaders({ 'token': this.TOKEN })
		return this.http.get(`${this.constant.BASE_URL}/tenant/${this.constant.TENANT_ID}/role/${id}`,
			{ 'headers': httpHeaders, observe: "response" });
	}

	setLoginSession(value: any) {
		// console.log(value)
		this.localService.setJsonValue("token", value.token)
		this.localService.setJsonValue("firstName", value.firstName)
		this.localService.setJsonValue("email", value.email)
		this.localService.setJsonValue('tenantID', value.tenantID)
		this.localService.setJsonValue('roleID', value.roleID)
		this.localService.setJsonValue('loginSessionID', value.loginSessionID)
		// this.localService.setJsonValue('departmentID', value.departmentID)
		this.localService.setJsonValue('credentialID', value.credentialID)//id from credentials table
		this.localService.setJsonValue('loginID', value.loginID)// id from one of six parent role tables

		this.TOKEN = value.token
		this.TENANT_ID = value.tenantID
		this.ROLE_ID = value.roleID
		this.CREDENTIAL_ID = value.credentialID
		this.LOGIN_ID = value.loginID
		this.FIRST_NAME = value.firstName
		this.EMAIL = value.email
		if (value.lastName != null) {
			this.LAST_NAME = value.lastName
			this.localService.setJsonValue("lastName", value.lastName)
		}
		if (value.departmentID != null) {
			this.DEPARTMENT_ID = value.departmentID
			this.localService.setJsonValue('departmentID', value.departmentID)
		}
	}

}

export interface ITokenDTO {
	name: string
	id: string
	tenantID: string
	token: string
	email: string
	loginID: string
	roleID: string
}

export interface ILoginInfo {
	token: string
	firstName: string
	lastName?: string
	email: string
	tenantID: string
	isActive: boolean
	roleID: string
	departmentID?: string
	credentialID: string
	loginID: string
}