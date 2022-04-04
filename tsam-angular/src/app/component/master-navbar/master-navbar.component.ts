import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { IMenu, IMenuDTO } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { Router } from '@angular/router';
import { UserloginService } from 'src/app/service/login/userlogin.service';
import { MatMenuTrigger } from '@angular/material/menu';
import { AccessLevel, Role } from 'src/app/service/constant';

@Component({
      selector: 'app-master-navbar',
      templateUrl: './master-navbar.component.html',
      styleUrls: ['./master-navbar.component.css']
})
export class MasterNavbarComponent implements OnInit {

      menuList: IMenuDTO[] = []
      menusWithParentList: IMenuDTO[]
      profileMenu: IMenuDTO
      selectedMainMenuName: string
      searchedMenuList: IMenuDTO[]
      adminMenuCategoryList: any[]
      salesPersonMenuCategoryList: any[]
      selectedSubMenuList: IMenuDTO[]
      selectedMainMenuOrder: number
      selectedMenuCategory: any

      @ViewChild("menuTrigger") menuTrigger: MatMenuTrigger
      @ViewChild("searchBar") searchBar: ElementRef

      // Access.
      access: any

      // Flags.
      showFullSubMenu: boolean

      constructor(
            private loginService: UserloginService,
            private router: Router,
            private util: UtilityService,
            private localService: LocalService,
            private role: Role,
            private accessLevel: AccessLevel,
      ) { }

      ngOnInit() {
            this.initializeVariables()
            this.checkLogin()
            this.getMenus()
      }

      // Initialize global variables.
      initializeVariables() {
            this.menuList = []
            this.searchedMenuList = []
            this.menusWithParentList = []
            this.selectedSubMenuList = []
            this.adminMenuCategoryList = [
                  {
                        mainMenuOrder: 3,
                        menuDivision: [{subMenuIndex: 0, categoryName: "Check"}, {subMenuIndex: 2, categoryName: "View/Add"},
                              {subMenuIndex: 7, categoryName: "Manage Reports"}]
                  },
                  {
                        mainMenuOrder: 4,
                        menuDivision: [{subMenuIndex: 0, categoryName: "View+Add"}, {subMenuIndex: 7, categoryName: "Add"},
                              {subMenuIndex: 11, categoryName: "Manage Reports"}]
                  },
                  {
                        mainMenuOrder: 5,
                        menuDivision: [{subMenuIndex: 0, categoryName: "Companies"}, {subMenuIndex: 5, categoryName: "Colleges"},
                              {subMenuIndex: 10}]
                  },
                  {
                        mainMenuOrder: 6,
                        menuDivision: [{subMenuIndex: 0, categoryName: "Check"}, {subMenuIndex: 3, categoryName: "View/Add"},
                              {subMenuIndex: 6, categoryName: "Manage"}]
                  }
            ]

            this.salesPersonMenuCategoryList = [
                  {
                        mainMenuOrder: 3,
                        menuDivision: [{subMenuIndex: 0, categoryName: "Check"}, {subMenuIndex: 2, categoryName: "View/Add"},
                              {subMenuIndex: 3, categoryName: "Manage Reports"}]
                  },
                  {
                        mainMenuOrder: 4,
                        menuDivision: [{subMenuIndex: 0, categoryName: "View+Add"}, {subMenuIndex: 1, categoryName: "Manage Reports"},
                              {subMenuIndex: 2}]
                  },
                  {
                        mainMenuOrder: 5,
                        menuDivision: [{subMenuIndex: 0, categoryName: "Companies"}, {subMenuIndex: 5, categoryName: "Colleges"},
                              {subMenuIndex: 10}]
                  },
            ]

            this.selectedMainMenuOrder = -1

            // Role.
            if (this.role.ADMIN == this.localService.getJsonValue("roleName") || this.role.SALES_PERSON == this.localService.getJsonValue("roleName")) {
                  this.access = this.accessLevel.ADMIN_AND_SALESPERSON
            }

            // Flags.
            this.showFullSubMenu = false
      }

      // Check if user is logged in or not.
      checkLogin() {
            let loginToken = this.localService.getJsonValue('token')
            if (loginToken == null) {
                  this.util.rediectTo("login")
            }
      }

      // Get menus from local storage.
      getMenus(): void {
            this.menuList = this.localService.getJsonValue("menus")
            this.menusWithParentList = this.localService.getJsonValue("menus")
            if (this.menuList?.length > 0){
                  this.removeParentMenuFromMenus(this.menuList)
                  this.profileMenu = this.menuList[this.menuList.length - 1]
                  this.menuList.pop()
                  let firstName = this.localService.getJsonValue('firstName')
                  this.profileMenu.menuName = "Hi, " + firstName
                  this.parentMenuHighlight()
                  if (this.role.FACULTY == this.localService.getJsonValue("roleName")){
                        this.addQueryParamsForAddNewMenuForFaculty()
                  }
            }
      }

