import { Injectable } from '@angular/core';
import { UserloginService } from './login/userlogin.service';
import { StorageService } from './storage/storage.service';
import { Constant } from './constant';

@Injectable({
      providedIn: 'root'
})
export class MainService {

      constructor(
            private loginService: UserloginService,
            private storage: StorageService,
            private constant: Constant
      ) { }

      userLogin(param: any) {
            return this.loginService.userLogin(param);
      }

      userLogout(tenantID: any, loginID: any, token: string) {
            return this.loginService.userLogout(tenantID, loginID, token)
      }


      // rwdPriority(role: any, action: any) {
      //       if (role.level <= 1) {
      //             return true;
      //       }

      //       if (role.level == 3) {
      //             if (action == this.constant.DELETE) {
      //                   return false
      //             }
      //             return true
      //       }

      //       if (this.constant.DELETE == action) {
      //             if (role.level == 1) {
      //                   return true
      //             }
      //             return false
      //       }

      //       if (this.constant.READ == action) {
      //             return true;
      //       }

      //       if (this.constant.EDIT == action) {
      //             if (role.level)
      //       }
      // }
}
