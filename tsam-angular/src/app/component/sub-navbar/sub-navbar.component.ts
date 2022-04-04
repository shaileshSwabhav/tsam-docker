import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatMenu } from '@angular/material/menu';
import { AddModalComponentService } from 'src/app/service/add-modal-component/add-model-component.service';
import { IMenu } from 'src/app/service/menu/menu.service';

@Component({
  selector: 'app-sub-navbar',
  templateUrl: './sub-navbar.component.html',
  styleUrls: ['./sub-navbar.component.css']
})
export class SubNavbarComponent implements OnInit {

  @ViewChild("menu", { static: true }) menu: MatMenu
  @Input() items: { menuName: string, menus: string[] }[]

  constructor(
    private addModalComponentService: AddModalComponentService,
  ) { }

  ngOnInit(): void {
  }

  // Add modal component by openinf modal of that component recognized by menu name.
  addModalComponent(menu: IMenu): void{
    this.addModalComponentService.openModalByMenuName(menu.url)
  }

}
