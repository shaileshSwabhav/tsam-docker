import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentBatchDetailsComponent } from './talent-batch-details.component';

describe('TalentBatchDetailsComponent', () => {
  let component: TalentBatchDetailsComponent;
  let fixture: ComponentFixture<TalentBatchDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentBatchDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentBatchDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
