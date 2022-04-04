import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentBatchSessionPlanComponent } from './talent-batch-session-plan.component';

describe('TalentBatchSessionPlanComponent', () => {
  let component: TalentBatchSessionPlanComponent;
  let fixture: ComponentFixture<TalentBatchSessionPlanComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentBatchSessionPlanComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentBatchSessionPlanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
