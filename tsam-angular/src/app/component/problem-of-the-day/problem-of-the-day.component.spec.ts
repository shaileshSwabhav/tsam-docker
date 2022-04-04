import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProblemOfTheDayComponent } from './problem-of-the-day.component';

describe('ProblemOfTheDayComponent', () => {
  let component: ProblemOfTheDayComponent;
  let fixture: ComponentFixture<ProblemOfTheDayComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProblemOfTheDayComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProblemOfTheDayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
