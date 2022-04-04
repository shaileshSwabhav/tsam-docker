import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'minute'
})
export class MinutePipe implements PipeTransform {

  transform(minutes: number): string {
    // let minutes: number = parseInt(strMinutes)
    if (minutes <= 0) {
      return "-"
    }

    // Calculate hours.
    let tempHours: any = Math.floor(minutes / 60)

    // Calculate minutes.
    let tempMinutes: any = minutes % 60

    // If hours is single digit then add 0 to it.
    if (tempHours.toString().length == 1){
      tempHours = "0" + tempHours
    }

    // If minutes is single digit then add 0 to it.
    if (tempMinutes.toString().length == 1){
      tempMinutes = "0" + tempMinutes
    }

    return (tempHours + ":" + tempMinutes + " hr(s)")
  }

}
