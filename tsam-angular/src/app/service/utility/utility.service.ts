import { Injectable } from '@angular/core';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { Router } from '@angular/router';
import { IMenu, IPermission } from '../menu/menu.service';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class UtilityService {

  private date: any;
  constructor(
    private router: Router,
    private localService: LocalService
  ) {
    this.date = new Date()
  }

  isObjectEmpty(obj: Object): boolean {
    return Object.keys(obj).length === 0
  }


  // Change url 
  rediectTo(url: string) {
    this.router.navigate([url])
  }

  // Change url with param
  redirectTowithParam(url: string, param: any) {
    this.router.navigate([url, param]);
  }

  parseRFC3339Date() {

  }

  // Parse Date
  parseDate(str: string) {
    if (str == null) {
      return "Not Date"
    }
    let strarray = str.split("T")
    str = strarray[0]
    if (strarray.length == 1) {
      strarray = str.split("/")
      if (str != undefined || str != null || str != "") {
        str = strarray[2] + '-' + strarray[1] + '-' + strarray[0]
      }
    }
    let date = new Date(str)
    str = date.toDateString();
    if (str != 'Invalid Date') {
      return str
    }
    return ''
  }

  getPassout(str: string) {
    if (str == null) {
      return "Not Date"
    }
    let strarray = str.split("T")
    str = strarray[0]
    if (strarray.length == 1) {
      strarray = str.split("/")
      if (str != undefined || str != null || str != "") {
        str = strarray[2] + '-' + strarray[1] + '-' + strarray[0]
      }
    }
    let date = new Date(str)
    str = date.toDateString();
    str = date.toDateString().split(' ')[1] + " " + date.getFullYear()
    if (str != 'Invalid Date') {
      return str
    }
    return ''
  }

  // info logger
  info(...param) {
    console.log(this.getTime(), ...param)
  }

  // Warning logger
  warn(...param) {
    console.warn(this.getTime(), ...param)
  }

  // Warning logger
  error(...param) {
    console.error(this.getTime(), ...param)
  }

  // Get Time
  private getTime(): string {
    let time = this.date.toISOString();
    return time
  }

  // Delete Key with Value null From Object
  deleteNullValueIDFromObject(target) {
    if (typeof target === 'object') {
      for (const key in target) {
        if (key == "id" && target[key] == null) {
          delete target[key]
        }
        this.deleteNullValueIDFromObject(target[key]);
      }
    }
  }

  deleteNullValuePropertyFromObject(target) {
    if (typeof target === 'object') {
      for (const key in target) {
        if (target[key] == null) {
          delete target[key]
        } else {
        }
        this.deleteNullValuePropertyFromObject(target[key])
      }
    }
  }

  // Get Current Date in DD/MM/YYYY Format
  getCurrentDate(): string {
    let datestr: string
    let date = new Date()
    let param = date.getDate()
    if (param < 10) {
      datestr = '0' + param + '/'
    } else {
      datestr = param + '/'
    }
    param = date.getUTCMonth() + 1
    if (param < 10) {
      datestr += '0' + param + '/'
    } else {
      datestr += param + '/'
    }
    param = date.getFullYear()
    datestr += param
    // console.log(datestr)
    return datestr
  }

  getDateINDDMMYYYY(data: string) {
    let temp = data.split("-")
    let datestr: string = ""
    let param: number
    param = parseInt(temp[2]);
    if (param < 10) {
      datestr = '0' + param + '/'
    } else {
      datestr = param + '/'
    }
    param = parseInt(temp[1])
    if (param < 10) {
      datestr += '0' + param + '/'
    } else {
      datestr += param + '/'
    }
    datestr += temp[0]
    return datestr
  }

  getTotalExperience(experience: any[]) {
    let totalExperience: number = 0;
    for (let index = 0; index < experience.length; index++) {
      totalExperience += experience[index].yearOfExperience
    }
    if (totalExperience < 1) {
      return 0
    }
    return totalExperience
  }

  getErrorString(err) {
    if (navigator.onLine) {
      if (typeof err == 'object') {
        if (typeof err.error == 'string') {
          if (err.error == "") {
            return "something wrong, Please try later"
          }
          return err.error
        }
        if (typeof err.error == 'object') {
          if (typeof err.error.error == 'string') {
            return err.error.error
          }
          if (typeof err.error.error == undefined) {
            return "Something wrong Please try later"
          }
        }
      }
    }
    return "Please Check your Internet Connection"
  }

  formatTimeString(str: string) {
    let str1 = str.split(":")

    //it is pm if hours from 12 onwards
    let suffix = (parseInt(str1[0]) >= 12) ? 'pm' : 'am';

    //only -12 from hours if it is greater than 12 (if not back at mid night)
    let hours = (parseInt(str1[0]) > 12) ? parseInt(str1[0]) - 12 : parseInt(str1[0]);

    //if 00 then it is 12 am
    hours = (str1[0] == '00') ? 12 : hours;

    return hours + ":" + str1[1] + " " + suffix
  }

  parseTime(str: string) {
    let arr = str.split(":")
    str = arr[0] + ":" + arr[1]
    return str
  }

  getPackageInLPA(currentPackage: string): string {

    return currentPackage
  }

  compareStringNumber(num: string) {
    let number = parseInt(num)
    if (number < 12) {
      return "AM"
    }
    return "PM"
  }

  // Return general Types Value by Passing key and List
  getValueByKey(key: number, list: any[]) {
    if (list != null) {
      for (let index = 0; index < list.length; index++) {
        if (list[index].key == key) {
          return list[index].value
        }
      }
    }
  }

  // Return general Types Value by Passing key and List
  getKeyByValue(value: string, list: any[]) {
    if (list != null) {
      for (let index = 0; index < list.length; index++) {
        if (list[index].value == value) {
          return list[index].key
        }
      }
    }
  }

  // Get permission by url of component
  getPermission(url: string): IPermission {
    return this.getMenu(this.localService.getJsonValue("menus"), url)
  }

  //get permission by recursively iterating menus
  getMenu(menus: IMenu[], url: string): IPermission {
    for (let i = 0; i < menus.length; i++) {
      if (menus[i].url === url) {
        return menus[i].permission
      }
      // ***check should be based on whether menus is present in menu. -N
      if (menus[i].menus?.length > 0) {
        let permission = this.getMenu(menus[i].menus, url)
        if (permission) {
          return permission
        }
      }
    }
    return null
  }

  /**
* getPreviousMonday
* @returns date depending on the weekIndex value, 0 is for this week,
*  -1 is for next week, -2 for next to next week and so on.
*  1 is for previous week and so on.
*/
  getMonday(weekIndex: number = 0): Date {
    let date = new Date();
    let day = date.getDay()
    let prevMonday = new Date()
    let index = weekIndex * 7
    let subDate = ((day + 6) % 7) + index
    prevMonday.setDate(date.getDate() - subDate)
    return prevMonday
  }

  // =============================================FORM=============================================

  updateValueAndValiditors(formGroup: FormGroup) {
    for (const field in formGroup.controls) {
      let formElement = formGroup.get(field)
      if (formElement instanceof FormControl) {
        formElement.updateValueAndValidity()
      }
      if (formElement instanceof FormGroup) {
        this.updateValueAndValiditors(formElement)
      }
      if (formElement instanceof FormArray) {
        for (const control of formElement.controls) {
          this.updateValueAndValiditors(control as FormGroup)
        }
      }
    }
  }

  markArrayDirty(formArray: FormArray) {
    formArray.controls.forEach(control => {
      switch (control.constructor.name) {
        case "FormGroup":
          this.markGroupDirty(control as FormGroup);
          break;
        case "FormArray":
          this.markArrayDirty(control as FormArray);
          break;
        case "FormControl":
          this.markControlDirty(control as FormControl);
          break;
      }
    });
  }

  markGroupDirty(formGroup: FormGroup) {
    Object.keys(formGroup.controls).forEach(key => {
      switch (formGroup.get(key).constructor.name) {
        case "FormGroup":
          this.markGroupDirty(formGroup.get(key) as FormGroup);
          break;
        case "FormArray":
          this.markArrayDirty(formGroup.get(key) as FormArray);
          break;
        case "FormControl":
          this.markControlDirty(formGroup.get(key) as FormControl);
          break;
      }
    });
  }

  markControlDirty(formControl: FormControl) {
    formControl.markAsDirty();
  }

  public findInvalidControlsRecursive(formToInvestigate: FormGroup | FormArray): string[] {
    var invalidControls: string[] = [];
    let recursiveFunc = (form: FormGroup | FormArray) => {
      Object.keys(form.controls).forEach(field => {
        const control = form.get(field);
        if (control.invalid) invalidControls.push(field);
        if (control instanceof FormGroup) {
          recursiveFunc(control);
        } else if (control instanceof FormArray) {
          recursiveFunc(control);
        }
      });
    }
    recursiveFunc(formToInvestigate);
    return invalidControls;
  }

  // Convert time to AM PM format and remove seconds.
  convertTime(time: any): string {
    if (time == null){
      return null
    }
    time = time.toString().match(/^([01]\d|2[0-3])(:)([0-5]\d)(:[0-5]\d)?$/) || [time]
    if (time.length > 1) {
      time.pop()
      time = time.slice(1)
      time[5] = +time[0] < 12 ? ' AM' : ' PM'
      time[0] = +time[0] % 12 || 12
    }
    return time.join('')
  }

}