      // Delete session while logging out.
      deleteSession() {
            let tenantID = this.localService.getJsonValue('tenantID')
            let credentialID = this.localService.getJsonValue('credentialID')
            let token = this.localService.getJsonValue('token')
            this.localService.clearToken()
            this.router.navigateByUrl("/login")
            this.loginService.userLogout(tenantID, credentialID, token).subscribe((response: any) => {
                  console.log("Logout successful")
            },
                  err => {
                        console.error("Error -" + err.error);
                  })
      }

      // Redirect to home page.
      redirectToHomePage(): void {
            if (this.menuList.length > 0){
                  this.util.rediectTo(this.menuList[0].url)
            }
      }

      // Highlight the parent menu on child menu selection.
      parentMenuHighlight(): void {
            let parentMenu: string = this.router.url.substring(1)
            parentMenu = parentMenu.substring(0, parentMenu.indexOf('/'))
            parentMenu = "/" + parentMenu
            for (let i = 0; i < this.menuList.length; i++) {
                  if (this.menuList[i].url == parentMenu) {
                        if (this.role.FACULTY == this.localService.getJsonValue("roleName") && this.menuList[i].menuName == "Add a New"){
                              continue
                        }
                        this.menuList[i].class = "active-menu"
                  }
            }
            if (this.profileMenu.url == parentMenu) {
                  this.profileMenu.class = "active-menu"
            }
      }

      // Add query params for 'add a new' menu for faculty. 
      addQueryParamsForAddNewMenuForFaculty(): void{
            for (let i = 0; i < this.menuList.length; i++){
                  if (this.menuList[i].menuName == "Add a New"){
                        for (let j = 0; j< this.menuList[i].menus?.length; j++){
                              this.menuList[i].menus[j].opertaion = "add"
                        }
                  }
            }
      }

      // Below functions are for search.-------------------------------------------------------

      // On search input focus out.
      onFocusOut() {
            this.searchBar.nativeElement.value = null
      }

      // On giving input to search input.
      onSearch() {
            this.searchedMenuList = []
            let tempSearchVaue: string = this.searchBar.nativeElement.value
            let value: string = this.searchBar.nativeElement.value
            if (value) {
                  this.filterMenus(value, this.menusWithParentList)
            }
            if (this.searchedMenuList.length == 0) {
                  this.menuTrigger.closeMenu()
            } else {
                  this.menuTrigger.openMenu()
            }
            this.searchBar.nativeElement.focus()
            this.searchBar.nativeElement.value = tempSearchVaue
      }

      // Search for menus having text as searched input text.
      filterMenus(value: string, menus: IMenu[]) {
            menus.forEach(menu => {
                  if (menu.menuName.toLowerCase().includes(value.toLowerCase()) && menu.menuID != null && menu.isVisible) {
                        this.searchedMenuList.push(menu)
                  }
                  if (menu.menus) {
                        return this.filterMenus(value, menu.menus)
                  }
            })
      }

      // Set parent menu of menus to null.
      removeParentMenuFromMenus(menus: IMenuDTO[]) {
            for (let i = 0; i < menus.length; i++) {
                  menus[i].parentMenu = null
                  if (menus[i].menus?.length > 0) {
                        this.removeParentMenuFromMenus(menus[i].menus)
                  }
            }
      }

      // Print all menus.
      printParentMenus(menus: IMenuDTO[]): void {
            for (let i = 0; i < menus.length; i++) {
                  if (menus[i].menus?.length > 0) {
                        this.printParentMenus(menus[i].menus)
                  }
                  console.log(menus[i].parentMenu)
            }
      }

      // Get menu by menu id.
      getMenuFromID(id: string): IMenuDTO {
            if (!id) {
                  return
            }
            const menu = this.menuList?.find((menu) => {
                  if (menu.id === id)
                        return menu
                  !!menu?.menus?.find(m => m.id === id)
            })
            return menu
      }

      // When main menu is hovered in.
      selectedMenuHoverIn(menuName: string): void {
            let tempMenuCategoryList: any[] = []
            if (this.role.ADMIN == this.localService.getJsonValue("roleName")){
                  tempMenuCategoryList = this.adminMenuCategoryList
            }
            if (this.role.SALES_PERSON == this.localService.getJsonValue("roleName")){
                  tempMenuCategoryList = this.salesPersonMenuCategoryList
            }
            for (let i = 0; i < this.menuList.length; i++){
                  if (this.menuList[i].menuName == menuName){
                        this.selectedSubMenuList = this.menuList[i].menus
                        this.selectedMainMenuName = this.menuList[i].menuName
                        this.selectedMainMenuOrder = this.menuList[i].order
                        for (let j = 0; j < tempMenuCategoryList.length; j++){
                              if (tempMenuCategoryList[j].mainMenuOrder == this.selectedMainMenuOrder){
                                    this.selectedMenuCategory = tempMenuCategoryList[j]
                              }
                        }
                  }
            }
            this.showFullSubMenu = true
      }

      // When main menu is hovered out.
      selectedMenuHoverOut(): void {
            this.showFullSubMenu = false
            this.selectedMainMenuName = null
      }
}

interface ISearchMenu {
      menu: IMenu
      parentName: string

}