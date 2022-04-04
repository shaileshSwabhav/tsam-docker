import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PracticeProblemListComponent } from './practice-problem-list.component';

describe('PracticeProblemListComponent', () => {
  let component: PracticeProblemListComponent;
  let fixture: ComponentFixture<PracticeProblemListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PracticeProblemListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PracticeProblemListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
