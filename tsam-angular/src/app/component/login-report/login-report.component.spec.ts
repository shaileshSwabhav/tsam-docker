import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoginReportComponent } from './login-report.component';

describe('LoginReportComponent', () => {
  let component: LoginReportComponent;
  let fixture: ComponentFixture<LoginReportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoginReportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoginReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
