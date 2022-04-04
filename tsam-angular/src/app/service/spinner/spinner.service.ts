import { Injectable } from '@angular/core';
import { NgxSpinnerService } from 'ngx-bootstrap-spinner';
import { ISpinner } from 'src/app/models/spinner/spinner';

@Injectable({
  providedIn: 'root'
})
export class SpinnerService {

  isSpinnerShown: boolean = false
  constructor(private spinner: NgxSpinnerService) { }

  private _isDisabled: boolean = false

  set isDisabled(flag: boolean) {
    this._isDisabled = flag
  }

  get isDisabled(): boolean {
    return this._isDisabled
  }

  private _ongoingOperations: number = 0

  set ongoingOperations(no: number) {
    this._ongoingOperations = no
  }

  get ongoingOperations(): number {
    return this._ongoingOperations
  }

  private _loadingMessage: string

  set loadingMessage(msg: string) {
    this._loadingMessage = msg
  }

  get loadingMessage() {
    return this._loadingMessage
  }

  // override feature can be added for multiple spinner if needed in future.
  startSpinner(spinner?: ISpinner, name?: string) {
    this.ongoingOperations++
    if (!this.isDisabled && !this.isSpinnerShown) {
      this.isSpinnerShown = true
      this.spinner.show(name, spinner)
    }
  }

  stopSpinner(name?: string, debounce?: number) {
    if (--this.ongoingOperations <= 0 && !this.isDisabled && this.isSpinnerShown) {
      // test code
      if (this.ongoingOperations < 0) {
        console.warn("stopSpinner() has been called more times than start.")
      }
      this.isSpinnerShown = false
      this.spinner.hide(name, debounce)
    }
  }
}
