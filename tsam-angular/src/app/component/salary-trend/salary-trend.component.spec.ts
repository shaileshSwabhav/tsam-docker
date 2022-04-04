import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SalaryTrendComponent } from './salary-trend.component';

describe('SalaryTrendComponent', () => {
  let component: SalaryTrendComponent;
  let fixture: ComponentFixture<SalaryTrendComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SalaryTrendComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SalaryTrendComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
