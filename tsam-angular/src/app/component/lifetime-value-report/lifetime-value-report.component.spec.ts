import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LifetimeValueReportComponent } from './lifetime-value-report.component';

describe('LifetimeValueReportComponent', () => {
  let component: LifetimeValueReportComponent;
  let fixture: ComponentFixture<LifetimeValueReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LifetimeValueReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LifetimeValueReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
