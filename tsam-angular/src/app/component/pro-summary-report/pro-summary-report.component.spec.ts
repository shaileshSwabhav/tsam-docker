import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProSummaryReportComponent } from './pro-summary-report.component';

describe('ProSummaryReportComponent', () => {
  let component: ProSummaryReportComponent;
  let fixture: ComponentFixture<ProSummaryReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProSummaryReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProSummaryReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
