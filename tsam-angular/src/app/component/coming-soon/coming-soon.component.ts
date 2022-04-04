import { Component, OnInit } from '@angular/core';
import { ActivatedRouteSnapshot, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;

@Component({
  selector: 'app-coming-soon',
  templateUrl: './coming-soon.component.html',
  styleUrls: ['./coming-soon.component.css']
})
export class ComingSoonComponent implements OnInit {



  constructor(
    private spinnerService: SpinnerService,
  ) {
    this.showSpinner()
    // window.location.reload();
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  showSpinner(): void {
    // this.spinnerService.loadingMessage = "Loading..."
    // 
    // setTimeout(() => {  }, 10);
  }

}
