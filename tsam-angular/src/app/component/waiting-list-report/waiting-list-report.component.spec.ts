import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WaitingListReportComponent } from './waiting-list-report.component';

describe('WaitingListReportComponent', () => {
  let component: WaitingListReportComponent;
  let fixture: ComponentFixture<WaitingListReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WaitingListReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WaitingListReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
