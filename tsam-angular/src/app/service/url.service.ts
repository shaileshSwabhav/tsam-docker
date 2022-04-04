import { Location } from "@angular/common";
import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import { LocalService } from "./storage/local.service";

@Injectable({
    providedIn: 'root'
})
export class UrlService {
    // private _previousUrl: string | UrlTree
    // private _currentUrl: string | UrlTree

    constructor(private localService: LocalService,
        private router: Router,
        private location: Location) {
    }

    // public set currentUrl(url: string | UrlTree) {
    //     this._currentUrl = url
    //     this.localService.setJsonValue('currentUrl', url)
    // }

    // public get currentUrl(): string | UrlTree {
    // this._currentUrl = this.localService.getJsonValue('currentUrl')
    //     return this._currentUrl
    // }


    // key should be the component name from where the redirect is happening.
    public setUrlFrom(key: string, url: string) {
        // this._previousUrl = url
        this.localService.setJsonValue(key, url)
    }

    getUrlFrom(key: string): string {
        // console.log(this.localService.getJsonValue(key))
        return this.localService.getJsonValue(key)
    }

    // key should be the component name where we came from
    // for multiple options, create multiple 'back' button with different keys.
    public goBack(key: string) {
        let value = this.getUrlFrom(key)
        // if user tries to access the URL directly or some disrupted flow, use location.back().
        if (value === "") {
            this.location.back()
            return
        }
        this.router.navigateByUrl(value)
        this.localService.removeToken(key)
    }
}