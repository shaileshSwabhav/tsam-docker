import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SessionPlanUpdateComponent } from './session-plan-update.component';

describe('SessionPlanUpdateComponent', () => {
  let component: SessionPlanUpdateComponent;
  let fixture: ComponentFixture<SessionPlanUpdateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SessionPlanUpdateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SessionPlanUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
