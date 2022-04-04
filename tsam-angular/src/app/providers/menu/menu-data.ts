import { Injectable } from "@angular/core";
import { IMenu } from "src/app/service/menu/menu.service";

@Injectable({
    providedIn: 'root'
})
export class MenuData {
    public menus: IMenu[]

    constructor() { }
}
