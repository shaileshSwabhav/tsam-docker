import { ChangeDetectorRef, Component } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from './service/spinner/spinner.service';

@Component({
      selector: 'app-root',
      templateUrl: './app.component.html',
      styleUrls: ['./app.component.css']
})
export class AppComponent {
      title = 'Talent Sourcing and Mentoring';

      constructor(
            private router: Router,
            private modalService: NgbModal,
            private spinnerService: SpinnerService,
            private changeDetection: ChangeDetectorRef
      ) { }
      loadingMessage: string = "Loading..."

      ngAfterViewChecked(): void {
            this.loadingMessage = this.spinnerService.loadingMessage
            this.changeDetection.detectChanges()
      }

      ngOnInit() {
            // used to close modal if the component is changed (and is allowedto)
            this.router.events.subscribe(event => {
                  if (event instanceof NavigationEnd) {
                        this.spinnerService.isDisabled = false
                        // close all open modals
                        this.modalService.dismissAll();
                  }
            });
      }
}
