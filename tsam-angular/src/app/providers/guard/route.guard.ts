import { Location } from '@angular/common';
import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { LocalService } from 'src/app/service/storage/local.service';
import { environment } from "src/environments/environment";

@Injectable({
  providedIn: 'root'
})
export class RouteGuard implements CanActivate {

  menuExists: boolean

  constructor(
    private localService: LocalService,
    private router: Router,
    private location: Location
  ) { }

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    // console.log("Guard -> " + next)
    // if (!environment.production) {
    //   return true
    // }

    if (!this.localService.getJsonValue("token")) {
      this.router.navigateByUrl("/login")
      // return false
    }

    let menus = this.localService.getJsonValue('menus')
    if (this.checkNavigation(next, menus, state)) {
      return true
    }
    alert("Not Authorized.")
    return false

  }

  // Need to change this -n
  checkNavigation(next: ActivatedRouteSnapshot, menus: any, state: RouterStateSnapshot): boolean {
    let currentURL: string = ""
    for (let i = 1; i < next.pathFromRoot.length; i++) {
      currentURL = currentURL + "/" + next.pathFromRoot[i].routeConfig.path
    }
    // console.log(currentURL)
    for (let i = 0; i < menus?.length; i++) {
      if (menus[i].url == currentURL) {
        return true
      }
      if (menus[i].menus?.length > 0) {
        if (this.checkNavigation(next, menus[i].menus, state)) {
          return true
        }
      }
    }
    return false
  }
}