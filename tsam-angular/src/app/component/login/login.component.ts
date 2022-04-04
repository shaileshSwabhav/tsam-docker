import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { MainService } from 'src/app/service/main.service';
import { IMenu, MenuService } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Router } from '@angular/router';
import { UserloginService } from 'src/app/service/login/userlogin.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { AdminService } from 'src/app/service/admin/admin.service';
import { HttpHeaders } from '@angular/common/http';
import { environment } from 'src/environments/environment';

@Component({
      selector: 'app-login',
      templateUrl: './login.component.html',
      styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

      loginForm: FormGroup;
      visibility: string;
      input: string;
      menus: IMenu[] = [];
      roleID: string;
      tenantID: string;
      clicked: boolean


      isLocal: boolean

      constructor(
            private formbuilder: FormBuilder,
            private util: UtilityService,
            private menuService: MenuService,
            private spinnerService: SpinnerService,
            private localService: LocalService,
            private loginService: UserloginService,
            private router: Router
      ) {
            this.clicked = false
            this.spinnerService.loadingMessage = "Logging In..."
            this.checkLogin()
            this.isLocal = false
            if (environment.production === false) {
                  this.isLocal = true
            }
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {
            this.createForm();
            this.visibility = "visibility_off";
            this.input = "password";
      }

      // create form using Formbuilder.
      private createForm(): void {
            this.loginForm = this.formbuilder.group({
                  email: ['', Validators.required],
                  password: ['', Validators.required]
            });
      }

      // Return Form Control.
      get formcontrol() {
            return this.loginForm.controls;
      }

      visibilityToggle() {
            if (this.visibility == "visibility") {
                  this.visibility = "visibility_off";
                  this.input = "password";
                  return;
            }
            this.visibility = "visibility";
            this.input = "text";
      }

      checkLogin() {
            this.spinnerService.loadingMessage = "Loading..."
            if (this.localService.getJsonValue('credentialID') != "") {

                  this.loginService.validateSession().subscribe((response: any) => {
                        // console.log(response);

                        if (response) {
                              if (this.localService.getJsonValue("roleName") == "Talent") {
                                    this.util.rediectTo("my-batches")
                                    return
                              }
                              if (this.localService.getJsonValue("menus")?.length > 0) {
                                    let menu: IMenu = this.localService.getJsonValue("menus")[0]
                                    this.router.navigateByUrl(menu?.url)
                              }
                        }
                  }, (err: any) => {

                        console.error(err);
                  })
            }
      }

      // userLogin does not hide the spinner as it contains async calls to role which again has an async call to menus
      // when spinner is hide in userLogin there is time gap between the redirection to dashboard page
      userLogin(emailID?: string, password?: string) {
            if (emailID) {
                  this.loginForm.setValue(
                        {
                              "email": emailID,
                              "password": password,
                        }
                  )
            }
            this.clicked = true
            this.spinnerService.loadingMessage = "Logging In..."

            this.loginService.userLogin(this.loginForm.value).subscribe((response) => {
                  this.clicked = false
                  this.roleID = response.roleID;
                  this.tenantID = response.tenantID;
                  this.loginService.setLoginSession(response)
                  this.getRoleName();
                  // 
            }, err => {
                  console.error(err);

                  this.clicked = false
                  if (err.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert("Error: Invalid Username or Password")
            })
      }

      getAllMenusByRole(): void {
            // 
            this.menuService.getAllMenusByRole(this.tenantID, this.roleID).subscribe((data) => {
                  this.menus = data.body;
                  console.log("menus",this.menus)
                  this.localService.setJsonValue("menus", this.menus)
                  if (this.menus.length > 0) {
                        this.util.rediectTo(this.menus[0]?.url)
                  }

            }, (error) => {
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            });
      }

      getRoleName(): void {
            // 
            this.loginService.getRole(this.roleID).subscribe((data) => {
                  let role: any = data.body;
                  this.localService.setJsonValue("roleName", role.roleName);
                  this.getAllMenusByRole();
            }, (error) => {

                  console.error(error);
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            });
      }
      /*setLoginCookie(value: any) {
       const expiryTime = new Date();
       expiryTime.setHours(expiryTime.getHours() + 4)
       console.log(expiryTime);
       this.cookieService.set("token", value, {expires:  expiryTime})
 
       // if secure is set to true then cookie can only be set for secure connection i.e. https
       // cookie wont be set for http connection
       // const secureFlag = true
       // this.cookieService.set("token", value, {expires:  expiryTime, secure: secureFlag})
       }     */

}
