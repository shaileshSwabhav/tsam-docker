import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class MenuService {
  private menuURL: string

  constructor(
    private httpClient: HttpClient,
    private constant: Constant,
    private localService: LocalService

  ) {
    this.menuURL = constant.BASE_URL
  }

  //get all menus
  getAllMenusByRole(tenantID: string, roleID: string): Observable<any> {
    let httpHeaders = new HttpHeaders({ 'token': this.localService.getJsonValue('token') })
    return this.httpClient.get(`${this.menuURL}/tenant/${tenantID}/menu/role/${roleID}`, { headers: httpHeaders, observe: "response" });
  }
}

//menu interface
export interface IMenu {
  id?: string,
  order: number,
  permission: IPermission,
  menuName: string,
  tenantID: string,
  menuID: string,
  roleID: string,
  url: string,
  menus: IMenuDTO[]
  isVisible: boolean
  class?: string
  opertaion?: string
}

//menu interface
export interface IMenuDTO {
  id?: string,
  order: number,
  permission: IPermission,
  menuName: string,
  tenantID: string,
  menuID: string,
  roleID: string,
  url: string,
  menus: IMenu[]
  isVisible: boolean
  class?: string
  parentMenu?: IParentMenu
}

// Parent menu.
export interface IParentMenu {
  id?: string,
  menuName: string
}


export interface IPermission {
  add: boolean,
  update: boolean,
  delete: boolean
}
