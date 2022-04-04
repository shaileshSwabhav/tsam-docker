import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PracticeProblemDetailsComponent } from './practice-problem-details.component';

describe('PracticeProblemDetailsComponent', () => {
  let component: PracticeProblemDetailsComponent;
  let fixture: ComponentFixture<PracticeProblemDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PracticeProblemDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PracticeProblemDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
