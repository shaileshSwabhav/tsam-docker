import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchTalentComponent } from './batch-talent.component';

describe('BatchTalentComponent', () => {
  let component: BatchTalentComponent;
  let fixture: ComponentFixture<BatchTalentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchTalentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchTalentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
