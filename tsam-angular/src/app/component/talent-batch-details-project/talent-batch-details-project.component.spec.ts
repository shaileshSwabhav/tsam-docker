import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentBatchDetailsProjectComponent } from './talent-batch-details-project.component';

describe('TalentBatchDetailsProjectComponent', () => {
  let component: TalentBatchDetailsProjectComponent;
  let fixture: ComponentFixture<TalentBatchDetailsProjectComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentBatchDetailsProjectComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentBatchDetailsProjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
