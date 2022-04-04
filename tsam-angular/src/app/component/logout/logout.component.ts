import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { UserloginService } from 'src/app/service/login/userlogin.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.css']
})
export class LogoutComponent implements OnInit {

  constructor(
    private loginService: UserloginService,
    private router: Router,
    private localService: LocalService
  ) {
    this.logout()
  }

  ngOnInit(): void {
  }

  // Logout and end session.
  logout() {
    let tenantID = this.localService.getJsonValue('tenantID')
    let credentialID = this.localService.getJsonValue('credentialID')
    let token = this.localService.getJsonValue('token')

    let queryParams: any = {
      loginSessionID: this.localService.getJsonValue("loginSessionID")
    }
    this.localService.clearToken()
    // console.log("queryParams -> ", queryParams);

    this.router.navigateByUrl("/login")

    this.loginService.userLogout(tenantID, credentialID, token, queryParams).subscribe((response: any) => {
      console.log("Logout successful")
    }, err => {
      console.error("Error -" + err.error);
    })
  }

}
