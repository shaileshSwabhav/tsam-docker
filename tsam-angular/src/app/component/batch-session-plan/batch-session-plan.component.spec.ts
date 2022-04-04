import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchSessionPlanComponent } from './batch-session-plan.component';

describe('BatchSessionPlanComponent', () => {
  let component: BatchSessionPlanComponent;
  let fixture: ComponentFixture<BatchSessionPlanComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchSessionPlanComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchSessionPlanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
