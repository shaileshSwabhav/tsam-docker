import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NextActionReportComponent } from './next-action-report.component';

describe('NextActionReportComponent', () => {
  let component: NextActionReportComponent;
  let fixture: ComponentFixture<NextActionReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NextActionReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NextActionReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
