import { Directive, HostListener, Input, Self } from '@angular/core';
import { NgControl } from '@angular/forms';

@Directive({
  selector: '[appAllowNumbersOnly]'
})
export class AllowNumbersOnlyDirective {

  @Input() upperLimit: number = Infinity
  constructor(@Self() private ngControl: NgControl) { }
  /**
   * Hosts listener
   * @param event 
   * @returns true if string is not ".", "e", "+", "-"  
   */
  @HostListener('keydown', ['$event']) onKeyDowns(event: KeyboardEvent): boolean {
    let blockedKeys: string[] = [".", "e", "+", "-"];
    // this.ngControl.valueChanges.subscribe((number) => {
    //   if (number > this.upperLimit) {
    //     console.log(number, this.upperLimit)
    //   }
    // })
    for (let i = 0; i < blockedKeys.length; i++) {
      if (event.key === blockedKeys[i]) {
        return false;
      }
    }
    return true
  }
}
