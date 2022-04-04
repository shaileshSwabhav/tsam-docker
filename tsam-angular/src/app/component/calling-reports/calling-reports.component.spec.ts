import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CallingReportsComponent } from './calling-reports.component';

describe('CallingReportsComponent', () => {
  let component: CallingReportsComponent;
  let fixture: ComponentFixture<CallingReportsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CallingReportsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CallingReportsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
