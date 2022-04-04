import { Directive, HostListener, Self } from '@angular/core';
import { NgControl } from '@angular/forms';

@Directive({
  selector: '[appAllowNumbersWithDecimal]'
})
export class AllowNumbersWithDecimalDirective {
  constructor(@Self() private ngControl: NgControl) { }
  /**
   * Hosts listener
   * @param event 
   * @returns true if string is not e", "+", "-"  
   */
  @HostListener('keydown', ['$event']) onKeyDowns(event: KeyboardEvent): boolean {
    let blockedKeys: string[] = ["e", "+", "-"];
    for (let i = 0; i < blockedKeys.length; i++) {
      if (event.key === blockedKeys[i]) {
        return false;
      }
    }
    return true
  }
}
