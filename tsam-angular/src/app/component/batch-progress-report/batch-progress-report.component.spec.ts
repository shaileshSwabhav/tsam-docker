import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchProgressReportComponent } from './batch-progress-report.component';

describe('BatchProgressReportComponent', () => {
  let component: BatchProgressReportComponent;
  let fixture: ComponentFixture<BatchProgressReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchProgressReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchProgressReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
